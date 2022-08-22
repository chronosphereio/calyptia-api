package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"testing"

	"github.com/calyptia/api/client"
	"github.com/calyptia/api/types"
)

const testFbitConfigWithAddr = `[INPUT]
	Name 			  forward
	Listen            0.0.0.0
	Port              24224`

const testFbitConfigWithAddr3 = `[INPUT]
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
	Mode     udp`

func TestClient_CreatePipeline(t *testing.T) {
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
	wantNoEqual(t, got.ID, "")
	wantEqual(t, got.Name, "test-pipeline")

	wantNoEqual(t, got.Config.ID, "")
	wantEqual(t, got.Config.RawConfig, testFbitConfigWithAddr)
	wantNoTimeZero(t, got.Config.CreatedAt)

	wantEqual(t, len(got.Secrets), 1)
	wantNoEqual(t, got.Secrets[0].ID, "")
	wantEqual(t, got.Secrets[0].Key, "testkey")
	wantNoEqual(t, got.Secrets[0].Value, []byte("test-value")) // should be encrypted now.
	wantNoTimeZero(t, got.Secrets[0].CreatedAt)
	wantNoTimeZero(t, got.Secrets[0].UpdatedAt)

	wantEqual(t, len(got.Files), 2) // Aditional "parsers" file should be created by default.

	// Sort files by name
	// since they are created at the same time
	// and the order may switch.
	sort.Slice(got.Files, func(i, j int) bool {
		return got.Files[i].Name < got.Files[j].Name
	})

	wantEqual(t, got.Files[0].Name, "parsers")

	wantNoEqual(t, got.Files[1].ID, "")
	wantEqual(t, got.Files[1].Name, "testfile")
	wantNoEqual(t, got.Files[1].Contents, []byte("test-contents")) // should be encrypted now.
	wantNoTimeZero(t, got.Files[1].CreatedAt)
	wantNoTimeZero(t, got.Files[1].UpdatedAt)

	wantNoEqual(t, got.Status.ID, "")
	wantEqual(t, got.Status.Config, got.Config)
	wantEqual(t, got.Status.Status, types.PipelineStatusNew)
	wantNoTimeZero(t, got.Status.CreatedAt)

	wantNoEqual(t, got.ResourceProfile.ID, "")
	wantEqual(t, got.ResourceProfile.Name, types.DefaultResourceProfileName)
	wantNoTimeZero(t, got.ResourceProfile.CreatedAt)
	wantNoTimeZero(t, got.ResourceProfile.UpdatedAt)

	wantEqual(t, got.ReplicasCount, uint(3))
	wantNoTimeZero(t, got.CreatedAt)
}

func TestClient_Pipelines(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	withToken := withToken(t, asUser)
	t.Run("ok", func(t *testing.T) {
		aggregator, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
			Name:                    "test-aggregator-one",
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
		wantEqual(t, len(got.Items), 2) // additional healthcheck pipeline should be created by default.

		wantEqual(t, got.Items[0].ID, pipeline.ID)
		wantEqual(t, got.Items[0].Name, pipeline.Name)
		wantEqual(t, got.Items[0].Config, pipeline.Config)
		wantEqual(t, got.Items[0].Status, pipeline.Status)
		wantEqual(t, got.Items[0].ResourceProfile, pipeline.ResourceProfile)
		wantEqual(t, got.Items[0].ReplicasCount, pipeline.ReplicasCount)
		wantNoEqual(t, got.Items[0].Metadata, nil)
		wantEqual(t, *got.Items[0].Metadata, rawMetadata)
		wantEqual(t, got.Items[0].Config.CreatedAt, pipeline.Config.CreatedAt)

		wantEqual(t, got.Items[1], *aggregator.HealthCheckPipeline)
	})

	t.Run("pagination", func(t *testing.T) {
		aggregator, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
			Name:                    "test-aggregator-two",
			AddHealthCheckPipeline:  true,
			HealthCheckPipelinePort: 2020,
		})
		wantEqual(t, err, nil)

		for i := 0; i < 10; i++ {
			_, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
				Name:      fmt.Sprintf("test-pipeline-%d", i),
				RawConfig: testFbitConfigWithAddr,
			})
			wantEqual(t, err, nil)
		}

		page1, err := asUser.Pipelines(ctx, aggregator.ID, types.PipelinesParams{
			Last: ptrUint(3),
		})
		wantEqual(t, err, nil)
		wantEqual(t, len(page1.Items), 3)
		wantNoEqual(t, page1.EndCursor, (*string)(nil))

		page2, err := asUser.Pipelines(ctx, aggregator.ID, types.PipelinesParams{
			Last:   ptrUint(3),
			Before: page1.EndCursor,
		})
		wantEqual(t, err, nil)
		wantEqual(t, len(page2.Items), 3)
		wantNoEqual(t, page2.EndCursor, (*string)(nil))

		wantNoEqual(t, page1.Items, page2.Items)
		wantNoEqual(t, *page1.EndCursor, *page2.EndCursor)
		wantEqual(t, page1.Items[2].CreatedAt.After(page2.Items[0].CreatedAt), true)
	})

	//nolint:dupl // this is a test and could be duplicated
	t.Run("tags", func(t *testing.T) {
		aggregator, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
			Name: "test-aggregator-three",
		})
		wantEqual(t, err, nil)
		for i := 10; i < 20; i++ {
			pipeline := types.CreatePipeline{
				Name:      fmt.Sprintf("test-pipeline-%d", i),
				RawConfig: testFbitConfigWithAddr,
			}
			if i >= 15 {
				pipeline.Tags = append(pipeline.Tags, "tagone", "tagthree")
			} else {
				pipeline.Tags = append(pipeline.Tags, "tagtwo", "tagthree")
			}
			_, err := asUser.CreatePipeline(ctx, aggregator.ID, pipeline)
			wantEqual(t, err, nil)
		}

		opts := types.PipelinesParams{}
		s := "tagone"
		opts.Tags = &s
		tag1, err := asUser.Pipelines(ctx, aggregator.ID, opts)
		wantEqual(t, err, nil)
		wantEqual(t, len(tag1.Items), 5)
		wantEqual(t, tag1.Items[0].Tags, []string{"tagone", "tagthree"})

		s2 := "tagtwo"
		opts.Tags = &s2
		tag2, err := asUser.Pipelines(ctx, aggregator.ID, opts)
		wantEqual(t, err, nil)
		wantEqual(t, len(tag2.Items), 5)
		wantEqual(t, tag2.Items[0].Tags, []string{"tagtwo", "tagthree"})

		s3 := "tagthree"
		opts.Tags = &s3
		tag3, err := asUser.Pipelines(ctx, aggregator.ID, opts)
		wantEqual(t, err, nil)
		wantEqual(t, len(tag3.Items), 10)
		wantEqual(t, tag3.Items[0].Tags, []string{"tagone", "tagthree"})

		s4 := "tagone AND tagtwo"
		opts.Tags = &s4
		tag4, err := asUser.Pipelines(ctx, aggregator.ID, opts)
		wantEqual(t, err, nil)
		wantEqual(t, len(tag4.Items), 0)
	})
}

func TestClient_ProjectPipelines(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	withToken := withToken(t, asUser)

	var want []types.CreatedPipeline

	for i := 0; i < 3; i++ {
		aggregator, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
			Name: fmt.Sprintf("test-aggregator-%d", i),
		})
		wantEqual(t, err, nil)

		pipeline, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
			Name:                      fmt.Sprintf("test-pipeline-%d", i),
			ReplicasCount:             1,
			RawConfig:                 testFbitConfigWithAddr,
			ResourceProfileName:       types.DefaultResourceProfileName,
			AutoCreatePortsFromConfig: true,
		})
		wantEqual(t, err, nil)

		want = append(want, pipeline)
	}

	project := defaultProject(t, asUser)

	got, err := asUser.ProjectPipelines(ctx, project.ID, types.PipelinesParams{})
	wantEqual(t, err, nil)
	wantEqual(t, len(got.Items), len(want))

	// Reverse order to do proper assertion since retuls come in descending order.
	sort.Slice(got.Items, func(i, j int) bool {
		return got.Items[i].CreatedAt.Before(got.Items[j].CreatedAt)
	})

	for i := range got.Items {
		wantEqual(t, got.Items[i].ID, want[i].ID)
		wantEqual(t, got.Items[i].Name, want[i].Name)
		wantEqual(t, got.Items[i].Config, want[i].Config)
		wantEqual(t, got.Items[i].Status, want[i].Status)
		wantEqual(t, got.Items[i].ResourceProfile, want[i].ResourceProfile)
		wantEqual(t, got.Items[i].ReplicasCount, want[i].ReplicasCount)
		wantEqual(t, got.Items[i].CreatedAt, want[i].CreatedAt)
	}

	t.Run("pagination", func(t *testing.T) {
		aggregator, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
			Name: "test-aggregator",
		})
		wantEqual(t, err, nil)
		for i := 0; i < 9; i++ {
			_, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
				Name:      fmt.Sprintf("test-pipeline-%d", i),
				RawConfig: testFbitConfigWithAddr,
			})
			wantEqual(t, err, nil)
		}

		page1, err := asUser.Pipelines(ctx, aggregator.ID, types.PipelinesParams{
			Last: ptrUint(3),
		})
		wantEqual(t, err, nil)
		wantEqual(t, len(page1.Items), 3)
		wantNoEqual(t, page1.EndCursor, (*string)(nil))

		page2, err := asUser.Pipelines(ctx, aggregator.ID, types.PipelinesParams{
			Last:   ptrUint(3),
			Before: page1.EndCursor,
		})
		wantEqual(t, err, nil)
		wantEqual(t, len(page2.Items), 3)
		wantNoEqual(t, page2.EndCursor, (*string)(nil))

		wantNoEqual(t, page1.Items, page2.Items)
		wantNoEqual(t, *page1.EndCursor, *page2.EndCursor)
		wantEqual(t, page1.Items[2].CreatedAt.After(page2.Items[0].CreatedAt), true)
	})

	//nolint:dupl // this is a test and could be duplicated
	t.Run("tags", func(t *testing.T) {
		aggregator, err := withToken.CreateAggregator(ctx, types.CreateAggregator{
			Name: "test-aggregator-one",
		})
		wantEqual(t, err, nil)
		for i := 10; i < 20; i++ {
			pipeline := types.CreatePipeline{
				Name:      fmt.Sprintf("test-pipeline-%d", i),
				RawConfig: testFbitConfigWithAddr,
			}
			if i >= 15 {
				pipeline.Tags = append(pipeline.Tags, "tagone", "tagthree")
			} else {
				pipeline.Tags = append(pipeline.Tags, "tagtwo", "tagthree")
			}
			_, err := asUser.CreatePipeline(ctx, aggregator.ID, pipeline)
			wantEqual(t, err, nil)
		}

		opts := types.PipelinesParams{}
		s := "tagone"
		opts.Tags = &s
		tag1, err := asUser.ProjectPipelines(ctx, project.ID, opts)
		wantEqual(t, err, nil)
		wantEqual(t, len(tag1.Items), 5)
		wantEqual(t, tag1.Items[0].Tags, []string{"tagone", "tagthree"})

		s2 := "tagtwo"
		opts.Tags = &s2
		tag2, err := asUser.ProjectPipelines(ctx, project.ID, opts)
		wantEqual(t, err, nil)
		wantEqual(t, len(tag2.Items), 5)
		wantEqual(t, tag2.Items[0].Tags, []string{"tagtwo", "tagthree"})

		s3 := "tagthree"
		opts.Tags = &s3
		tag3, err := asUser.ProjectPipelines(ctx, project.ID, opts)
		wantEqual(t, err, nil)
		wantEqual(t, len(tag3.Items), 10)
		wantEqual(t, tag3.Items[0].Tags, []string{"tagone", "tagthree"})

		s4 := "tagone AND tagtwo"
		opts.Tags = &s4
		tag4, err := asUser.ProjectPipelines(ctx, project.ID, opts)
		wantEqual(t, err, nil)
		wantEqual(t, len(tag4.Items), 0)
	})
}

func TestClient_Pipeline(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))

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

	got, err := asUser.Pipeline(ctx, pipeline.ID, types.PipelineParams{})
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
}

func TestClient_UpdatePipeline(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))

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
		Name:          ptrStr("test-pipeline-updated"),
		ReplicasCount: ptrUint(4),
		RawConfig:     ptrStr(testFbitConfigWithAddr3),
		Secrets: []types.UpdatePipelineSecret{
			{
				Key:   ptrStr("testkeyupdated"),
				Value: ptrBytes([]byte("test-value-updated")),
			},
		},
		Files: []types.UpdatePipelineFile{
			{
				Name:      ptrStr("testfileupdated"),
				Contents:  ptrBytes([]byte("test-contents-updated")),
				Encrypted: ptrBool(true),
			},
		},
		Status:                    (*types.PipelineStatusKind)(ptrStr(string(types.PipelineStatusStarted))),
		ResourceProfile:           ptrStr(string(types.ResourceProfileHighPerformanceOptimalThroughput)),
		AutoCreatePortsFromConfig: ptrBool(true),
	})
	wantEqual(t, err, nil)

	wantEqual(t, len(got.AddedPorts), 2) // new config has 2 additional ports tcp/5140 and udp/5141.
	wantEqual(t, len(got.RemovedPorts), 0)

	// Sort added ports by protocold and port number
	// since they are created at the same time
	// and the order may switch.
	sort.Slice(got.AddedPorts, func(i, j int) bool {
		if got.AddedPorts[i].Protocol == got.AddedPorts[j].Protocol {
			if got.AddedPorts[i].FrontendPort == got.AddedPorts[j].FrontendPort {
				return got.AddedPorts[i].BackendPort < got.AddedPorts[j].BackendPort
			}
			return got.AddedPorts[i].FrontendPort < got.AddedPorts[j].FrontendPort
		}
		return got.AddedPorts[i].Protocol < got.AddedPorts[j].Protocol
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
	aggregator := setupAggregator(t, withToken(t, asUser))

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

func setupPipeline(t *testing.T, asUser *client.Client, aggregatorID string) types.CreatedPipeline {
	t.Helper()
	ctx := context.Background()

	pipeline, err := asUser.CreatePipeline(ctx, aggregatorID, types.CreatePipeline{
		Name:                      "test-pipeline",
		ReplicasCount:             1,
		RawConfig:                 testFbitConfigWithAddr,
		ResourceProfileName:       types.DefaultResourceProfileName,
		AutoCreatePortsFromConfig: true,
	})
	wantEqual(t, err, nil)

	return pipeline
}
