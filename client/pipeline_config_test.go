//go:build integration
// +build integration

package client_test

import (
	"context"
	"testing"

	"github.com/calyptia/api/types"
)

func TestClient_PipelineConfigHistory(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	_, err := asUser.UpdatePipeline(ctx, pipeline.ID, types.UpdatePipeline{
		RawConfig: ptrStr(testFbitConfigWithAddr3),
	})
	wantEqual(t, err, nil)

	got, err := asUser.PipelineConfigHistory(ctx, pipeline.ID, types.PipelineConfigHistoryParams{})
	wantEqual(t, err, nil)

	wantEqual(t, len(got), 2) // Initial config should be already there by default.

	wantNoEqual(t, got[0].ID, "")
	wantEqual(t, got[0].RawConfig, testFbitConfigWithAddr3)
	wantNoTimeZero(t, got[0].CreatedAt)

	wantNoEqual(t, got[1].ID, "")
	wantEqual(t, got[1].RawConfig, testFbitConfigWithAddr)
	wantNoTimeZero(t, got[1].CreatedAt)
}
