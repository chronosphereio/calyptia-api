package client_test

import (
	"context"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/calyptia/api/types"
)

func TestClient_CreatePipelineCheck(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	got, err := asUser.CreatePipelineCheck(ctx, pipeline.ID,
		types.CreatePipelineCheck{
			Protocol: types.PipelineProtocolTCP,
			Status:   types.CheckStatusOK,
			Retries:  3,
			Host:     "localhost",
			Port:     randPort(),
		},
	)
	wantEqual(t, err, nil)
	wantNoEqual(t, got.ID, "")
	wantNoTimeZero(t, got.CreatedAt)
}

func TestClient_PipelineChecks(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	t.Run("ok", func(t *testing.T) {
		check, err := asUser.CreatePipelineCheck(ctx, pipeline.ID,
			types.CreatePipelineCheck{
				Protocol: types.PipelineProtocolTCP,
				Status:   types.CheckStatusOK,
				Retries:  3,
				Host:     "localhost",
				Port:     randPort(),
			},
		)
		wantEqual(t, err, nil)

		got, err := asUser.PipelineChecks(ctx, pipeline.ID, types.PipelineChecksParams{})
		wantEqual(t, err, nil)
		wantEqual(t, len(got.Items), 1)

		wantEqual(t, got.Items[0].ID, check.ID)
		wantEqual(t, got.Items[0].CreatedAt, check.CreatedAt)
		wantEqual(t, got.Items[0].UpdatedAt, check.CreatedAt)

		pp, err := asUser.Pipeline(ctx, pipeline.ID, types.PipelineParams{})
		assert.NoError(t, err)
		assert.Equal(t, pp.ChecksOK, 1)
		assert.Equal(t, pp.ChecksTotal, 1)

		_, err = asUser.CreatePipelineCheck(ctx, pipeline.ID,
			types.CreatePipelineCheck{
				Protocol: types.PipelineProtocolTCP,
				Status:   types.CheckStatusNew,
				Retries:  3,
				Host:     "localhost",
				Port:     randPort(),
			},
		)
		wantEqual(t, err, nil)

		pp, err = asUser.Pipeline(ctx, pipeline.ID, types.PipelineParams{})
		assert.NoError(t, err)
		assert.Equal(t, pp.ChecksOK, 1)
		assert.Equal(t, pp.ChecksTotal, 2)
	})

	t.Run("pagination", func(t *testing.T) {
		for i := 0; i < 9; i++ {
			_, err := asUser.CreatePipelineCheck(ctx, pipeline.ID,
				types.CreatePipelineCheck{
					Protocol: types.PipelineProtocolTCP,
					Status:   types.CheckStatusOK,
					Retries:  3,
					Host:     "localhost",
					Port:     randPort(),
				},
			)
			wantEqual(t, err, nil)
		}

		allPipelineChecks, err := asUser.PipelineChecks(ctx, pipeline.ID, types.PipelineChecksParams{})
		wantEqual(t, err, nil)
		page1, err := asUser.PipelineChecks(ctx, pipeline.ID, types.PipelineChecksParams{Last: ptrUint(3)})
		wantEqual(t, err, nil)
		page2, err := asUser.PipelineChecks(ctx, pipeline.ID, types.PipelineChecksParams{Last: ptrUint(3), Before: page1.EndCursor})
		wantEqual(t, err, nil)

		want := allPipelineChecks.Items[3:6]
		wantEqual(t, page2.Items, want)
	})
}

func TestClient_PipelineCheck(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	check, err := asUser.CreatePipelineCheck(ctx, pipeline.ID,
		types.CreatePipelineCheck{
			Protocol: types.PipelineProtocolTCP,
			Status:   types.CheckStatusOK,
			Retries:  3,
			Host:     "localhost",
			Port:     randPort(),
		},
	)
	wantEqual(t, err, nil)

	got, err := asUser.PipelineCheck(ctx, check.ID)
	wantEqual(t, err, nil)
	wantEqual(t, got.ID, check.ID)
	wantEqual(t, got.CreatedAt, check.CreatedAt)
	wantEqual(t, got.UpdatedAt, check.CreatedAt)
}

func TestClient_UpdatePipelineCheck(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	check, err := asUser.CreatePipelineCheck(ctx, pipeline.ID,
		types.CreatePipelineCheck{
			Protocol: types.PipelineProtocolTCP,
			Status:   types.CheckStatusOK,
			Retries:  3,
			Host:     "localhost",
			Port:     randPort(),
		},
	)
	wantEqual(t, err, nil)

	err = asUser.UpdatePipelineCheck(ctx, check.ID,
		types.UpdatePipelineCheck{
			Protocol: ptr(types.PipelineProtocolUDP),
		},
	)
	wantEqual(t, err, nil)
}

func TestClient_DeletePipelineCheck(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	check, err := asUser.CreatePipelineCheck(ctx, pipeline.ID,
		types.CreatePipelineCheck{
			Protocol: types.PipelineProtocolTCP,
			Status:   types.CheckStatusOK,
			Retries:  3,
			Host:     "localhost",
			Port:     randPort(),
		},
	)
	wantEqual(t, err, nil)
	err = asUser.DeletePipelineCheck(ctx, check.ID)
	wantEqual(t, err, nil)
}
