package client_test

import (
	"context"
	"fmt"
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

		_, err = withToken.CreateAggregator(ctx, types.CreateAggregator{
			Name: "duplicate",
		})

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

	t.Run("pagination", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			_, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
				Name:    fmt.Sprintf("test-aggregator-%d", i),
				Version: types.DefaultAggregatorVersion,
			})
			wantEqual(t, err, nil)
		}

		page1, err := asUser.Aggregators(ctx, project.ID, types.AggregatorsParams{
			Last: ptrUint64(3),
		})
		wantEqual(t, err, nil)
		wantEqual(t, len(page1.Items), 3)
		wantNoEqual(t, page1.EndCursor, (*string)(nil))

		page2, err := asUser.Aggregators(ctx, project.ID, types.AggregatorsParams{
			Last:   ptrUint64(3),
			Before: page1.EndCursor,
		})
		wantEqual(t, err, nil)
		wantEqual(t, len(page2.Items), 3)
		wantNoEqual(t, page2.EndCursor, (*string)(nil))

		wantNoEqual(t, page1.Items, page2.Items)
		wantNoEqual(t, *page1.EndCursor, *page2.EndCursor)
		wantEqual(t, page1.Items[2].CreatedAt.After(page2.Items[0].CreatedAt), true)
	})

	t.Run("tags", func(t *testing.T) {
		for i := 10; i < 20; i++ {
			aggregator := types.CreateAggregator{
				Name:    fmt.Sprintf("test-aggregator-%d", i),
				Version: types.DefaultAggregatorVersion,
			}
			if i >= 15 {
				aggregator.Tags = append(aggregator.Tags, "tagone", "tagthree")
			} else {
				aggregator.Tags = append(aggregator.Tags, "tagtwo", "tagthree")
			}
			_, err := withToken.CreateAggregator(ctx, aggregator)
			wantEqual(t, err, nil)
		}

		opts := types.AggregatorsParams{}
		s := "tagone"
		opts.Tags = &s
		tag1, err := asUser.Aggregators(ctx, project.ID, opts)
		wantEqual(t, err, nil)
		wantEqual(t, len(tag1.Items), 5)
		wantEqual(t, tag1.Items[0].Tags, []string{"tagone", "tagthree"})

		s2 := "tagtwo"
		opts.Tags = &s2
		tag2, err := asUser.Aggregators(ctx, project.ID, opts)
		wantEqual(t, err, nil)
		wantEqual(t, len(tag2.Items), 5)
		wantEqual(t, tag2.Items[0].Tags, []string{"tagtwo", "tagthree"})

		s3 := "tagthree"
		opts.Tags = &s3
		tag3, err := asUser.Aggregators(ctx, project.ID, opts)
		wantEqual(t, err, nil)
		wantEqual(t, len(tag3.Items), 10)
		wantEqual(t, tag3.Items[0].Tags, []string{"tagone", "tagthree"})

		s4 := "tagone AND tagtwo"
		opts.Tags = &s4
		tag4, err := asUser.Aggregators(ctx, project.ID, opts)
		wantEqual(t, err, nil)
		wantEqual(t, len(tag4.Items), 0)
	})
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

		wantErrMsg(t, err, "invalid aggregator version")

		err = withToken.UpdateAggregator(ctx, created.ID, types.UpdateAggregator{
			Version: ptrStr(""),
		})

		wantErrMsg(t, err, "invalid aggregator version")
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
