package client_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"

	"github.com/calyptia/api/client"
	"github.com/calyptia/api/processingrule"
	"github.com/calyptia/api/types"
)

var sampleProcessingRuleActions = []types.RuleAction{{
	Kind:        types.RuleActionKindAdd,
	Enabled:     true,
	Description: "Add foobar key-value pair",
	Add:         &types.LogAttr{Key: "foo", Value: "bar"},
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

	t.Run("invalid_match", func(t *testing.T) {
		_, err := asUser.CreateProcessingRule(ctx, types.CreateProcessingRule{
			PipelineID: pip.ID,
		})
		assert.EqualError(t, err, "invalid processing rule match")
	})

	t.Run("invalid_lang", func(t *testing.T) {
		_, err := asUser.CreateProcessingRule(ctx, types.CreateProcessingRule{
			PipelineID: pip.ID,
			Match:      "*",
		})
		assert.EqualError(t, err, "invalid processing rule language")
	})

	t.Run("invalid_action_desc", func(t *testing.T) {
		_, err := asUser.CreateProcessingRule(ctx, types.CreateProcessingRule{
			PipelineID: pip.ID,
			Match:      "*",
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
			Match:      "*",
			Language:   types.ProcessingRuleLanguageLua,
			Actions:    []types.RuleAction{{}},
		})
		assert.EqualError(t, err, "invalid processing rule action kind")
	})

	t.Run("invalid_action_selector_kind", func(t *testing.T) {
		_, err := asUser.CreateProcessingRule(ctx, types.CreateProcessingRule{
			PipelineID: pip.ID,
			Match:      "*",
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
			Match:      "*",
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
			Match:      "*",
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
			Match:      "*",
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
			Match:      "*",
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
			Match:      "*",
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
			Match:      "*",
			Language:   types.ProcessingRuleLanguageLua,
			Actions:    sampleProcessingRuleActions,
		}
		got, err := asUser.CreateProcessingRule(ctx, in)
		assert.NoError(t, err)
		assert.NotZero(t, got)

		pr, err := asUser.ProcessingRule(ctx, got.ID)
		assert.NoError(t, err)

		assert.Equal(t, got.ID, pr.ID)
		assert.Equal(t, in.Match, pr.Match)
		assert.Equal(t, in.IsMatchRegexp, pr.IsMatchRegexp)
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

		t.Run("file", func(t *testing.T) {
			file, err := asUser.PipelineFile(ctx, got.FileID)
			assert.NoError(t, err)

			code, err := processingrule.ToLua(in.Actions)
			assert.NoError(t, err)
			assert.Equal(t, string(file.Contents), code)

			cs, err := asUser.ConfigSection(ctx, got.ConfigSectionID)
			assert.NoError(t, err)
			scriptPropVal, ok := cs.Properties.Get("script")
			assert.True(t, ok)
			scriptProp, ok := scriptPropVal.(string)
			assert.True(t, ok)
			assert.Equal(t, scriptProp, fmt.Sprintf("{{files.%s}}", file.Name))
		})
	})
}

func TestClient_ProcessingRules(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)

	t.Run("invalid_pipeline_id", func(t *testing.T) {
		_, err := asUser.ProcessingRules(ctx, types.ProcessingRulesParams{})
		assert.EqualError(t, err, "invalid pipeline ID")
	})

	t.Run("pipeline_not_found", func(t *testing.T) {
		_, err := asUser.ProcessingRules(ctx, types.ProcessingRulesParams{
			PipelineID: randUUID(t),
		})
		assert.EqualError(t, err, "pipeline not found")
	})

	pr := setupProcessingRule(t, asUser)

	t.Run("ok", func(t *testing.T) {
		got, err := asUser.ProcessingRules(ctx, types.ProcessingRulesParams{
			PipelineID: pr.PipelineID,
		})
		assert.NoError(t, err)
		assert.NotZero(t, got)

		found := got.Items[0]

		// revert transport layer
		// automatically converting nil slices into empty arrays.
		for i, a := range found.Actions {
			if len(a.Selectors) == 0 && a.Selectors != nil {
				found.Actions[i].Selectors = nil
			}
		}

		assert.Equal(t, pr, found)
	})
}

func TestClient_ProcessingRule(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)

	t.Run("invalid_id", func(t *testing.T) {
		_, err := asUser.ProcessingRule(ctx, "")
		assert.EqualError(t, err, "invalid processing rule ID")
	})

	t.Run("not_found", func(t *testing.T) {
		_, err := asUser.ProcessingRule(ctx, randUUID(t))
		assert.EqualError(t, err, "processing rule not found")
	})

	pr := setupProcessingRule(t, asUser)

	t.Run("ok", func(t *testing.T) {
		_, err := asUser.ProcessingRule(ctx, pr.ID)
		assert.NoError(t, err)
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
			Match:            ptr("foo.*"),
			IsMatchRegexp:    ptr(true),
			Language:         ptr(types.ProcessingRuleLanguageLua),
			Actions:          &sampleProcessingRuleActions,
		}
		got, err := asUser.UpdateProcessingRule(ctx, in)
		assert.NoError(t, err)
		assert.NotZero(t, got)

		found, err := asUser.ProcessingRule(ctx, pr.ID)
		assert.NoError(t, err)

		assert.Equal(t, *in.Match, found.Match)
		assert.Equal(t, *in.IsMatchRegexp, found.IsMatchRegexp)

		// revert transport layer
		// automatically converting nil slices into empty arrays.
		for i, a := range found.Actions {
			if len(a.Selectors) == 0 && a.Selectors != nil {
				found.Actions[i].Selectors = nil
			}
		}

		assert.Equal(t, *in.Actions, found.Actions)
		assert.Equal(t, got.UpdatedAt, found.UpdatedAt)

		t.Run("file", func(t *testing.T) {
			file, err := asUser.PipelineFile(ctx, found.FileID)
			assert.NoError(t, err)

			code, err := processingrule.ToLua(*in.Actions)
			assert.NoError(t, err)
			assert.Equal(t, string(file.Contents), code)
		})
	})
}

func TestClient_DeleteProcessingRule(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)

	t.Run("invalid_id", func(t *testing.T) {
		err := asUser.DeleteProcessingRule(ctx, "@nope@")
		assert.EqualError(t, err, "invalid processing rule ID")
	})

	t.Run("not_found_ok", func(t *testing.T) {
		err := asUser.DeleteProcessingRule(ctx, randUUID(t))
		assert.NoError(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		pr := setupProcessingRule(t, asUser)

		err := asUser.DeleteProcessingRule(ctx, pr.ID)
		assert.NoError(t, err)

		t.Run("idempotent", func(t *testing.T) {
			err := asUser.DeleteProcessingRule(ctx, pr.ID)
			assert.NoError(t, err)
		})
	})
}

func TestClient_PreviewProcessingRule(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)

	now := time.Now().UTC().Truncate(time.Second)

	logs := []types.FluentBitLog{
		{
			Timestamp: types.FluentBitTime(now.Unix()),
			Attrs: types.FluentBitLogAttrs{
				"msg":         "a log message",
				"some_number": float64(33), // JSON will convert it to float64 anyways.
				"some_bool":   true,
				"some_array":  []any{"test"},
				"some_object": map[string]any{
					"test_key": "test_value",
				},
				"copy_me":   "copied_value",
				"rename_me": "renamed_value",
				"mask_me":   "gone",
				"remove_me": "gone",
			},
		},
	}

	actions := []types.RuleAction{
		{
			Kind:    types.RuleActionKindAdd,
			Enabled: true,
			Add:     &types.LogAttr{Key: "added_key", Value: "added_value"},
		},
		{
			Kind:    types.RuleActionKindCopy,
			Enabled: true,
			Selectors: []types.LogSelector{
				{
					Kind: types.LogSelectorKindKey,
					Op:   types.LogSelectorOpKindEqual,
					Expr: "copy_me",
				},
			},
			CopyAs: ptr("copied"),
		},
		{
			Kind:    types.RuleActionKindRename,
			Enabled: true,
			Selectors: []types.LogSelector{
				{
					Kind: types.LogSelectorKindKey,
					Op:   types.LogSelectorOpKindEqual,
					Expr: "rename_me",
				},
			},
			RenameTo: ptr("renamed"),
		},
		{
			Kind:    types.RuleActionKindMask,
			Enabled: true,
			Selectors: []types.LogSelector{
				{
					Kind: types.LogSelectorKindKey,
					Op:   types.LogSelectorOpKindEqual,
					Expr: "mask_me",
				},
			},
			MaskWith: ptr("masked"),
		},
		{
			Kind:    types.RuleActionKindRemove,
			Enabled: true,
			Selectors: []types.LogSelector{
				{
					Kind: types.LogSelectorKindKey,
					Op:   types.LogSelectorOpKindEqual,
					Expr: "remove_me",
				},
			},
		},
	}

	got, err := asUser.PreviewProcessingRule(ctx, types.PreviewProcessingRule{
		Language: types.ProcessingRuleLanguageLua,
		Actions:  actions,
		Logs:     logs,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(got))

	assert.Equal(t, types.FluentBitLogAttrs{
		"msg":         "a log message",
		"some_number": float64(33),
		"some_bool":   true,
		"some_array":  []any{"test"},
		"some_object": map[string]any{
			"test_key": "test_value",
		},
		"added_key": "added_value",
		"copy_me":   "copied_value",
		"copied":    "copied_value",
		"renamed":   "renamed_value",
		"mask_me":   "masked",
	}, got[0].Attrs)
}

func setupProcessingRule(t *testing.T, client *client.Client) types.ProcessingRule {
	t.Helper()

	ctx := context.Background()
	agg := setupAggregator(t, withToken(t, client))
	pip := setupPipeline(t, client, agg.ID)
	in := types.CreateProcessingRule{
		PipelineID:    pip.ID,
		Match:         "*",
		IsMatchRegexp: false,
		Language:      types.ProcessingRuleLanguageLua,
		Actions:       sampleProcessingRuleActions,
	}
	created, err := client.CreateProcessingRule(ctx, in)
	assert.NoError(t, err)
	return types.ProcessingRule{
		ID:              created.ID,
		PipelineID:      pip.ID,
		ConfigSectionID: created.ConfigSectionID,
		FileID:          created.FileID,
		Match:           in.Match,
		IsMatchRegexp:   in.IsMatchRegexp,
		Language:        in.Language,
		Actions:         in.Actions,
		CreatedAt:       created.CreatedAt,
		UpdatedAt:       created.CreatedAt,
	}
}
