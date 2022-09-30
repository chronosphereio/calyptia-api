package client_test

import (
	"context"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/calyptia/api/types"
)

func TestClient_CreateCoreInstanceCheck(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	coreInstance := setupAggregator(t, withToken(t, asUser))

	got, err := asUser.CreateCoreInstanceCheck(ctx, coreInstance.ID,
		types.CreateCoreInstanceCheck{
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

func TestClient_CoreInstanceChecks(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	coreInstance := setupAggregator(t, withToken(t, asUser))

	t.Run("ok", func(t *testing.T) {
		check, err := asUser.CreateCoreInstanceCheck(ctx, coreInstance.ID,
			types.CreateCoreInstanceCheck{
				Protocol: types.PipelineProtocolTCP,
				Status:   types.CheckStatusOK,
				Retries:  3,
				Host:     "localhost",
				Port:     randPort(),
			},
		)
		wantEqual(t, err, nil)

		got, err := asUser.CoreInstanceChecks(ctx, coreInstance.ID, types.CoreInstanceChecksParams{})
		wantEqual(t, err, nil)
		wantEqual(t, len(got.Items), 1)

		wantEqual(t, got.Items[0].ID, check.ID)
		wantEqual(t, got.Items[0].CreatedAt, check.CreatedAt)
		wantEqual(t, got.Items[0].UpdatedAt, check.CreatedAt)

		_, err = asUser.Aggregator(ctx, coreInstance.ID)
		assert.NoError(t, err)

		_, err = asUser.CreateCoreInstanceCheck(ctx, coreInstance.ID,
			types.CreateCoreInstanceCheck{
				Protocol: types.PipelineProtocolTCP,
				Status:   types.CheckStatusNew,
				Retries:  3,
				Host:     "localhost",
				Port:     randPort(),
			},
		)
		wantEqual(t, err, nil)

		_, err = asUser.Aggregator(ctx, coreInstance.ID)
		assert.NoError(t, err)
	})

	t.Run("pagination", func(t *testing.T) {
		for i := 0; i < 9; i++ {
			_, err := asUser.CreateCoreInstanceCheck(ctx, coreInstance.ID,
				types.CreateCoreInstanceCheck{
					Protocol: types.PipelineProtocolTCP,
					Status:   types.CheckStatusOK,
					Retries:  3,
					Host:     "localhost",
					Port:     randPort(),
				},
			)
			wantEqual(t, err, nil)
		}

		allCoreInstanceChecks, err := asUser.CoreInstanceChecks(ctx, coreInstance.ID, types.CoreInstanceChecksParams{})
		wantEqual(t, err, nil)
		page1, err := asUser.CoreInstanceChecks(ctx, coreInstance.ID, types.CoreInstanceChecksParams{Last: ptrUint(3)})
		wantEqual(t, err, nil)
		page2, err := asUser.CoreInstanceChecks(ctx, coreInstance.ID, types.CoreInstanceChecksParams{Last: ptrUint(3), Before: page1.EndCursor})
		wantEqual(t, err, nil)

		want := allCoreInstanceChecks.Items[3:6]
		wantEqual(t, page2.Items, want)
	})
}

func TestClient_CoreInstanceCheck(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	coreInstance := setupAggregator(t, withToken(t, asUser))

	check, err := asUser.CreateCoreInstanceCheck(ctx, coreInstance.ID,
		types.CreateCoreInstanceCheck{
			Protocol: types.PipelineProtocolTCP,
			Status:   types.CheckStatusOK,
			Retries:  3,
			Host:     "localhost",
			Port:     randPort(),
		},
	)
	wantEqual(t, err, nil)

	got, err := asUser.CoreInstanceCheck(ctx, check.ID)
	wantEqual(t, err, nil)
	wantEqual(t, got.ID, check.ID)
	wantEqual(t, got.CreatedAt, check.CreatedAt)
	wantEqual(t, got.UpdatedAt, check.CreatedAt)
}

func TestClient_UpdateCoreInstanceCheck(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	coreInstance := setupAggregator(t, withToken(t, asUser))

	check, err := asUser.CreateCoreInstanceCheck(ctx, coreInstance.ID,
		types.CreateCoreInstanceCheck{
			Protocol: types.PipelineProtocolTCP,
			Status:   types.CheckStatusOK,
			Retries:  3,
			Host:     "localhost",
			Port:     randPort(),
		},
	)
	wantEqual(t, err, nil)

	err = asUser.UpdateCoreInstanceCheck(ctx, check.ID,
		types.UpdateCoreInstanceCheck{
			Protocol: ptr(types.PipelineProtocolUDP),
		},
	)
	wantEqual(t, err, nil)
}

func TestClient_DeleteCoreInstanceCheck(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	coreInstance := setupAggregator(t, withToken(t, asUser))

	check, err := asUser.CreateCoreInstanceCheck(ctx, coreInstance.ID,
		types.CreateCoreInstanceCheck{
			Protocol: types.PipelineProtocolTCP,
			Status:   types.CheckStatusOK,
			Retries:  3,
			Host:     "localhost",
			Port:     randPort(),
		},
	)
	wantEqual(t, err, nil)
	err = asUser.DeleteCoreInstanceCheck(ctx, check.ID)
	wantEqual(t, err, nil)
}
