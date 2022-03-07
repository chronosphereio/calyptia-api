package client_test

import (
	"context"
	"strings"
	"testing"

	"github.com/calyptia/api/client"
	"github.com/calyptia/api/types"
)

func wantErrMsg(t *testing.T, err error, msg string) {
	t.Helper()
	if err == nil || (err != nil && !strings.Contains(err.Error(), msg)) {
		t.Fatalf("want msg %q; got %v", msg, err)
	}
}

func TestClient_CreateAggregator(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)
	withToken := withToken(t, asUser)

	got, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
		Name:                   "test-aggregator",
		Version:                types.DefaultAggregatorVersion,
		AddHealthCheckPipeline: true,
	})

	wantEqual(t, err, nil)
	wantEqual(t, got.Version, types.DefaultAggregatorVersion)
	wantEqual(t, got.Name, "test-aggregator")
	wantNoEqual(t, got.Token, "")
	wantNoEqual(t, got.PrivateRSAKey, "")
	wantNoEqual(t, got.PublicRSAKey, "")
	wantNoTimeZero(t, got.CreatedAt)
	wantNoEqual(t, got.HealthCheckPipeline, nil)
	wantEqual(t, len(got.ResourceProfiles), 3)

	t.Run("name exists", func(t *testing.T) {
		_, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
			Name: "duplicate",
		})

		wantEqual(t, err, nil)

		got, err = withToken.CreateAggregator(ctx, types.CreateAggregator{
			Name: "duplicate",
		})

		wantEqual(t, got, types.CreatedAggregator{})
		wantErrMsg(t, err, "aggregator name already exists")
	})
}

func TestClient_Aggregators(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)
	withToken := withToken(t, asUser)

	created, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
		Name:                   "test-aggregator",
		Version:                types.DefaultAggregatorVersion,
		AddHealthCheckPipeline: true,
	})
	wantEqual(t, err, nil)

	project := defaultProject(t, asUser)
	got, err := asUser.Aggregators(ctx, project.ID, types.AggregatorsParams{})

	wantEqual(t, err, nil)
	wantEqual(t, len(got.Items), 1)
	wantEqual(t, got.Items[0].ID, created.ID)
	wantEqual(t, got.Items[0].Name, created.Name)
	wantEqual(t, got.Items[0].Version, created.Version)
	wantEqual(t, got.Items[0].Token, created.Token)
	wantEqual(t, got.Items[0].PipelinesCount, uint64(1))
	wantEqual(t, got.Items[0].CreatedAt, created.CreatedAt)
	wantEqual(t, got.Items[0].UpdatedAt, created.CreatedAt)
}

func TestClient_Aggregator(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)
	withToken := withToken(t, asUser)

	created, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
		Name:                   "test-aggregator",
		Version:                types.DefaultAggregatorVersion,
		AddHealthCheckPipeline: true,
	})
	wantEqual(t, err, nil)

	got, err := asUser.Aggregator(ctx, created.ID)
	wantEqual(t, err, nil)
	wantEqual(t, got.ID, created.ID)
	wantEqual(t, got.Name, created.Name)
	wantEqual(t, got.Version, created.Version)
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

	t.Run("version", func(t *testing.T) {
		err = withToken.UpdateAggregator(ctx, created.ID, types.UpdateAggregator{
			Version: ptrStr(types.DefaultAggregatorVersion),
		})

		wantEqual(t, err, nil)

		err = withToken.UpdateAggregator(ctx, created.ID, types.UpdateAggregator{
			Version: ptrStr("non-semver-version"),
		})

		wantNoEqual(t, err, nil)

		err = withToken.UpdateAggregator(ctx, created.ID, types.UpdateAggregator{
			Version: ptrStr(""),
		})

		wantNoEqual(t, err, nil)
	})
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
		Name:    "test-aggregator",
		Version: types.DefaultAggregatorVersion,
	})
	wantEqual(t, err, nil)

	return aggregator
}
