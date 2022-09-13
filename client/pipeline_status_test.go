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

	t.Run("ok", func(t *testing.T) {
		_, err := withToken.UpdatePipeline(ctx, pipeline.ID, types.UpdatePipeline{
			Status: (*types.PipelineStatusKind)(ptrStr(string(types.PipelineStatusStarting))),
			// Pending acceptance in cloud
			// Events: []types.PipelineEvent{
			//	{
			//		Source:   types.PipelineEventSourceDeployment,
			//		Reason:   "Testing",
			//		Message:  "",
			//		LoggedAt: time.Now(),
			//	},
			// },
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
	})

	t.Run("pagination", func(t *testing.T) {
		for i := 0; i < 9; i++ {
			_, err := withToken.UpdatePipeline(ctx, pipeline.ID, types.UpdatePipeline{
				Status: (*types.PipelineStatusKind)(ptrStr(string(types.PipelineStatusStarting))),
			})
			wantEqual(t, err, nil)
			if i%2 == 0 {
				_, err := withToken.UpdatePipeline(ctx, pipeline.ID, types.UpdatePipeline{
					Status: (*types.PipelineStatusKind)(ptrStr(string(types.PipelineStatusFailed))),
				})
				wantEqual(t, err, nil)
			}
		}
		allStatusHistory, err := asUser.PipelineStatusHistory(ctx, pipeline.ID, types.PipelineStatusHistoryParams{})
		wantEqual(t, err, nil)
		page1, err := asUser.PipelineStatusHistory(ctx, pipeline.ID, types.PipelineStatusHistoryParams{Last: ptrUint(3)})
		wantEqual(t, err, nil)
		page2, err := asUser.PipelineStatusHistory(ctx, pipeline.ID, types.PipelineStatusHistoryParams{Last: ptrUint(3), Before: page1.EndCursor})
		wantEqual(t, err, nil)

		want := allStatusHistory.Items[3:6]
		wantEqual(t, page2.Items, want)
	})

	t.Run("pipeline statuses", func(t *testing.T) {
		statuses := []types.PipelineStatusKind{
			types.PipelineStatusFailed,
			types.PipelineStatusStarted,
			types.PipelineStatusStarting,
			types.PipelineStatusNew,
			types.PipelineStatusScaling,
		}

		for _, status := range statuses {
			_, err := withToken.UpdatePipeline(ctx, pipeline.ID, types.UpdatePipeline{
				Status: (*types.PipelineStatusKind)(ptrStr(string(status))),
			})
			wantEqual(t, err, nil)
		}
	})
}
