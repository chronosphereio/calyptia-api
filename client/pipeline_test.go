package client_test

import (
	"context"
	"encoding/json"
	"sort"
	"testing"

	"github.com/calyptia/api/types"
)

const testFbitConfigWithAddr = `
[INPUT]
	Name 			  forward
	Listen            0.0.0.0
	Port              24224
`

const testFbitConfigWithAddr3 = `
[INPUT]
	Name              forward
	Listen            0.0.0.0
	Port              24224
[INPUT]
	Name     syslog
	Listen   0.0.0.0
	Port     5140
	Mode     tcp
[INPUT]
	Name     syslog.1
	Listen   0.0.0.0
	Port     5140
	Mode     udp
`

func TestClient_CreatePipeline(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	withToken := withToken(t, asUser)

	aggregator, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
		Name: "test-aggregator",
	})
	wantEqual(t, err, nil)

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
	wantNoEqual(t, got.ID, "")
	wantEqual(t, got.Name, "test-pipeline")

	wantNoEqual(t, got.Config.ID, "")
	wantEqual(t, got.Config.RawConfig, testFbitConfigWithAddr)
	wantNoTimeZero(t, got.Config.CreatedAt)
	wantNoTimeZero(t, got.Config.UpdatedAt)

	wantEqual(t, len(got.Secrets), 1)
	wantNoEqual(t, got.Secrets[0].ID, "")
	wantEqual(t, got.Secrets[0].Key, "testkey")
	wantNoEqual(t, got.Secrets[0].Value, []byte("test-value")) // should be encrypted now.
	wantNoTimeZero(t, got.Secrets[0].CreatedAt)
	wantNoTimeZero(t, got.Secrets[0].UpdatedAt)

	wantEqual(t, len(got.Files), 2) // Aditional "parsers" file should be created by default.
	wantNoEqual(t, got.Files[0].ID, "")
	wantEqual(t, got.Files[0].Name, "testfile")
	wantNoEqual(t, got.Files[0].Contents, []byte("test-contents")) // should be encrypted now.
	wantNoTimeZero(t, got.Files[0].CreatedAt)
	wantNoTimeZero(t, got.Files[0].UpdatedAt)

	wantNoEqual(t, got.Status.ID, "")
	wantEqual(t, got.Status.Config, got.Config)
	wantEqual(t, got.Status.Status, types.PipelineStatusNew)
	wantNoTimeZero(t, got.Status.CreatedAt)
	wantNoTimeZero(t, got.Status.UpdatedAt)

	wantNoEqual(t, got.ResourceProfile.ID, "")
	wantEqual(t, got.ResourceProfile.Name, types.DefaultResourceProfileName)
	wantNoTimeZero(t, got.ResourceProfile.CreatedAt)
	wantNoTimeZero(t, got.ResourceProfile.UpdatedAt)

	wantEqual(t, got.ReplicasCount, uint64(3))
	wantNoTimeZero(t, got.CreatedAt)
}

func TestClient_Pipelines(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	withToken := withToken(t, asUser)

	aggregator, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
		Name:                    "test-aggregator",
		AddHealthCheckPipeline:  true,
		HealthCheckPipelinePort: 2020,
	})
	wantEqual(t, err, nil)

	jsonMetadata, err := json.Marshal(map[string]interface{}{
		"test-key": "test-value",
	})
	wantEqual(t, err, nil)

	rawMetadata := json.RawMessage(jsonMetadata)

	pipeline, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
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

	got, err := asUser.Pipelines(ctx, aggregator.ID, types.PipelinesParams{})
	wantEqual(t, err, nil)
	wantEqual(t, len(got), 2) // additional healthcheck pipeline should be created by default.

	wantEqual(t, got[0].ID, pipeline.ID)
	wantEqual(t, got[0].Name, pipeline.Name)
	wantEqual(t, got[0].Config, pipeline.Config)
	wantEqual(t, got[0].Status, pipeline.Status)
	wantEqual(t, got[0].ResourceProfile, pipeline.ResourceProfile)
	wantEqual(t, got[0].ReplicasCount, pipeline.ReplicasCount)
	wantNoEqual(t, got[0].Metadata, nil)
	wantEqual(t, *got[0].Metadata, rawMetadata)
	wantEqual(t, got[0].Config.CreatedAt, pipeline.Config.CreatedAt)
	wantEqual(t, got[0].Config.UpdatedAt, pipeline.Config.UpdatedAt)

	wantEqual(t, got[1], *aggregator.HealthCheckPipeline)
}

func TestClient_Pipeline(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	withToken := withToken(t, asUser)

	aggregator, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
		Name: "test-aggregator",
	})
	wantEqual(t, err, nil)

	jsonMetadata, err := json.Marshal(map[string]interface{}{
		"test-key": "test-value",
	})
	wantEqual(t, err, nil)

	rawMetadata := json.RawMessage(jsonMetadata)

	pipeline, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
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

	got, err := asUser.Pipeline(ctx, pipeline.ID)
	wantEqual(t, err, nil)

	wantEqual(t, got.ID, pipeline.ID)
	wantEqual(t, got.Name, pipeline.Name)
	wantEqual(t, got.Config, pipeline.Config)
	wantEqual(t, got.Status, pipeline.Status)
	wantEqual(t, got.ResourceProfile, pipeline.ResourceProfile)
	wantEqual(t, got.ReplicasCount, pipeline.ReplicasCount)
	wantNoEqual(t, got.Metadata, nil)
	wantEqual(t, *got.Metadata, rawMetadata)
	wantEqual(t, got.Config.CreatedAt, pipeline.Config.CreatedAt)
	wantEqual(t, got.Config.UpdatedAt, pipeline.Config.UpdatedAt)
}

func TestClient_UpdatePipeline(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	withToken := withToken(t, asUser)

	aggregator, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
		Name: "test-aggregator",
	})
	wantEqual(t, err, nil)

	jsonMetadata, err := json.Marshal(map[string]interface{}{
		"test-key": "test-value",
	})
	wantEqual(t, err, nil)

	rawMetadata := json.RawMessage(jsonMetadata)

	pipeline, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
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

	got, err := asUser.UpdatePipeline(ctx, pipeline.ID, types.UpdatePipeline{
		Name:                      ptrStr("test-pipeline-updated"),
		ReplicasCount:             ptrUint64(4),
		RawConfig:                 ptrStr(testFbitConfigWithAddr3),
		Status:                    (*types.PipelineStatusKind)(ptrStr(string(types.PipelineStatusStarted))),
		ResourceProfile:           ptrStr(string(types.ResourceProfileHighPerformanceOptimalThroughput)),
		AutoCreatePortsFromConfig: ptrBool(true),
	})
	wantEqual(t, err, nil)

	wantEqual(t, len(got.AddedPorts), 2) // new config has 2 additional ports tcp/5140 and udp/5141.

	// Sort added ports by protocold and port number
	// since they are created at the same time
	// and the order may switch.
	sort.Slice(got.AddedPorts, func(i, j int) bool {
		return got.AddedPorts[i].Protocol < got.AddedPorts[j].Protocol &&
			got.AddedPorts[i].FrontendPort < got.AddedPorts[j].FrontendPort &&
			got.AddedPorts[i].BackendPort < got.AddedPorts[j].BackendPort
	})

	wantNoEqual(t, got.AddedPorts[0].ID, "")
	wantEqual(t, got.AddedPorts[0].Protocol, "tcp")
	wantEqual(t, got.AddedPorts[0].FrontendPort, uint(5140))
	wantEqual(t, got.AddedPorts[0].BackendPort, uint(5140))
	wantEqual(t, got.AddedPorts[0].Endpoint, "")
	wantNoTimeZero(t, got.AddedPorts[0].CreatedAt)
	wantNoTimeZero(t, got.AddedPorts[0].UpdatedAt)

	wantNoEqual(t, got.AddedPorts[1].ID, "")
	wantEqual(t, got.AddedPorts[1].Protocol, "udp")
	wantEqual(t, got.AddedPorts[1].FrontendPort, uint(5140))
	wantEqual(t, got.AddedPorts[1].BackendPort, uint(5140))
	wantEqual(t, got.AddedPorts[1].Endpoint, "")
	wantNoTimeZero(t, got.AddedPorts[1].CreatedAt)
	wantNoTimeZero(t, got.AddedPorts[1].UpdatedAt)
}

func TestClient_DeletePipeline(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	withToken := withToken(t, asUser)

	aggregator, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
		Name: "test-aggregator",
	})
	wantEqual(t, err, nil)

	jsonMetadata, err := json.Marshal(map[string]interface{}{
		"test-key": "test-value",
	})
	wantEqual(t, err, nil)

	rawMetadata := json.RawMessage(jsonMetadata)

	pipeline, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
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

	err = asUser.DeletePipeline(ctx, pipeline.ID)
	wantEqual(t, err, nil)
}
