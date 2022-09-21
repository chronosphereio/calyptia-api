package client_test

import (
	"context"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/calyptia/api/types"
)

func TestClient_CreateConfigSection(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)
	proj := defaultProject(t, asUser)

	t.Run("invalid_project_id", func(t *testing.T) {
		_, err := asUser.CreateConfigSection(ctx, "@nope@", types.CreateConfigSection{})
		assert.EqualError(t, err, "invalid project ID")
	})

	t.Run("permission_denied", func(t *testing.T) {
		_, err := asUser.CreateConfigSection(ctx, randUUID(t), types.CreateConfigSection{})
		assert.EqualError(t, err, "permission denied")
	})

	t.Run("unsuported_config_section", func(t *testing.T) {
		_, err := asUser.CreateConfigSection(ctx, proj.ID, types.CreateConfigSection{
			Kind: types.ConfigSectionKind("@nope@"),
		})
		assert.EqualError(t, err, "unsupported config section")
	})

	t.Run("missing_output_name", func(t *testing.T) {
		_, err := asUser.CreateConfigSection(ctx, proj.ID, types.CreateConfigSection{
			Kind: types.SectionKindOutput,
		})
		// TODO: differentita between missing `Name` property
		// and unknown plugin name.
		assert.EqualError(t, err, `invalid pipeline config: unknown plugin ""`)
	})

	t.Run("unknown_output", func(t *testing.T) {
		_, err := asUser.CreateConfigSection(ctx, proj.ID, types.CreateConfigSection{
			Kind: types.SectionKindOutput,
			Properties: types.Pairs{
				{Key: "name", Value: "@nope@"},
			},
		})
		assert.EqualError(t, err, `invalid pipeline config: unknown plugin "@nope@"`)
	})

	t.Run("invalid_property", func(t *testing.T) {
		_, err := asUser.CreateConfigSection(ctx, proj.ID, types.CreateConfigSection{
			Kind: types.SectionKindOutput,
			Properties: types.Pairs{
				{Key: "name", Value: "http"},
				{Key: "port", Value: "@nope@"},
			},
		})
		assert.EqualError(t, err, `invalid pipeline config: http: expected "port" to be a valid integer; got "@nope@"`)
	})

	t.Run("ok", func(t *testing.T) {
		got, err := asUser.CreateConfigSection(ctx, proj.ID, types.CreateConfigSection{
			Kind: types.SectionKindOutput,
			Properties: types.Pairs{
				{Key: "name", Value: "stdout"},
				{Key: "match", Value: "*"},
			},
		})
		assert.NoError(t, err)
		assert.NotZero(t, got)
	})
}

func TestClient_ConfigSections(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)
	proj := defaultProject(t, asUser)

	t.Run("invalid_project_id", func(t *testing.T) {
		_, err := asUser.ConfigSections(ctx, "@nope@", types.ConfigSectionsParams{})
		assert.EqualError(t, err, "invalid project ID")
	})

	t.Run("permission_denied", func(t *testing.T) {
		_, err := asUser.ConfigSections(ctx, randUUID(t), types.ConfigSectionsParams{})
		assert.EqualError(t, err, "permission denied")
	})

	t.Run("empty", func(t *testing.T) {
		got, err := asUser.ConfigSections(ctx, proj.ID, types.ConfigSectionsParams{})
		assert.NoError(t, err)
		assert.Equal(t, []types.ConfigSection{}, got.Items)
		assert.Zero(t, got.EndCursor)
	})

	t.Run("ok", func(t *testing.T) {
		in := types.CreateConfigSection{
			Kind: types.SectionKindOutput,
			Properties: types.Pairs{
				{Key: "name", Value: "http"},
				{Key: "port", Value: 80},
			},
		}

		created, err := asUser.CreateConfigSection(ctx, proj.ID, in)
		assert.NoError(t, err)

		got, err := asUser.ConfigSections(ctx, proj.ID, types.ConfigSectionsParams{})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(got.Items))
		assert.Equal(t, 2, len(got.Items[0].Properties))
		assert.NotZero(t, got.EndCursor)

		assert.Equal(t, types.ConfigSection{
			ID:         created.ID,
			ProjectID:  proj.ID,
			Kind:       in.Kind,
			Properties: in.Properties,
			CreatedAt:  created.CreatedAt,
			UpdatedAt:  created.CreatedAt,
		}, got.Items[0])
	})
}

func TestClient_ConfigSection(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)
	proj := defaultProject(t, asUser)

	t.Run("invalid_config_section_id", func(t *testing.T) {
		_, err := asUser.ConfigSection(ctx, "@nope@")
		assert.EqualError(t, err, "invalid config section ID")
	})

	t.Run("config_section_not_found", func(t *testing.T) {
		_, err := asUser.ConfigSection(ctx, randUUID(t))
		assert.EqualError(t, err, "config section not found")
	})

	t.Run("ok", func(t *testing.T) {
		in := types.CreateConfigSection{
			Kind: types.SectionKindOutput,
			Properties: types.Pairs{
				{Key: "name", Value: "http"},
				{Key: "port", Value: 80},
			},
		}

		created, err := asUser.CreateConfigSection(ctx, proj.ID, in)
		assert.NoError(t, err)

		got, err := asUser.ConfigSection(ctx, created.ID)
		assert.NoError(t, err)

		assert.Equal(t, types.ConfigSection{
			ID:         created.ID,
			ProjectID:  proj.ID,
			Kind:       in.Kind,
			Properties: in.Properties,
			CreatedAt:  created.CreatedAt,
			UpdatedAt:  created.CreatedAt,
		}, got)
	})
}

func TestClient_UpdateConfigSection(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)
	proj := defaultProject(t, asUser)

	t.Run("invalid_config_section_id", func(t *testing.T) {
		_, err := asUser.UpdateConfigSection(ctx, "@nope@", types.UpdateConfigSection{})
		assert.EqualError(t, err, "invalid config section ID")
	})

	t.Run("not_found", func(t *testing.T) {
		_, err := asUser.UpdateConfigSection(ctx, randUUID(t), types.UpdateConfigSection{})
		assert.EqualError(t, err, "config section not found")
	})

	input := types.CreateConfigSection{
		Kind: types.SectionKindOutput,
		Properties: types.Pairs{
			{Key: "name", Value: "http"},
			{Key: "port", Value: 80},
		},
	}
	configSection, err := asUser.CreateConfigSection(ctx, proj.ID, input)
	assert.NoError(t, err)

	t.Run("empty", func(t *testing.T) {
		got, err := asUser.UpdateConfigSection(ctx, configSection.ID, types.UpdateConfigSection{})
		assert.NoError(t, err)
		assert.NotZero(t, got.UpdatedAt)

		found, err := asUser.ConfigSection(ctx, configSection.ID)
		assert.NoError(t, err)
		assert.Equal(t, input.Properties, found.Properties)
	})

	t.Run("ok", func(t *testing.T) {
		newProperties := types.Pairs{
			{Key: "name", Value: "http"},
			{Key: "port", Value: 443},
		}
		got, err := asUser.UpdateConfigSection(ctx, configSection.ID, types.UpdateConfigSection{
			Properties: &newProperties,
		})
		assert.NoError(t, err)
		assert.NotZero(t, got.UpdatedAt)

		found, err := asUser.ConfigSection(ctx, configSection.ID)
		assert.NoError(t, err)
		assert.Equal(t, newProperties, found.Properties)
	})
}

func TestClient_DeleteConfigSection(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)
	proj := defaultProject(t, asUser)

	t.Run("invalid_config_section_id", func(t *testing.T) {
		err := asUser.DeleteConfigSection(ctx, "@nope@")
		assert.EqualError(t, err, "invalid config section ID")
	})

	t.Run("not_found", func(t *testing.T) {
		err := asUser.DeleteConfigSection(ctx, randUUID(t))
		assert.NoError(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		cs, err := asUser.CreateConfigSection(ctx, proj.ID, types.CreateConfigSection{
			Kind: types.SectionKindOutput,
			Properties: types.Pairs{
				{Key: "name", Value: "http"},
				{Key: "port", Value: 80},
			},
		})
		assert.NoError(t, err)

		err = asUser.DeleteConfigSection(ctx, cs.ID)
		assert.NoError(t, err)
	})
}

func TestClient_UpdateConfigSectionSet(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)
	project := defaultProject(t, asUser)

	t.Run("invalid_pipeline_id", func(t *testing.T) {
		err := asUser.UpdateConfigSectionSet(ctx, "@nope@")
		assert.EqualError(t, err, "invalid pipeline ID")
	})

	t.Run("pipeline_not_found", func(t *testing.T) {
		err := asUser.UpdateConfigSectionSet(ctx, randUUID(t))
		assert.EqualError(t, err, "pipeline not found")
	})

	withToken := withToken(t, asUser)
	aggregator, err := withToken.CreateAggregator(ctx, types.CreateAggregator{})
	assert.NoError(t, err)

	pipeline, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
		RawConfig: testFbitConfigWithAddr,
	})
	assert.NoError(t, err)

	t.Run("zero_config_section_ids", func(t *testing.T) {
		err := asUser.UpdateConfigSectionSet(ctx, pipeline.ID)
		assert.NoError(t, err)
	})

	t.Run("invalid_config_section_id", func(t *testing.T) {
		err := asUser.UpdateConfigSectionSet(ctx, pipeline.ID, "@nope@")
		assert.EqualError(t, err, "invalid config section ID")
	})

	t.Run("config_section_set_size_overflow", func(t *testing.T) {
		ids := make([]string, 51)
		for i := 0; i < len(ids); i++ {
			ids[i] = randUUID(t)
		}
		err := asUser.UpdateConfigSectionSet(ctx, pipeline.ID, ids...)
		assert.EqualError(t, err, "config section set size overflow")
	})

	t.Run("not_found", func(t *testing.T) {
		err := asUser.UpdateConfigSectionSet(ctx, pipeline.ID, randUUID(t))
		assert.EqualError(t, err, "config section not found")
	})

	configSection, err := asUser.CreateConfigSection(ctx, project.ID, types.CreateConfigSection{
		Kind: types.SectionKindOutput,
		Properties: types.Pairs{
			{Key: "name", Value: "http"},
			{Key: "port", Value: 80},
		},
	})
	assert.NoError(t, err)

	t.Run("ok", func(t *testing.T) {
		err := asUser.UpdateConfigSectionSet(ctx, pipeline.ID, configSection.ID)
		assert.NoError(t, err)

		foundPipeline, err := asUser.Pipeline(ctx, pipeline.ID, types.PipelineParams{})
		assert.NoError(t, err)
		assert.NotEqual(t, pipeline.CreatedAt, foundPipeline.UpdatedAt)
	})
}
