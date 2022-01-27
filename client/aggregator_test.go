//go:build integration
// +build integration

package client_test

import (
	"context"
	"testing"

	"github.com/calyptia/api/client"
	"github.com/calyptia/api/types"
)

func TestClient_CreateAggregator(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)
	withToken := withToken(t, asUser)

	got, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
		Name:                   "test-aggregator",
		AddHealthCheckPipeline: true,
	})
	wantEqual(t, err, nil)
	wantEqual(t, got.Name, "test-aggregator")
	wantNoEqual(t, got.Token, "")
	wantNoEqual(t, got.PrivateRSAKey, "")
	wantNoEqual(t, got.PublicRSAKey, "")
	wantNoTimeZero(t, got.CreatedAt)
	wantNoEqual(t, got.HealthCheckPipeline, nil)
	wantEqual(t, len(got.ResourceProfiles), 3)
}

func TestClient_Aggregators(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)
	withToken := withToken(t, asUser)

	created, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
		Name:                   "test-aggregator",
		AddHealthCheckPipeline: true,
	})
	wantEqual(t, err, nil)

	project := defaultProject(t, asUser)
	got, err := asUser.Aggregators(ctx, project.ID, types.AggregatorsParams{})
	wantEqual(t, err, nil)
	wantEqual(t, len(got), 1)
	wantEqual(t, got[0].ID, created.ID)
	wantEqual(t, got[0].Name, created.Name)
	wantEqual(t, got[0].Token, created.Token)
	wantEqual(t, got[0].PipelinesCount, uint64(1))
	wantEqual(t, got[0].CreatedAt, created.CreatedAt)
	wantEqual(t, got[0].UpdatedAt, created.CreatedAt)
}

func TestClient_Aggregator(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)
	withToken := withToken(t, asUser)

	created, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
		Name:                   "test-aggregator",
		AddHealthCheckPipeline: true,
	})
	wantEqual(t, err, nil)

	got, err := asUser.Aggregator(ctx, created.ID)
	wantEqual(t, err, nil)
	wantEqual(t, got.ID, created.ID)
	wantEqual(t, got.Name, created.Name)
	wantEqual(t, got.Token, created.Token)
	wantEqual(t, got.PipelinesCount, uint64(1))
	wantEqual(t, got.CreatedAt, created.CreatedAt)
	wantEqual(t, got.UpdatedAt, created.CreatedAt)
}

func TestClient_UpdateAggregator(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)
	withToken := withToken(t, asUser)

	created, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
		Name:                   "test-aggregator",
		AddHealthCheckPipeline: true,
	})
	wantEqual(t, err, nil)

	err = withToken.UpdateAggregator(ctx, created.ID, types.UpdateAggregator{
		Name: ptrStr("test-aggregator-updated"),
	})
	wantEqual(t, err, nil)
}

func TestClient_DeleteAggregator(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)
	withToken := withToken(t, asUser)

	created, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
		Name:                   "test-aggregator",
		AddHealthCheckPipeline: true,
	})
	wantEqual(t, err, nil)

	err = withToken.DeleteAggregator(ctx, created.ID)
	wantEqual(t, err, nil)
}

func setupAggregator(t *testing.T, withToken *client.Client) types.CreatedAggregator {
	t.Helper()
	ctx := context.Background()

	aggregator, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
		Name: "test-aggregator",
	})
	wantEqual(t, err, nil)

	return aggregator
}
