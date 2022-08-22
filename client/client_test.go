package client_test

//nolint:gci // using goimports
import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	math_rand "math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"golang.org/x/oauth2"

	"github.com/calyptia/api/client"
	"github.com/calyptia/api/types"
)

const (
	dockerHostGateway          = "host-gateway"
	defaultCloudImageLatestTag = "main"
)

var (
	dockerPool                *dockertest.Pool
	testBearerTokenPrivateKey jwk.RSAPrivateKey
	testCloudURL              string
)

var (
	hostIP                             = env("HOST_IP", dockerHostGateway)
	testCloudImage                     = env("TEST_CLOUD_IMAGE", "ghcr.io/calyptia/cloud:main")
	testCloudPort                      = env("TEST_CLOUD_PORT", "5000")
	testFluentbitConfigValidatorAPIKey = os.Getenv("TEST_FLUENTBIT_CONFIG_VALIDATOR_API_KEY")
	testFluentdConfigValidatorAPIKey   = os.Getenv("TEST_FLUENTD_CONFIG_VALIDATOR_API_KEY")
	testSMTPHost                       = env("TEST_SMTP_HOST", "smtp.mailtrap.io")
	testSMTPPort                       = env("TEST_SMTP_PORT", "465")
	testSMTPUsername                   = os.Getenv("TEST_SMTP_USERNAME")
	testSMTPPassword                   = os.Getenv("TEST_SMTP_PASSWORD")
	kubeConfig                         = os.Getenv("KUBECONFIG")
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(testMain(m))
}

func setupDockerPool() (*dockertest.Pool, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, err
	}

	return pool, pool.Client.Ping()
}

//nolint:gocyclo // this function setups all the components required by tests.
func testMain(m *testing.M) int {
	jwksServer, privateKey, err := setupJWKSServer()
	if err != nil {
		fmt.Printf("could not setup jwks server: %v\n", err)
		return 1
	}

	defer jwksServer.Close()

	testBearerTokenPrivateKey = privateKey

	dockerPool, err = setupDockerPool()
	if err != nil {
		fmt.Printf("error setting up docker pool: %v\n", err)
		return 1
	}

	postgres, err := setupPostgres(dockerPool)
	if err != nil {
		fmt.Printf("could not setup postgres: %v\n", err)
		return 1
	}

	defer func(postgres *dockertest.Resource) {
		err := postgres.Close()
		if err != nil {
			return
		}
	}(postgres)

	defer func(pool *dockertest.Pool, r *dockertest.Resource) {
		err := pool.Purge(r)
		if err != nil {
			return
		}
	}(dockerPool, postgres)

	influx, err := setupInflux(dockerPool)
	if err != nil {
		fmt.Printf("could not setup influx: %v\n", err)
		return 1
	}

	defer func(influx *dockertest.Resource) {
		err := influx.Close()
		if err != nil {
			return
		}
	}(influx)

	defer func(pool *dockertest.Pool, r *dockertest.Resource) {
		err := pool.Purge(r)
		if err != nil {
			return
		}
	}(dockerPool, influx)

	err = pingPostgres(postgres)
	if err != nil {
		fmt.Printf("could not ping postgres: %v\n", err)
		return 1
	}

	err = pingInflux(influx)
	if err != nil {
		fmt.Printf("could not ping influx: %v\n", err)
		return 1
	}

	jwksURL, err := url.Parse(jwksServer.URL)
	if err != nil {
		fmt.Printf("could not parse jwks url: %v\n", err)
		return 1
	}

	jwksURL.Host = "host.docker.internal:" + jwksURL.Port()
	jwksURL.Path = "/.well-known/jwks.json"

	cloud, err := setupCloud(dockerPool, setupCloudConfig{
		port:                           testCloudPort,
		jwksURL:                        jwksURL.String(),
		accessTokenAudience:            "http://cloud-api-testing.localhost",
		accessTokenIssuer:              "http://cloud-api-testing.localhost",
		postgresDSN:                    "postgresql://postgres@host.docker.internal:" + postgres.GetPort("5432/tcp") + "?sslmode=disable",
		influxServer:                   "http://host.docker.internal:" + influx.GetPort("8086/tcp"),
		fluentBitConfigValidatorAPIKey: testFluentbitConfigValidatorAPIKey,
		fluentdConfigValidatorAPIKey:   testFluentdConfigValidatorAPIKey,
		smtpHost:                       testSMTPHost,
		smtpPort:                       testSMTPPort,
		smtpUsername:                   testSMTPUsername,
		smtpPassword:                   testSMTPPassword,
	})
	if err != nil {
		fmt.Printf("could not setup cloud: %v\n", err)
		return 1
	}

	defer func(cloud *dockertest.Resource) {
		err := cloud.Close()
		if err != nil {
			return
		}
	}(cloud)

	defer func(pool *dockertest.Pool, r *dockertest.Resource) {
		err := pool.Purge(r)
		if err != nil {
			return
		}
	}(dockerPool, cloud)

	if testing.Verbose() {
		for _, name := range []string{cloud.Container.ID} {
			name := name
			go func() {
				err := dockerPool.Client.Logs(docker.LogsOptions{
					Container:   name,
					ErrorStream: os.Stderr,
					Stderr:      true,
					Follow:      true,
				})
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}()
		}
	}

	err = pingCloud(cloud)
	if err != nil {
		fmt.Printf("could not get cloud base url: %v\n", err)
		return 1
	}

	testCloudURL = "http://" + cloud.GetHostPort(testCloudPort+"/tcp")

	return m.Run()
}

func setupJWKSServer() (*httptest.Server, jwk.RSAPrivateKey, error) {
	raw, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("could not generate new RSA private key: %w", err)
	}

	key, err := jwk.New(raw)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create symmetric key: %w", err)
	}

	err = key.Set(jwk.KeyIDKey, "test-kid")
	if err != nil {
		return nil, nil, fmt.Errorf("could not set jwt key id: %w", err)
	}

	priv, ok := key.(jwk.RSAPrivateKey)
	if !ok {
		return nil, nil, fmt.Errorf("expected jwk.RSAPrivateKey, got %T", key)
	}

	pub, err := priv.PublicKey()
	if err != nil {
		return nil, nil, fmt.Errorf("could not get public key: %w", err)
	}

	resp, err := json.Marshal(pub)
	if err != nil {
		return nil, nil, fmt.Errorf("could not json marshal public key: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(resp)
		if err != nil && !errors.Is(err, syscall.EPIPE) {
			fmt.Printf("could not write jkws response: %v\n", err)
		}
	})

	srv := httptest.NewUnstartedServer(mux)

	defer srv.Start()

	// The behavior of httptest.NewServer differs from OS X / Linux
	// the binding by default will do 127.0.0.1:0 and that is translatable
	// as the host-gateway on containers, that fails on Linux as the scope
	// of the binding is limited to local and doesn't expands to *:0 so
	// we force a particular known address (172.17.0.1) as the default
	// docker bridge address.
	if hostIP != dockerHostGateway {
		l, err := net.Listen("tcp", fmt.Sprintf("%s:0", hostIP))
		if err != nil {
			return nil, nil, err
		}

		srv.Listener = l
	}

	return srv, priv, nil
}

type bearerTokenClaims struct {
	Email         string
	EmailVerified bool
	Name          string
}

func issueBearerToken(key jwk.RSAPrivateKey, claims bearerTokenClaims) (*oauth2.Token, error) {
	now := time.Now()
	exp := now.Add(time.Hour * 24 * 7)
	tok, err := jwt.NewBuilder().
		Audience([]string{"http://cloud-api-testing.localhost"}).
		Issuer("http://cloud-api-testing.localhost").
		IssuedAt(now).
		NotBefore(now).
		Expiration(exp).
		Subject("test-subject").
		Claim("https://cloud.calyptia.com/email", claims.Email).
		Claim("https://cloud.calyptia.com/name", claims.Name).
		Claim("https://cloud.calyptia.com/email_verified", claims.EmailVerified).
		Build()
	if err != nil {
		return nil, fmt.Errorf("could not build jwt: %w", err)
	}

	headers := jws.NewHeaders()
	err = headers.Set("kid", "test-kid")
	if err != nil {
		return nil, fmt.Errorf("could not set jwt kid header: %w", err)
	}

	b, err := jwt.Sign(tok, jwa.RS256, key, jwt.WithHeaders(headers))
	if err != nil {
		return nil, fmt.Errorf("could not sign jwt: %w", err)
	}

	return &oauth2.Token{
		AccessToken:  string(b),
		TokenType:    "Bearer",
		RefreshToken: "",
		Expiry:       exp,
	}, nil
}

func setupPostgres(pool *dockertest.Pool) (*dockertest.Resource, error) {
	return pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "postgres",
		Env:          []string{"POSTGRES_HOST_AUTH_METHOD=trust"},
		ExposedPorts: []string{"5432"},
	}, func(hc *docker.HostConfig) {
		hc.AutoRemove = true
		hc.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
}

func pingPostgres(postgres *dockertest.Resource) error {
	return retry(func() error {
		hostPort := postgres.GetHostPort("5432/tcp")
		if hostPort == "" {
			return errors.New("postgres host-port not ready")
		}

		db, err := sql.Open("postgres", "postgresql://postgres@"+hostPort+"?sslmode=disable")
		if err != nil {
			return fmt.Errorf("could not open postgres db: %w", err)
		}

		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				return
			}
		}(db)

		if err := db.Ping(); err != nil {
			return fmt.Errorf("could not ping postgres db: %w", err)
		}

		return nil
	})
}

func setupInflux(pool *dockertest.Pool) (*dockertest.Resource, error) {
	return pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "influxdb",
		Env: []string{
			"DOCKER_INFLUXDB_INIT_MODE=setup",
			"DOCKER_INFLUXDB_INIT_USERNAME=my-user",
			"DOCKER_INFLUXDB_INIT_PASSWORD=my-password",
			"DOCKER_INFLUXDB_INIT_ORG=cloud-api",
			"DOCKER_INFLUXDB_INIT_BUCKET=cloud-api",
			"DOCKER_INFLUXDB_INIT_ADMIN_TOKEN=cloud-api",
		},
		ExposedPorts: []string{"8086"},
	}, func(hc *docker.HostConfig) {
		hc.AutoRemove = true
		hc.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
}

func pingInflux(influx *dockertest.Resource) error {
	return retry(func() error {
		hostPort := influx.GetHostPort("8086/tcp")
		if hostPort == "" {
			return errors.New("influx host-port not ready")
		}

		client := influxdb.NewClient("http://"+hostPort, "cloud-api")
		defer client.Close()

		ok, err := client.Ping(context.Background())
		if err != nil {
			return err
		}

		if !ok {
			return errors.New("influx is not ready")
		}

		return nil
	})
}

type setupCloudConfig struct {
	port                           string
	jwksURL                        string
	accessTokenAudience            string
	accessTokenIssuer              string
	postgresDSN                    string
	influxServer                   string
	fluentBitConfigValidatorAPIKey string
	fluentdConfigValidatorAPIKey   string
	smtpHost                       string
	smtpPort                       string
	smtpUsername                   string
	smtpPassword                   string
}

func getAuthConfigForImage(image string) (docker.AuthConfiguration, error) {
	var authConfig docker.AuthConfiguration
	if !strings.Contains(image, "://") {
		image = "//" + image
	}
	parsedURL, err := url.Parse(image)
	if err != nil {
		return authConfig, errors.New("local image, skipping auth config")
	}

	authConfs, err := docker.NewAuthConfigurationsFromDockerCfg()
	if err != nil {
		return authConfig, fmt.Errorf("could not read auth config: %w", err)
	}

	authConfig, ok := authConfs.Configs[parsedURL.Host]
	if !ok {
		return authConfig, fmt.Errorf("not found auth config for host: %q", parsedURL.Host)
	}

	return authConfig, nil
}

func setupCloud(pool *dockertest.Pool, conf setupCloudConfig) (*dockertest.Resource, error) {
	var cloudImageTag string

	splitDockerImage := strings.SplitN(testCloudImage, ":", 2)
	if len(splitDockerImage) == 2 {
		cloudImageTag = splitDockerImage[1]
		testCloudImage = splitDockerImage[0]
	} else {
		cloudImageTag = defaultCloudImageLatestTag
	}

	authConfig, err := getAuthConfigForImage(splitDockerImage[0])
	if err != nil {
		fmt.Println(err.Error())
	}

	return pool.RunWithOptions(&dockertest.RunOptions{
		Auth:       authConfig,
		Repository: testCloudImage,
		Tag:        cloudImageTag,
		Env: []string{
			"PORT=" + conf.port,
			"ORIGIN=http://localhost:" + conf.port,
			"JWKS_URL=" + conf.jwksURL,
			"ACCESS_TOKEN_AUD=" + conf.accessTokenAudience,
			"ACCESS_TOKEN_ISS=" + conf.accessTokenIssuer,
			"POSTGRES_DSN=" + conf.postgresDSN,
			"INFLUX_SERVER=" + conf.influxServer,
			"FLUENTBIT_CONFIG_VALIDATOR_API_KEY=" + conf.fluentBitConfigValidatorAPIKey,
			"FLUENTD_CONFIG_VALIDATOR_API_KEY=" + conf.fluentdConfigValidatorAPIKey,
			"SMTP_HOST=" + conf.smtpHost,
			"SMTP_PORT=" + conf.smtpPort,
			"SMTP_USERNAME=" + conf.smtpUsername,
			"SMTP_PASSWORD=" + conf.smtpPassword,
			"ALLOWED_ORIGINS=http://cloud-api-testing.localhost",
			"DEBUG=true",
		},
		ExposedPorts: []string{conf.port},
		ExtraHosts:   []string{"host.docker.internal:" + hostIP},
	}, func(hc *docker.HostConfig) {
		hc.AutoRemove = true
		hc.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
}

func pingCloud(cloud *dockertest.Resource) error {
	return retry(func() error {
		hostPort := cloud.GetHostPort(testCloudPort + "/tcp")
		if hostPort == "" {
			return errors.New("cloud host-port not ready")
		}

		//nolint //http.Get is okay on this context
		resp, err := http.Get("http://" + hostPort + "/healthz")
		if err != nil {
			return err
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("healthz returned %d", resp.StatusCode)
		}

		var want struct {
			OK bool `json:"ok"`
		}
		err = json.NewDecoder(resp.Body).Decode(&want)
		if err != nil {
			return err
		}

		if !want.OK {
			return errors.New("cloud not ready")
		}

		return nil
	})
}

func userClient(t *testing.T) *client.Client {
	t.Helper()

	if testCloudURL == "" {
		t.Skip("testBaseURL not set")
	}

	if testBearerTokenPrivateKey == nil {
		t.Skip("bearer token private key not set")
	}

	name := randStr(t)

	ctx := context.Background()
	tok, err := issueBearerToken(testBearerTokenPrivateKey, bearerTokenClaims{
		Email:         name + "@example.org",
		Name:          name,
		EmailVerified: true,
	})
	if err != nil {
		fmt.Printf("could not retrieve access token: %v\n", err)
		return nil
	}

	return &client.Client{
		BaseURL: testCloudURL,
		Client:  oauth2.NewClient(ctx, oauth2.StaticTokenSource(tok)),
	}
}

func withToken(t *testing.T, asUser *client.Client) *client.Client {
	t.Helper()

	tok := defaultToken(t, asUser)
	asUser.SetProjectToken(tok.Token)
	return asUser
}

var errNoKubeConfig = errors.New("no kubeconfig")

func setupCoreInstance(pool *dockertest.Pool, baseURL string, token types.Token) (*dockertest.Resource, error) {
	// Determine kubernetes config location
	if kubeConfig == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		kubeConfig = filepath.Join(home, ".kube", "config")
	}

	if _, err := os.Stat(kubeConfig); err != nil {
		if os.IsNotExist(err) {
			return nil, errNoKubeConfig
		}

		return nil, err
	}

	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	res, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "ghcr.io/calyptia/core",
		Tag:        "latest",
		Privileged: true,
		Env: []string{
			"AGGREGATOR_FLUENTBIT_CLOUD_URL=" + "http://host.docker.internal:" + parsed.Port(),
			"AGGREGATOR_FLUENTBIT_TLS=off",
			"AGGREGATOR_FLUENTBIT_TLS_VERIFY=off",
			"PROJECT_TOKEN=" + token.Token,
			"KUBECONFIG=/opt/calyptia-aggregator/kubeconfig",
		},
		ExtraHosts: []string{"host.docker.internal:" + hostIP},
	}, func(hc *docker.HostConfig) {
		hc.Mounts = []docker.HostMount{
			{
				Target: "/opt/calyptia-aggregator/kubeconfig",
				Source: kubeConfig,
				Type:   "bind",
			},
		}
	})
	if err != nil {
		return nil, err
	}

	if testing.Verbose() {
		go func() {
			name := res.Container.ID
			err := pool.Client.Logs(docker.LogsOptions{
				Container:    name,
				OutputStream: os.Stdout,
				ErrorStream:  os.Stderr,
				Stdout:       true,
				Stderr:       true,
				Follow:       true,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to get logs for container %q: %v\n", name, err)
			}
		}()
	}

	return res, nil
}

func defaultProject(t *testing.T, asUser *client.Client) types.Project {
	t.Helper()

	ctx := context.Background()
	// every new user should have a default project.
	pp, err := asUser.Projects(ctx, types.ProjectsParams{Last: ptrUint(1)})
	if err != nil {
		t.Fatal(err)
	}

	if len(pp.Items) == 0 {
		t.Fatal("no default project")
	}

	return pp.Items[0]
}

func defaultToken(t *testing.T, asUser *client.Client) types.Token {
	t.Helper()

	ctx := context.Background()
	project := defaultProject(t, asUser)

	// every new project should have a default token.
	tt, err := asUser.Tokens(ctx, project.ID, types.TokensParams{Last: ptrUint(1)})
	if err != nil {
		t.Fatal(err)
	}

	if len(tt.Items) == 0 {
		t.Fatal("no default token")
	}

	return tt.Items[0]
}

func randStr(t *testing.T) string {
	t.Helper()
	math_rand.Seed(time.Now().UnixNano())
	//nolint // math_rand uses math/rand import
	return fmt.Sprintf("%x", math_rand.Int63())
}

func wantNoTimeZero(t *testing.T, ts time.Time) {
	t.Helper()
	if ts.IsZero() {
		t.Fatal("time is zero")
	}
}

func wantNoEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if reflect.DeepEqual(got, want) {
		t.Fatalf("want not %+v; got %+v", want, got)
	}
}

func wantEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("want %+v; got %+v", want, got)
	}
}

func wantErrMsg(t *testing.T, err error, msg string) {
	t.Helper()
	if err == nil || (err != nil && !strings.Contains(err.Error(), msg)) {
		t.Fatalf("want msg %q; got %v", msg, err)
	}
}

func ptrStr(s string) *string {
	return &s
}

func ptrBool(b bool) *bool {
	return &b
}

func ptrUint(u uint) *uint {
	return &u
}

func ptrBytes(b []byte) *[]byte {
	return &b
}

func randUUID(t *testing.T) string {
	t.Helper()
	return uuid.Must(uuid.NewRandom()).String()
}

func retry(op func() error) error {
	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = time.Second * 5
	bo.MaxElapsedTime = time.Minute
	if err := backoff.Retry(op, bo); err != nil {
		if bo.NextBackOff() == backoff.Stop {
			return fmt.Errorf("reached retry deadline: %w", err)
		}
		return err
	}
	return nil
}

func env(key, fallback string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return v
}
