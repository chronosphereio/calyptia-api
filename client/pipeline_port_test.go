//go:build integration
// +build integration

package client_test

import (
	"context"
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

	port, err := asUser.CreatePipelinePort(ctx, pipeline.ID, types.CreatePipelinePort{
		Protocol:     "tcp",
		FrontendPort: 4000,
		BackendPort:  4000,
		Endpoint:     "http://localhost:4000",
	})
	wantEqual(t, err, nil)

	got, err := asUser.PipelinePorts(ctx, pipeline.ID, types.PipelinePortsParams{})
	wantEqual(t, err, nil)

	wantEqual(t, len(got), 2) // The setup pipeline already contains ports in its config.

	wantEqual(t, got[0].ID, port.ID)
	wantEqual(t, got[0].Protocol, "tcp")
	wantEqual(t, got[0].FrontendPort, uint(4000))
	wantEqual(t, got[0].BackendPort, uint(4000))
	wantEqual(t, got[0].Endpoint, "http://localhost:4000")
	wantEqual(t, got[0].CreatedAt, port.CreatedAt)
	wantEqual(t, got[0].UpdatedAt, port.CreatedAt)

	wantNoEqual(t, got[1].ID, "")
	wantEqual(t, got[1].Protocol, "tcp")           // from the test config.
	wantEqual(t, got[1].FrontendPort, uint(24224)) // from the test config.
	wantEqual(t, got[1].BackendPort, uint(24224))  // from the test config.
	wantEqual(t, got[1].Endpoint, "")
	wantNoTimeZero(t, got[1].CreatedAt)
	wantNoTimeZero(t, got[1].UpdatedAt)
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
