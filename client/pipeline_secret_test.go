package client_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/calyptia/api/types"
)

func TestClient_CreatePipelineSecret(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	got, err := asUser.CreatePipelineSecret(ctx, pipeline.ID, types.CreatePipelineSecret{
		Key:   "testsecret",
		Value: []byte("test-value"),
	})
	wantEqual(t, err, nil)
	wantNoEqual(t, got.ID, "")
	wantNoTimeZero(t, got.CreatedAt)
}

func TestClient_PipelineSecrets(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	secret, err := asUser.CreatePipelineSecret(ctx, pipeline.ID, types.CreatePipelineSecret{
		Key:   "testsecret",
		Value: []byte("test-value"),
	})
	wantEqual(t, err, nil)

	got, err := asUser.PipelineSecrets(ctx, pipeline.ID, types.PipelineSecretsParams{})
	wantEqual(t, err, nil)
	wantEqual(t, len(got), 1)

	wantEqual(t, got[0].ID, secret.ID)
	wantEqual(t, got[0].Key, "testsecret")
	wantNoEqual(t, got[0].Value, []byte("test-value"))
	wantEqual(t, got[0].CreatedAt, secret.CreatedAt)
	wantEqual(t, got[0].UpdatedAt, secret.CreatedAt)
}

func TestClient_PipelineSecret(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	secret, err := asUser.CreatePipelineSecret(ctx, pipeline.ID, types.CreatePipelineSecret{
		Key:   "testsecret",
		Value: []byte("test-value"),
	})
	wantEqual(t, err, nil)

	got, err := asUser.PipelineSecret(ctx, secret.ID)
	wantEqual(t, err, nil)
	wantEqual(t, got.ID, secret.ID)
	wantEqual(t, got.Key, "testsecret")
	wantNoEqual(t, got.Value, []byte("test-value"))
	wantEqual(t, got.CreatedAt, secret.CreatedAt)
	wantEqual(t, got.UpdatedAt, secret.CreatedAt)
}

func TestClient_UpdatePipelineSecret(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	secret, err := asUser.CreatePipelineSecret(ctx, pipeline.ID, types.CreatePipelineSecret{
		Key:   "testsecret",
		Value: []byte("test-value"),
	})
	wantEqual(t, err, nil)

	err = asUser.UpdatePipelineSecret(ctx, secret.ID, types.UpdatePipelineSecret{
		Key:   ptrStr("testsecretupdated"),
		Value: ptrBytes([]byte("test-value-updated")),
	})
	wantEqual(t, err, nil)
}

func TestClient_DeletePipelineSecret(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	secret, err := asUser.CreatePipelineSecret(ctx, pipeline.ID, types.CreatePipelineSecret{
		Key:   "testsecret",
		Value: []byte("test-value"),
	})
	wantEqual(t, err, nil)

	err = asUser.DeletePipelineSecret(ctx, secret.ID)
	wantEqual(t, err, nil)
}

func Test_decryptPipelineSecret(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	secret, err := asUser.CreatePipelineSecret(ctx, pipeline.ID, types.CreatePipelineSecret{
		Key:   "testsecret",
		Value: []byte("test-value"),
	})
	wantEqual(t, err, nil)

	got, err := asUser.PipelineSecret(ctx, secret.ID)
	wantEqual(t, err, nil)

	decrypted := decrypt(t, got.Value, aggregator.PrivateRSAKey)
	wantEqual(t, decrypted, []byte("test-value"))
}

// decrypt an encrypted value using implementation from cloud internals.
// You should pass a private RSA key to decrypt the value.
func decrypt(t *testing.T, val, key []byte) []byte {
	t.Helper()

	block, _ := pem.Decode(key)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	wantEqual(t, err, nil)

	hash := sha512.New()
	decrypted, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, val, nil)
	wantEqual(t, err, nil)

	return decrypted
}
