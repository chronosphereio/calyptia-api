package client_test

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/calyptia/api/types"
)

func TestClient_PipelineConfigHistory(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	t.Run("ok", func(t *testing.T) {
		_, err := asUser.UpdatePipeline(ctx, pipeline.ID, types.UpdatePipeline{
			RawConfig: ptrStr(testFbitConfigWithAddr3),
		})
		wantEqual(t, err, nil)

		got, err := asUser.PipelineConfigHistory(ctx, pipeline.ID, types.PipelineConfigHistoryParams{})
		wantEqual(t, err, nil)

		wantEqual(t, len(got.Items), 2) // Initial config should be already there by default.

		wantNoEqual(t, got.Items[0].ID, "")
		wantEqual(t, got.Items[0].RawConfig, testFbitConfigWithAddr3)
		wantNoTimeZero(t, got.Items[0].CreatedAt)

		wantNoEqual(t, got.Items[1].ID, "")
		wantEqual(t, got.Items[1].RawConfig, testFbitConfigWithAddr)
		wantNoTimeZero(t, got.Items[1].CreatedAt)
	})

	t.Run("pagination", func(t *testing.T) {
		for i := 0; i < 9; i++ {
			_, err := asUser.UpdatePipeline(ctx, pipeline.ID, types.UpdatePipeline{
				RawConfig: ptrStr(strings.ReplaceAll(testFbitConfigWithAddr, "24224", strconv.Itoa(24224+i))),
			})
			wantEqual(t, err, nil)
		}
		allConfigs, err := asUser.PipelineConfigHistory(ctx, pipeline.ID, types.PipelineConfigHistoryParams{})
		wantEqual(t, err, nil)
		page1, err := asUser.PipelineConfigHistory(ctx, pipeline.ID, types.PipelineConfigHistoryParams{Last: ptrUint(3)})
		wantEqual(t, err, nil)
		page2, err := asUser.PipelineConfigHistory(ctx, pipeline.ID, types.PipelineConfigHistoryParams{Last: ptrUint(3), Before: page1.EndCursor})
		wantEqual(t, err, nil)

		want := allConfigs.Items[3:6]
		wantEqual(t, page2.Items, want)
	})
}
