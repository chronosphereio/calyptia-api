package client_test

import (
	"context"
	"testing"

	"github.com/calyptia/api/types"
)

func TestClient_PipelineStatusHistory(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	withToken := withToken(t, asUser)
	aggregator := setupAggregator(t, withToken)
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	_, err := withToken.UpdatePipeline(ctx, pipeline.ID, types.UpdatePipeline{
		Status: (*types.PipelineStatusKind)(ptrStr(string(types.PipelineStatusStarting))),
	})
	wantEqual(t, err, nil)

	got, err := asUser.PipelineStatusHistory(ctx, pipeline.ID, types.PipelineStatusHistoryParams{})
	wantEqual(t, err, nil)

	wantEqual(t, len(got.Items), 2) // Initial status "new" should be already there by default.

	wantNoEqual(t, got.Items[0].ID, "")
	wantEqual(t, got.Items[0].Config, got.Items[1].Config) // config should be the same as it wasn't updated.
	wantEqual(t, got.Items[0].Status, types.PipelineStatusStarting)
	wantNoTimeZero(t, got.Items[0].CreatedAt)

	wantNoEqual(t, got.Items[1].ID, "")
	wantNoEqual(t, got.Items[0].Config.ID, "")
	wantEqual(t, got.Items[0].Config.RawConfig, testFbitConfigWithAddr)
	wantNoTimeZero(t, got.Items[0].Config.CreatedAt)
	wantEqual(t, got.Items[1].Status, types.PipelineStatusNew)
	wantNoTimeZero(t, got.Items[1].CreatedAt)
}
