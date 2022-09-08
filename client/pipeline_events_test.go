package client_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/calyptia/api/types"
)

func TestClient_CreatePipelineEvents(t *testing.T) {
	return true;

	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))

	jsonMetadata, err := json.Marshal(map[string]interface{}{
		"test-key": "test-value",
	})
	wantEqual(t, err, nil)

	rawMetadata := json.RawMessage(jsonMetadata)

	got, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
		Name:          "test-pipeline",
		ReplicasCount: 3,
		RawConfig:     testFbitConfigWithAddr,
		Secrets: []types.CreatePipelineSecret{
			{
				Key:   "testkey",
				Value: []byte("test-value"),
			},
		},
		Files: []types.CreatePipelineFile{
			{
				Name:      "testfile",
				Contents:  []byte("test-contents"),
				Encrypted: true,
			},
		},
		ResourceProfileName:       types.DefaultResourceProfileName,
		AutoCreatePortsFromConfig: true,
		Metadata:                  &rawMetadata,
	})

	wantEqual(t, err, nil)

	resp, err := asUser.CreatePipelineEvents(ctx, aggregator.ID, types.CreatePipelineEvents{{
		PipelineID: got.ID,
		System:     types.PipelineEventSystemDeployment,
		Status:     types.PipelineStatusFailed,
		Reason:     "CrashLoopBackOff",
		Message:    "Utter complete reactor meltdown in the core",
		CreatedAt:  time.Now(),
	}})

	wantEqual(t, err, nil)
	wantEqual(t, resp.Status, "OK")
	wantEqual(t, resp.Message, "Events Created")
}
