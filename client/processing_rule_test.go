package client_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"

	"github.com/calyptia/api/client"
	"github.com/calyptia/api/types"
)

var sampleProcessingRuleActions = []types.RuleAction{{
	Kind:        types.RuleActionKindAdd,
	Enabled:     true,
	Description: "Add foobar key-value pair",
	Add:         &types.LogKeyVal{Key: "foo", Value: "bar"},
}, {
	Kind:        types.RuleActionKindRename,
	Enabled:     true,
	Description: `Rename key "rename_me" into "renamed"`,
	Selectors: []types.LogSelector{{
		Kind: types.LogSelectorKindKey,
		Op:   types.LogSelectorOpKindEqual,
		Expr: "rename_me",
	}},
	RenameTo: ptr("renamed"),
}, {
	Kind:        types.RuleActionKindCopy,
	Enabled:     true,
	Description: `Copy key "copy_me" as "copied"`,
	Selectors: []types.LogSelector{{
		Kind: types.LogSelectorKindKey,
		Op:   types.LogSelectorOpKindEqual,
		Expr: "copy_me",
	}},
	CopyAs: ptr("copied"),
}, {
	Kind:        types.RuleActionKindMask,
	Enabled:     true,
	Description: `Mask key "mask_me" with "masked"`,
	Selectors: []types.LogSelector{{
		Kind: types.LogSelectorKindKey,
		Op:   types.LogSelectorOpKindEqual,
		Expr: "mask_me",
	}},
	MaskWith: ptr("masked"),
}, {
	Kind:        types.RuleActionKindRemove,
	Enabled:     true,
	Description: `Remove key "remove_me"`,
	Selectors: []types.LogSelector{{
		Kind: types.LogSelectorKindKey,
		Op:   types.LogSelectorOpKindEqual,
		Expr: "remove_me",
	}},
}, {
	Kind:        types.RuleActionKindSkip,
	Enabled:     false,
	Description: `Skip log record if key "skip_me" is present`,
	Selectors: []types.LogSelector{{
		Kind: types.LogSelectorKindKey,
		Op:   types.LogSelectorOpKindEqual,
		Expr: "skip_me",
	}},
}}

func TestClient_CreateProcessingRule(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)

	t.Run("invalid_pipeline_id", func(t *testing.T) {
		_, err := asUser.CreateProcessingRule(ctx, types.CreateProcessingRule{})
		assert.EqualError(t, err, "invalid pipeline ID")
	})

	t.Run("pipeline_not_found", func(t *testing.T) {
		_, err := asUser.CreateProcessingRule(ctx, types.CreateProcessingRule{
			PipelineID: randUUID(t),
		})
		assert.EqualError(t, err, "pipeline not found")
	})

	agg := setupAggregator(t, withToken(t, asUser))
	pip := setupPipeline(t, asUser, agg.ID)

	t.Run("invalid_lang", func(t *testing.T) {
		_, err := asUser.CreateProcessingRule(ctx, types.CreateProcessingRule{
			PipelineID: pip.ID,
		})
		assert.EqualError(t, err, "invalid processing rule language")
	})

	t.Run("invalid_action_desc", func(t *testing.T) {
		_, err := asUser.CreateProcessingRule(ctx, types.CreateProcessingRule{
			PipelineID: pip.ID,
			Language:   types.ProcessingRuleLanguageLua,
			Actions: []types.RuleAction{{
				Description: strings.Repeat("x", 101),
			}},
		})
		assert.EqualError(t, err, "invalid processing rule action desc")
	})

	t.Run("invalid_action_kind", func(t *testing.T) {
		_, err := asUser.CreateProcessingRule(ctx, types.CreateProcessingRule{
			PipelineID: pip.ID,
			Language:   types.ProcessingRuleLanguageLua,
			Actions:    []types.RuleAction{{}},
		})
		assert.EqualError(t, err, "invalid processing rule action kind")
	})

	t.Run("invalid_action_selector_kind", func(t *testing.T) {
		_, err := asUser.CreateProcessingRule(ctx, types.CreateProcessingRule{
			PipelineID: pip.ID,
			Language:   types.ProcessingRuleLanguageLua,
			Actions: []types.RuleAction{{
				Kind:      types.RuleActionKindAdd,
				Selectors: []types.LogSelector{{}},
			}},
		})
		assert.EqualError(t, err, "invalid processing rule action selector kind")
	})

	t.Run("invalid_action_selector_op", func(t *testing.T) {
		_, err := asUser.CreateProcessingRule(ctx, types.CreateProcessingRule{
			PipelineID: pip.ID,
			Language:   types.ProcessingRuleLanguageLua,
			Actions: []types.RuleAction{{
				Kind: types.RuleActionKindAdd,
				Selectors: []types.LogSelector{{
					Kind: types.LogSelectorKindKey,
				}},
			}},
		})
		assert.EqualError(t, err, "invalid processing rule action selector op")
	})

	t.Run("invalid_action_args_add", func(t *testing.T) {
		_, err := asUser.CreateProcessingRule(ctx, types.CreateProcessingRule{
			PipelineID: pip.ID,
			Language:   types.ProcessingRuleLanguageLua,
			Actions: []types.RuleAction{{
				Kind: types.RuleActionKindAdd,
			}},
		})
		assert.EqualError(t, err, "invalid processing rule action args")
	})

	t.Run("invalid_action_args_rename", func(t *testing.T) {
		_, err := asUser.CreateProcessingRule(ctx, types.CreateProcessingRule{
			PipelineID: pip.ID,
			Language:   types.ProcessingRuleLanguageLua,
			Actions: []types.RuleAction{{
				Kind: types.RuleActionKindRename,
			}},
		})
		assert.EqualError(t, err, "invalid processing rule action args")
	})

	t.Run("invalid_action_args_copy", func(t *testing.T) {
		_, err := asUser.CreateProcessingRule(ctx, types.CreateProcessingRule{
			PipelineID: pip.ID,
			Language:   types.ProcessingRuleLanguageLua,
			Actions: []types.RuleAction{{
				Kind: types.RuleActionKindCopy,
			}},
		})
		assert.EqualError(t, err, "invalid processing rule action args")
	})

	t.Run("invalid_action_args_mask", func(t *testing.T) {
		_, err := asUser.CreateProcessingRule(ctx, types.CreateProcessingRule{
			PipelineID: pip.ID,
			Language:   types.ProcessingRuleLanguageLua,
			Actions: []types.RuleAction{{
				Kind: types.RuleActionKindMask,
			}},
		})
		assert.EqualError(t, err, "invalid processing rule action args")
	})

	t.Run("ok", func(t *testing.T) {
		in := types.CreateProcessingRule{
			PipelineID: pip.ID,
			Language:   types.ProcessingRuleLanguageLua,
			Actions:    sampleProcessingRuleActions,
		}
		got, err := asUser.CreateProcessingRule(ctx, in)
		assert.NoError(t, err)
		assert.NotZero(t, got)

		pr, err := asUser.ProcessingRuleFromPipeline(ctx, pip.ID)
		assert.NoError(t, err)

		assert.Equal(t, got.ID, pr.ID)
		assert.Equal(t, in.Language, pr.Language)

		// revert transport layer
		// automatically converting nil slices into empty arrays.
		for i, a := range pr.Actions {
			if len(a.Selectors) == 0 && a.Selectors != nil {
				pr.Actions[i].Selectors = nil
			}
		}

		assert.Equal(t, in.Actions, pr.Actions)
		assert.Equal(t, got.CreatedAt, pr.CreatedAt)

		testProcessingRuleConfigSection(t, asUser, pr, pr.CreatedAt)
	})

	t.Run("exists", func(t *testing.T) {
		_, err := asUser.CreateProcessingRule(ctx, types.CreateProcessingRule{
			PipelineID: pip.ID,
			Language:   types.ProcessingRuleLanguageLua,
		})
		assert.EqualError(t, err, "processing rule exists")
	})
}

func TestClient_ProcessingRuleFromPipeline(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)

	t.Run("invalid_pipeline_id", func(t *testing.T) {
		_, err := asUser.ProcessingRuleFromPipeline(ctx, "")
		assert.EqualError(t, err, "invalid pipeline ID")
	})

	t.Run("pipeline_not_found", func(t *testing.T) {
		_, err := asUser.ProcessingRuleFromPipeline(ctx, randUUID(t))
		assert.EqualError(t, err, "pipeline not found")
	})

	agg := setupAggregator(t, withToken(t, asUser))
	pip := setupPipeline(t, asUser, agg.ID)

	t.Run("not_found", func(t *testing.T) {
		_, err := asUser.ProcessingRuleFromPipeline(ctx, pip.ID)
		assert.EqualError(t, err, "processing rule not found")
	})

	// ok case already tested on create/update
}

func TestClient_UpdateProcessingRule(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)

	t.Run("invalid_id", func(t *testing.T) {
		_, err := asUser.UpdateProcessingRule(ctx, types.UpdateProcessingRule{})
		assert.EqualError(t, err, "invalid processing rule ID")
	})

	t.Run("not_found", func(t *testing.T) {
		_, err := asUser.UpdateProcessingRule(ctx, types.UpdateProcessingRule{
			ProcessingRuleID: randUUID(t),
		})
		assert.EqualError(t, err, "processing rule not found")
	})

	pr := setupProcessingRule(t, asUser)

	t.Run("invalid_lang", func(t *testing.T) {
		_, err := asUser.UpdateProcessingRule(ctx, types.UpdateProcessingRule{
			ProcessingRuleID: pr.ID,
			Language:         (*types.ProcessingRuleLanguage)(ptr("nope")),
		})
		assert.EqualError(t, err, "invalid processing rule language")
	})

	t.Run("invalid_action_kind", func(t *testing.T) {
		_, err := asUser.UpdateProcessingRule(ctx, types.UpdateProcessingRule{
			ProcessingRuleID: pr.ID,
			Language:         ptr(types.ProcessingRuleLanguageLua),
			Actions: &[]types.RuleAction{
				{Kind: types.RuleActionKind("nope")},
			},
		})
		assert.EqualError(t, err, "invalid processing rule action kind")
	})

	t.Run("ok", func(t *testing.T) {
		in := types.UpdateProcessingRule{
			ProcessingRuleID: pr.ID,
			Language:         ptr(types.ProcessingRuleLanguageLua),
			Actions:          &sampleProcessingRuleActions,
		}
		got, err := asUser.UpdateProcessingRule(ctx, in)
		assert.NoError(t, err)
		assert.NotZero(t, got)

		found, err := asUser.ProcessingRuleFromPipeline(ctx, pr.PipelineID)
		assert.NoError(t, err)

		// revert transport layer
		// automatically converting nil slices into empty arrays.
		for i, a := range found.Actions {
			if len(a.Selectors) == 0 && a.Selectors != nil {
				found.Actions[i].Selectors = nil
			}
		}

		assert.Equal(t, *in.Actions, found.Actions)
		assert.Equal(t, got.UpdatedAt, found.UpdatedAt)

		testProcessingRuleConfigSection(t, asUser, pr, got.UpdatedAt)
	})
}

func testProcessingRuleConfigSection(t *testing.T, client *client.Client, pr types.ProcessingRule, ts time.Time) {
	ctx := context.Background()

	t.Run("pip_config_sections", func(t *testing.T) {
		pip, err := client.Pipeline(ctx, pr.PipelineID, types.PipelineParams{})
		assert.NoError(t, err)
		assert.Equal(t, ts, pip.UpdatedAt)
		assert.NotZero(t, pip.ConfigSections)

		cs := pip.ConfigSections[len(pip.ConfigSections)-1]
		assert.Equal(t, types.SectionKindFilter, cs.Kind)
		assert.Equal(t, pr.ID, *cs.ProcessingRuleID)
		assert.Equal(t, types.Pair{Key: "name", Value: "lua"}, cs.Properties[0])
		assert.Equal(t, types.Pair{Key: "match", Value: "*"}, cs.Properties[1])
		assert.Equal(t, types.Pair{Key: "call", Value: "processing_rule"}, cs.Properties[2])
		// not asserting lua code since it is too much to test,
		// and generated lua code is already tested
		// on the processingrule package.
		// assert.Equal(t, types.Pair{Key: "code", Value: "TODO"}, cs.Properties[3])
	})

	t.Run("proj_config_sections", func(t *testing.T) {
		proj := defaultProject(t, client)
		cc, err := client.ConfigSections(ctx, proj.ID, types.ConfigSectionsParams{
			IncludeProcessingRules: false,
		})
		assert.NoError(t, err)
		assert.Zero(t, cc.Items)

		cc, err = client.ConfigSections(ctx, proj.ID, types.ConfigSectionsParams{
			IncludeProcessingRules: true,
		})
		assert.NoError(t, err)
		assert.NotZero(t, cc.Items)

		cs := cc.Items[len(cc.Items)-1]
		assert.Equal(t, types.SectionKindFilter, cs.Kind)
		assert.Equal(t, pr.ID, *cs.ProcessingRuleID)
		assert.Equal(t, types.Pair{Key: "name", Value: "lua"}, cs.Properties[0])
		assert.Equal(t, types.Pair{Key: "match", Value: "*"}, cs.Properties[1])
		assert.Equal(t, types.Pair{Key: "call", Value: "processing_rule"}, cs.Properties[2])
		// not asserting lua code since it is too much to test,
		// and generated lua code is already tested
		// on the processingrule package.
		// assert.Equal(t, types.Pair{Key: "code", Value: "TODO"}, cs.Properties[3])
	})
}

func setupProcessingRule(t *testing.T, client *client.Client) types.ProcessingRule {
	t.Helper()

	ctx := context.Background()
	agg := setupAggregator(t, withToken(t, client))
	pip := setupPipeline(t, client, agg.ID)
	in := types.CreateProcessingRule{
		PipelineID: pip.ID,
		Language:   types.ProcessingRuleLanguageLua,
		Actions:    sampleProcessingRuleActions,
	}
	created, err := client.CreateProcessingRule(ctx, in)
	assert.NoError(t, err)
	return types.ProcessingRule{
		ID:         created.ID,
		PipelineID: pip.ID,
		Language:   in.Language,
		Actions:    in.Actions,
		CreatedAt:  created.CreatedAt,
		UpdatedAt:  created.CreatedAt,
	}
}
