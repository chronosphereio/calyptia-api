package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	math_rand "math/rand"
	"sort"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"

	"github.com/calyptia/api/client"
	"github.com/calyptia/api/types"
)

const testFluentBitConfigWithInvalidOutput = `[OUTPUT]
	Name tcp
	Port 5556`

const testFbitConfigWithMultipleOutputs = `[INPUT]
	Name 			  forward
	Host              0.0.0.0
	Port              24224
[OUTPUT]
	Name tcp
	Host 127.0.0.1
	Port 5556
[OUTPUT]
	Name opensearch
	Host api.elasticsearch.com
[OUTPUT]
	Name syslog
	Host 192.168.0.1
	Mode udp
[OUTPUT]
	Name syslog
	Host 192.168.0.1
	Mode tcp
[OUTPUT]
	Name gelf
	Host 127.0.0.1
	Port 3333`

const testFbitConfigWithSingleOutput = `[INPUT]
	Name 			  forward
	Host              0.0.0.0
	Port              24224
[OUTPUT]
	Name syslog
	Host 192.168.0.1
	Mode udp`

const testFbitConfigWithAddr = `[INPUT]
	Name 			  forward
	Host              0.0.0.0
	Port              24224`

const testFbitConfigWithAddr3 = `[INPUT]
	Name              forward
	Host              0.0.0.0
	Port              24224
[INPUT]
	Name     syslog
	Host     0.0.0.0
	Port     5140
	Mode     tcp
[INPUT]
	Name     syslog
	Host     0.0.0.0
	Port     5140
	Mode     udp`

func TestClient_CreatePipeline(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))

	t.Run("ok", func(t *testing.T) {
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
	})

	t.Run("ok valid pipeline kind", func(t *testing.T) {
		got, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
			Name:                      "test-pipeline-01",
			Kind:                      types.PipelineKindDeployment,
			ReplicasCount:             1,
			RawConfig:                 testFbitConfigWithAddr,
			AutoCreatePortsFromConfig: true,
		})

		wantEqual(t, err, nil)
		wantEqual(t, got.Kind, types.PipelineKindDeployment)
	})

	t.Run("ok valid pipeline kind", func(t *testing.T) {
		got, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
			Name:                      "test-pipeline-02",
			ReplicasCount:             1,
			RawConfig:                 testFbitConfigWithAddr,
			AutoCreatePortsFromConfig: true,
		})

		wantEqual(t, err, nil)
		wantEqual(t, got.Kind, types.PipelineKindDeployment)
	})

	t.Run("ok valid pipeline kind daemonSet", func(t *testing.T) {
		got, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
			Name:      "test-pipeline-03",
			Kind:      types.PipelineKindDaemonSet,
			RawConfig: testFbitConfigWithAddr,
		})

		wantEqual(t, err, nil)
		wantEqual(t, got.Kind, types.PipelineKindDaemonSet)
	})

	t.Run("not ok pipeline kind", func(t *testing.T) {
		_, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
			Name:          "test-pipeline-04",
			Kind:          "invalid",
			ReplicasCount: 1,
			RawConfig:     testFbitConfigWithAddr,
		})

		wantNoEqual(t, err, nil)
	})

	t.Run("create with auto create checks with invalid config", func(t *testing.T) {
		_, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
			Name:                       "test-pipeline-09",
			Kind:                       types.PipelineKindDeployment,
			RawConfig:                  testFluentBitConfigWithInvalidOutput,
			AutoCreateChecksFromConfig: true,
		})

		wantNoEqual(t, err, nil)
		wantErrMsg(t, err, "invalid pipeline config: host property not defined in section \"tcp.0\"")
	})

	t.Run("create with auto create checks enabled and multi output configuration", func(t *testing.T) {
		got, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
			Name:                       "test-pipeline-10",
			Kind:                       types.PipelineKindDeployment,
			RawConfig:                  testFbitConfigWithMultipleOutputs,
			AutoCreateChecksFromConfig: true,
		})

		wantEqual(t, err, nil)
		wantEqual(t, len(got.Checks), 5)
	})
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

	t.Run("render_config_sections", func(t *testing.T) {
		agg, err := withToken.CreateAggregator(ctx, types.CreateAggregator{})
		assert.NoError(t, err)

		createdPip, err := asUser.CreatePipeline(ctx, agg.ID, types.CreatePipeline{
			RawConfig: newClassicConf(`
				[INPUT]
					Name dummy
			`),
		})
		assert.NoError(t, err)

		defaultProj := defaultProject(t, asUser)

		csInput := types.CreateConfigSection{
			Kind: types.SectionKindOutput,
			Properties: types.Pairs{
				{Key: "Name", Value: "http"},
				{Key: "Match", Value: "*"},
				{Key: "host", Value: "localhost"},
				{Key: "port", Value: "80"},
			},
		}
		cs, err := asUser.CreateConfigSection(ctx, defaultProj.ID, csInput)
		assert.NoError(t, err)

		err = asUser.UpdateConfigSectionSet(ctx, createdPip.ID, cs.ID)
		assert.NoError(t, err)

		pp, err := asUser.Pipelines(ctx, agg.ID, types.PipelinesParams{
			RenderWithConfigSections: true,
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(pp.Items))

		pip := pp.Items[0]

		assert.Equal(t, 1, len(pip.ConfigSections))
		assert.Equal(t, types.ConfigSection{
			ID:         cs.ID,
			Kind:       csInput.Kind,
			ProjectID:  defaultProj.ID,
			Properties: csInput.Properties,
			CreatedAt:  cs.CreatedAt,
			UpdatedAt:  cs.CreatedAt,
		}, pip.ConfigSections[0])

		assert.Equal(t, newClassicConf(`
			[INPUT]
			    Name dummy
			[OUTPUT]
			    Name http
			    Match *
			    host localhost
			    port 80
		`), pip.Config.RawConfig)
	})
	t.Run("auto create checks from config", func(t *testing.T) {
		agg, err := withToken.CreateAggregator(ctx, types.CreateAggregator{})
		assert.NoError(t, err)

		createdPip, err := asUser.CreatePipeline(ctx, agg.ID, types.CreatePipeline{
			RawConfig:                  testFbitConfigWithMultipleOutputs,
			AutoCreateChecksFromConfig: true,
		})
		assert.NoError(t, err)
		assert.Equal(t, len(createdPip.Checks), 5)

		pp, err := asUser.Pipelines(ctx, agg.ID, types.PipelinesParams{})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(pp.Items))
		assert.Equal(t, pp.Items[0].ChecksOK, 0)
		assert.Equal(t, pp.Items[0].ChecksTotal, 5)
	})

	t.Run("update await_for_checks property", func(t *testing.T) {
		agg, err := withToken.CreateAggregator(ctx, types.CreateAggregator{})
		assert.NoError(t, err)
		createdPip, err := asUser.CreatePipeline(ctx, agg.ID, types.CreatePipeline{
			RawConfig: testFbitConfigWithAddr,
		})
		assert.NoError(t, err)
		assert.Equal(t, createdPip.WaitForChecksBeforeDeploying, false)

		_, err = asUser.UpdatePipeline(ctx, createdPip.ID, types.UpdatePipeline{
			WaitForChecksBeforeDeploying: ptr(true),
		})

		assert.NoError(t, err)

		updatedPip, err := asUser.Pipeline(ctx, createdPip.ID, types.PipelineParams{})
		assert.NoError(t, err)
		assert.Equal(t, updatedPip.WaitForChecksBeforeDeploying, true)

		createdPip, err = asUser.CreatePipeline(ctx, agg.ID, types.CreatePipeline{
			RawConfig:                    testFbitConfigWithAddr,
			WaitForChecksBeforeDeploying: true,
		})
		assert.NoError(t, err)
		assert.Equal(t, createdPip.WaitForChecksBeforeDeploying, true)

		_, err = asUser.UpdatePipeline(ctx, createdPip.ID, types.UpdatePipeline{
			WaitForChecksBeforeDeploying: ptr(false),
		})

		assert.NoError(t, err)

		updatedPip, err = asUser.Pipeline(ctx, createdPip.ID, types.PipelineParams{})
		assert.NoError(t, err)
		assert.Equal(t, updatedPip.WaitForChecksBeforeDeploying, false)
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

	createdPip, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
		Name:          "test-pipeline",
		ReplicasCount: 3,
		RawConfig: newClassicConf(`
			[INPUT]
				Name dummy
		`),
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

	pip, err := asUser.Pipeline(ctx, createdPip.ID, types.PipelineParams{})
	wantEqual(t, err, nil)

	wantEqual(t, pip.ID, createdPip.ID)
	wantEqual(t, pip.Name, createdPip.Name)
	wantEqual(t, pip.Config, createdPip.Config)
	wantEqual(t, pip.Status, createdPip.Status)
	wantEqual(t, pip.ResourceProfile, createdPip.ResourceProfile)
	wantEqual(t, pip.ReplicasCount, createdPip.ReplicasCount)
	wantNoEqual(t, pip.Metadata, nil)
	wantEqual(t, *pip.Metadata, rawMetadata)
	wantEqual(t, pip.Config.CreatedAt, createdPip.Config.CreatedAt)

	t.Run("render_config_sections", func(t *testing.T) {
		defaultProj := defaultProject(t, asUser)

		csInput := types.CreateConfigSection{
			Kind: types.SectionKindOutput,
			Properties: types.Pairs{
				{Key: "Name", Value: "http"},
				{Key: "Match", Value: "*"},
				{Key: "host", Value: "localhost"},
				{Key: "port", Value: "80"},
			},
		}
		cs, err := asUser.CreateConfigSection(ctx, defaultProj.ID, csInput)
		assert.NoError(t, err)

		err = asUser.UpdateConfigSectionSet(ctx, createdPip.ID, cs.ID)
		assert.NoError(t, err)

		pip, err := asUser.Pipeline(ctx, createdPip.ID, types.PipelineParams{
			RenderWithConfigSections: true,
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(pip.ConfigSections))
		assert.Equal(t, types.ConfigSection{
			ID:         cs.ID,
			Kind:       csInput.Kind,
			ProjectID:  defaultProj.ID,
			Properties: csInput.Properties,
			CreatedAt:  cs.CreatedAt,
			UpdatedAt:  cs.CreatedAt,
		}, pip.ConfigSections[0])

		assert.Equal(t, newClassicConf(`
			[INPUT]
			    Name dummy
			[OUTPUT]
			    Name http
			    Match *
			    host localhost
			    port 80
		`), pip.Config.RawConfig)
	})
}

func TestClient_UpdatePipeline(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))

	t.Run("ok", func(t *testing.T) {
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
			Name:          ptr("test-pipeline-updated"),
			ReplicasCount: ptrUint(4),
			RawConfig:     ptr(testFbitConfigWithAddr3),
			Secrets: []types.UpdatePipelineSecret{
				{
					Key:   ptr("testkeyupdated"),
					Value: ptr([]byte("test-value-updated")),
				},
			},
			Files: []types.UpdatePipelineFile{
				{
					Name:      ptr("testfileupdated"),
					Contents:  ptr([]byte("test-contents-updated")),
					Encrypted: ptr(true),
				},
			},
			Status:                    (*types.PipelineStatusKind)(ptr(string(types.PipelineStatusStarted))),
			ResourceProfile:           ptr(types.ResourceProfileHighPerformanceOptimalThroughput),
			AutoCreatePortsFromConfig: ptr(true),
		})
		wantEqual(t, err, nil)
		wantEqual(t, len(got.AddedPorts), 2) // new config has 2 additional ports tcp/5140 and udp/5141.
		wantEqual(t, len(got.RemovedPorts), 0)

		// Sort added ports by protocol and port number
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
	})
	t.Run("not ok pipeline kind", func(t *testing.T) {
		got, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
			Name:      "test-pipeline-03",
			Kind:      types.PipelineKindDaemonSet,
			RawConfig: testFbitConfigWithAddr,
		})

		wantEqual(t, err, nil)
		wantEqual(t, got.Kind, types.PipelineKindDaemonSet)

		pipelineKind := types.PipelineKindDeployment
		_, err = asUser.UpdatePipeline(ctx, got.ID, types.UpdatePipeline{
			Kind: &pipelineKind,
		})
		wantNoEqual(t, err, nil)
		wantErrMsg(t, err, "pipeline kind cannot be updated")
	})
	t.Run("not ok pipeline replicaSet for daemonSet", func(t *testing.T) {
		got, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
			Name:      "test-pipeline-04",
			Kind:      types.PipelineKindDaemonSet,
			RawConfig: testFbitConfigWithAddr,
		})

		wantEqual(t, err, nil)
		wantEqual(t, got.Kind, types.PipelineKindDaemonSet)

		var replicaSet uint = 3
		_, err = asUser.UpdatePipeline(ctx, got.ID, types.UpdatePipeline{
			ReplicasCount: &replicaSet,
		})
		wantNoEqual(t, err, nil)
		wantErrMsg(t, err, "pipeline replicas can only be set for pipelines of kind deployment")
	})
	t.Run("null status events", func(t *testing.T) {
		got, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
			Name:      "test-pipeline-05",
			Kind:      types.PipelineKindDeployment,
			RawConfig: testFbitConfigWithAddr,
		})
		wantEqual(t, err, nil)

		_, err = asUser.UpdatePipeline(ctx, got.ID, types.UpdatePipeline{
			Events: nil,
		})
		wantEqual(t, err, nil)

		pip, err := asUser.Pipeline(ctx, got.ID, types.PipelineParams{})
		wantEqual(t, err, nil)
		wantEqual(t, len(pip.Status.Events), 0)
	})
	t.Run("empty status events", func(t *testing.T) {
		got, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
			Name:      "test-pipeline-06",
			Kind:      types.PipelineKindDeployment,
			RawConfig: testFbitConfigWithAddr,
		})
		wantEqual(t, err, nil)

		_, err = asUser.UpdatePipeline(ctx, got.ID, types.UpdatePipeline{
			Events: []types.PipelineEvent{},
		})
		wantEqual(t, err, nil)

		pip, err := asUser.Pipeline(ctx, got.ID, types.PipelineParams{})
		wantEqual(t, err, nil)
		wantEqual(t, len(pip.Status.Events), 0)
	})
	t.Run("single status event", func(t *testing.T) {
		got, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
			Name:      "test-pipeline-07",
			Kind:      types.PipelineKindDeployment,
			RawConfig: testFbitConfigWithAddr,
		})
		wantEqual(t, err, nil)

		_, err = asUser.UpdatePipeline(ctx, got.ID, types.UpdatePipeline{
			Events: []types.PipelineEvent{
				{
					Source:   "k8s:pod",
					Reason:   "event",
					Message:  "message",
					LoggedAt: time.Now(),
				},
			},
		})

		wantEqual(t, err, nil)
		pip, err := asUser.Pipeline(ctx, got.ID, types.PipelineParams{})
		wantEqual(t, err, nil)
		wantEqual(t, len(pip.Status.Events), 1)

		_, err = asUser.UpdatePipeline(ctx, got.ID, types.UpdatePipeline{
			ReplicasCount: ptrUint(2),
		})

		wantEqual(t, err, nil)

		pip, err = asUser.Pipeline(ctx, got.ID, types.PipelineParams{})
		wantEqual(t, err, nil)
		wantEqual(t, len(pip.Status.Events), 1)
	})
	t.Run("multi status event", func(t *testing.T) {
		got, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
			Name:      "test-pipeline-08",
			Kind:      types.PipelineKindDeployment,
			RawConfig: testFbitConfigWithAddr,
		})
		wantEqual(t, err, nil)

		events := []types.PipelineEvent{
			{
				Source:   types.PipelineEventSourcePod,
				Reason:   "event",
				Message:  "message",
				LoggedAt: time.Now(),
			},
			{
				Source:   types.PipelineEventSourcePod,
				Reason:   "event",
				Message:  "message",
				LoggedAt: time.Now(),
			},
			{
				Source:   types.PipelineEventSourcePod,
				Reason:   "event",
				Message:  "message",
				LoggedAt: time.Now(),
			},
			{
				Source:   types.PipelineEventSourcePod,
				Reason:   "event",
				Message:  "message",
				LoggedAt: time.Now(),
			},
		}

		_, err = asUser.UpdatePipeline(ctx, got.ID, types.UpdatePipeline{
			Events: events,
		})

		wantEqual(t, err, nil)
		pip, err := asUser.Pipeline(ctx, got.ID, types.PipelineParams{})
		wantEqual(t, err, nil)
		wantEqual(t, len(pip.Status.Events), 4)

		_, err = asUser.UpdatePipeline(ctx, got.ID, types.UpdatePipeline{
			ReplicasCount: ptrUint(2),
		})

		wantEqual(t, err, nil)

		pip, err = asUser.Pipeline(ctx, got.ID, types.PipelineParams{})
		wantEqual(t, err, nil)
		wantEqual(t, len(pip.Status.Events), 4)
		wantEqual(t, pip.Status.Events[0].Source, events[0].Source)
		wantEqual(t, pip.Status.Events[0].Reason, events[0].Reason)
		wantEqual(t, pip.Status.Events[0].Message, events[0].Message)

		// safety check for sequential loads
		_, err = asUser.Pipelines(ctx, aggregator.ID, types.PipelinesParams{})
		wantEqual(t, err, nil)
	})
	t.Run("update with auto create checks with invalid config", func(t *testing.T) {
		_, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
			Name:                       "test-pipeline-09",
			Kind:                       types.PipelineKindDeployment,
			RawConfig:                  testFluentBitConfigWithInvalidOutput,
			AutoCreateChecksFromConfig: true,
		})

		wantNoEqual(t, err, nil)
		wantErrMsg(t, err, "invalid pipeline config: host property not defined in section \"tcp.0\"")
	})

	t.Run("update with auto create checks enabled and multi output configuration", func(t *testing.T) {
		got, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
			Name:                       "test-pipeline-10",
			Kind:                       types.PipelineKindDeployment,
			RawConfig:                  testFbitConfigWithSingleOutput,
			AutoCreateChecksFromConfig: true,
		})

		wantEqual(t, err, nil)
		wantEqual(t, len(got.Checks), 1)

		updated, err := asUser.UpdatePipeline(ctx, got.ID, types.UpdatePipeline{
			RawConfig:                  ptr(testFbitConfigWithMultipleOutputs),
			AutoCreateChecksFromConfig: ptr(true),
		})

		wantEqual(t, err, nil)
		wantEqual(t, len(updated.AddedChecks), 5)
		wantEqual(t, len(updated.RemovedChecks), 1)
	})

	t.Run("update with auto create checks enabled from disabled", func(t *testing.T) {
		got, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
			Name:                       "test-pipeline-11",
			Kind:                       types.PipelineKindDeployment,
			RawConfig:                  testFbitConfigWithSingleOutput,
			AutoCreateChecksFromConfig: false,
		})

		wantEqual(t, err, nil)
		wantEqual(t, len(got.Checks), 0)

		updated, err := asUser.UpdatePipeline(ctx, got.ID, types.UpdatePipeline{
			RawConfig:                  ptr(testFbitConfigWithMultipleOutputs),
			AutoCreateChecksFromConfig: ptr(true),
		})

		wantEqual(t, err, nil)
		wantEqual(t, len(updated.AddedChecks), 5)
		wantEqual(t, len(updated.RemovedChecks), 0)
	})
}

func TestClient_DeletePipeline(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))

	t.Run("ok", func(t *testing.T) {
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
	})
}

func TestClient_DeletePipelines(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	agg := setupAggregator(t, withToken(t, asUser))

	pip1, err := asUser.CreatePipeline(ctx, agg.ID, types.CreatePipeline{
		Name:      "test-pipeline-1",
		RawConfig: testFbitConfigWithAddr,
	})
	wantEqual(t, err, nil)

	pip2, err := asUser.CreatePipeline(ctx, agg.ID, types.CreatePipeline{
		Name:      "test-pipeline-2",
		RawConfig: testFbitConfigWithAddr,
	})
	wantEqual(t, err, nil)

	defer func() {
		err := asUser.DeletePipeline(ctx, pip2.ID)
		wantEqual(t, err, nil)
	}()

	err = asUser.DeletePipelines(ctx, agg.ID, pip1.ID)
	wantEqual(t, err, nil)

	_, err = asUser.Pipeline(ctx, pip1.ID, types.PipelineParams{})
	wantErrMsg(t, err, "pipeline not found")

	_, err = asUser.Pipeline(ctx, pip2.ID, types.PipelineParams{})
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

func randPort() uint {
	return uint(randInt(1, 65535))
}

func randInt(min, max int) int {
	math_rand.Seed(time.Now().UnixNano())
	return min + math_rand.Intn(max-min+1)
}
