package client_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/calyptia/api/types"
)

func TestClient_CreatePipelinePort(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)

	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	got, err := asUser.CreatePipelinePort(ctx, pipeline.ID, types.CreatePipelinePort{
		Protocol:     "tcp",
		FrontendPort: 4000,
		BackendPort:  4000,
		Endpoint:     "http://localhost:4000",
	})
	wantEqual(t, err, nil)
	wantNoEqual(t, got.ID, "")
	wantNoTimeZero(t, got.CreatedAt)
}

func TestClient_PipelinePorts(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)

	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)
	t.Run("ok", func(t *testing.T) {
		port, err := asUser.CreatePipelinePort(ctx, pipeline.ID, types.CreatePipelinePort{
			Protocol:     "tcp",
			FrontendPort: 4000,
			BackendPort:  4000,
			Endpoint:     "http://localhost:4000",
		})
		wantEqual(t, err, nil)

		got, err := asUser.PipelinePorts(ctx, pipeline.ID, types.PipelinePortsParams{})
		wantEqual(t, err, nil)

		wantEqual(t, len(got.Items), 2) // The setup pipeline already contains ports in its config.

		wantEqual(t, got.Items[0].ID, port.ID)
		wantEqual(t, got.Items[0].Protocol, "tcp")
		wantEqual(t, got.Items[0].FrontendPort, uint(4000))
		wantEqual(t, got.Items[0].BackendPort, uint(4000))
		wantEqual(t, got.Items[0].Endpoint, "http://localhost:4000")
		wantEqual(t, got.Items[0].CreatedAt, port.CreatedAt)
		wantEqual(t, got.Items[0].UpdatedAt, port.CreatedAt)

		wantNoEqual(t, got.Items[1].ID, "")
		wantEqual(t, got.Items[1].Protocol, "tcp")           // from the test config.
		wantEqual(t, got.Items[1].FrontendPort, uint(24224)) // from the test config.
		wantEqual(t, got.Items[1].BackendPort, uint(24224))  // from the test config.
		wantEqual(t, got.Items[1].Endpoint, "")
		wantNoTimeZero(t, got.Items[1].CreatedAt)
		wantNoTimeZero(t, got.Items[1].UpdatedAt)
	})

	t.Run("pagination", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			_, err := asUser.CreatePipelinePort(ctx, pipeline.ID, types.CreatePipelinePort{
				Protocol:     "tcp",
				FrontendPort: 400 + uint(i),
				BackendPort:  400 + uint(i),
				Endpoint:     fmt.Sprintf("http://localhost:400%d", i),
			})
			wantEqual(t, err, nil)
		}
		allPorts, err := asUser.PipelinePorts(ctx, pipeline.ID, types.PipelinePortsParams{})
		wantEqual(t, err, nil)
		page1, err := asUser.PipelinePorts(ctx, pipeline.ID, types.PipelinePortsParams{Last: ptrUint64(3)})
		wantEqual(t, err, nil)
		page2, err := asUser.PipelinePorts(ctx, pipeline.ID, types.PipelinePortsParams{Last: ptrUint64(3), Before: page1.EndCursor})
		wantEqual(t, err, nil)

		want := allPorts.Items[3:6]
		wantEqual(t, page2.Items, want)
	})
}

func TestClient_PipelinePort(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)

	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	port, err := asUser.CreatePipelinePort(ctx, pipeline.ID, types.CreatePipelinePort{
		Protocol:     "tcp",
		FrontendPort: 4000,
		BackendPort:  4000,
		Endpoint:     "http://localhost:4000",
	})
	wantEqual(t, err, nil)

	got, err := asUser.PipelinePort(ctx, port.ID)
	wantEqual(t, err, nil)

	wantEqual(t, got.ID, port.ID)
	wantEqual(t, got.Protocol, "tcp")
	wantEqual(t, got.FrontendPort, uint(4000))
	wantEqual(t, got.BackendPort, uint(4000))
	wantEqual(t, got.Endpoint, "http://localhost:4000")
	wantEqual(t, got.CreatedAt, port.CreatedAt)
	wantEqual(t, got.UpdatedAt, port.CreatedAt)
}

func TestClient_UpdatePipelinePort(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	withToken := withToken(t, asUser)

	aggregator := setupAggregator(t, withToken)
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	port, err := asUser.CreatePipelinePort(ctx, pipeline.ID, types.CreatePipelinePort{
		Protocol:     "tcp",
		FrontendPort: 4000,
		BackendPort:  4000,
		Endpoint:     "http://localhost:4000",
	})
	wantEqual(t, err, nil)

	err = withToken.UpdatePipelinePort(ctx, port.ID, types.UpdatePipelinePort{
		Protocol:     ptrStr("udp"),
		FrontendPort: ptrUint(4001),
		BackendPort:  ptrUint(4001),
		Endpoint:     ptrStr("http://localhost:4001"),
	})
	wantEqual(t, err, nil)
}

func TestClient_DeletePipelinePort(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)

	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	port, err := asUser.CreatePipelinePort(ctx, pipeline.ID, types.CreatePipelinePort{
		Protocol:     "tcp",
		FrontendPort: 4000,
		BackendPort:  4000,
		Endpoint:     "http://localhost:4000",
	})
	wantEqual(t, err, nil)

	err = asUser.DeletePipelinePort(ctx, port.ID)
	wantEqual(t, err, nil)
}
