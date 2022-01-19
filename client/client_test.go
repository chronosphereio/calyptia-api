package client_test

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/calyptia/api/client"
	"github.com/calyptia/api/types"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

// var testBaseURL = "https://cloud-api-dev.calyptia.com"
var testBaseURL = "http://localhost:4000"

var (
	auth0Client = setupAuth0Client()
	// asGuest     = &client.Client{BaseURL: testBaseURL, Client: http.DefaultClient}
)

func setupAuth0Client() *Auth0Client {
	_ = godotenv.Load("../.env")

	clientID, ok := os.LookupEnv("TEST_AUTH0_CLIENT_ID")
	if !ok {
		fmt.Println("TEST_AUTH0_CLIENT_ID not set")
		return nil
	}

	clientSecret, ok := os.LookupEnv("TEST_AUTH0_CLIENT_SECRET")
	if !ok {
		fmt.Println("TEST_AUTH0_CLIENT_SECRET not set")
		return nil
	}

	domain, ok := os.LookupEnv("TEST_AUTH0_DOMAIN")
	if !ok {
		domain = "sso.calyptia.com"
	}
	if !strings.HasPrefix(domain, "http") {
		domain = "https://" + domain
	}

	aud, ok := os.LookupEnv("TEST_AUTH0_AUDIENCE")
	if !ok {
		aud = "http://localhost:3000"
	}

	return &Auth0Client{
		BaseURL:      domain,
		Client:       http.DefaultClient,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Audience:     aud,
		UserAgent:    "calyptia-api-client-test",
	}
}

func userClient(t *testing.T) *client.Client {
	t.Helper()

	if auth0Client == nil {
		t.Skip("auth client not initialized")
	}

	name := randStr(t)

	ctx := context.Background()
	tok, err := auth0Client.AccessToken(ctx, Auth0UserClaims{
		Email:         name + "@example.org",
		EmailVerified: true,
		Name:          name,
	})
	if err != nil {
		fmt.Printf("could not retrieve access token: %v\n", err)
		return nil
	}

	time.Sleep(time.Second) // prevent `Token used before issued`

	return &client.Client{
		BaseURL: testBaseURL,
		Client:  oauth2.NewClient(ctx, oauth2.StaticTokenSource(tok)),
	}
}

func withToken(t *testing.T, asUser *client.Client) *client.Client {
	t.Helper()

	tok := defaultToken(t, asUser)
	asUser.SetProjectToken(tok.Token)
	return asUser
}

func defaultProject(t *testing.T, asUser *client.Client) types.Project {
	t.Helper()

	ctx := context.Background()
	// every new user should have a default project.
	pp, err := asUser.Projects(ctx, types.ProjectsParams{Last: ptrUint64(1)})
	if err != nil {
		t.Fatal(err)
	}

	if len(pp) == 0 {
		t.Fatal("no default project")
	}

	return pp[0]
}

func defaultToken(t *testing.T, asUser *client.Client) types.Token {
	t.Helper()

	ctx := context.Background()
	project := defaultProject(t, asUser)

	// every new project should have a default token.
	tt, err := asUser.Tokens(ctx, project.ID, types.TokensParams{Last: ptrUint64(1)})
	if err != nil {
		t.Fatal(err)
	}

	if len(tt) == 0 {
		t.Fatal("no default token")
	}

	return tt[0]
}

func randStr(t *testing.T) string {
	t.Helper()
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%x", rand.Int63())
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

func ptrStr(s string) *string {
	return &s
}

func ptrBool(b bool) *bool {
	return &b
}

func ptrUint64(u uint64) *uint64 {
	return &u
}

func ptrUint(u uint) *uint {
	return &u
}
