package client_test

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/calyptia/api/types"
)

func TestClient_PipelineConfig(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))

	t.Run("fix for cloud/632", func(t *testing.T) {
		pip, err := asUser.CreatePipeline(ctx, aggregator.ID, types.CreatePipeline{
			Name: "test",
			Kind: types.PipelineKindDeployment,
			RawConfig: `
{
   "pipeline":{
      "inputs":[
         {
            "dummy":{
               "dummy":"{\"message\":\"dummy\"}",
               "rate":"1",
               "samples":"0",
               "start_time_sec":"-1",
               "start_time_nsec":"-1",
               "Name":"dummy",
               "tag":"dummy.2a227ce1-6a50-4ad2-ba37-9e436c454116"
            }
         }
      ],
      "outputs":[
         {
            "http":{
               "host":"calyptia-vivo",
               "port":"5489",
               "uri":"/console",
               "format":"json",
               "Name":"http",
               "Match":"*"
            }
         }
      ]
   }
}
`,
			ConfigFormat: types.ConfigFormatJSON,
		})

		assert.NoError(t, err)
		assert.NotZero(t, pip)

		pr, err := asUser.CreateProcessingRule(ctx, types.CreateProcessingRule{
			PipelineID:    pip.ID,
			Match:         "*",
			IsMatchRegexp: false,
			Language:      types.ProcessingRuleLanguageLua,
			Actions: []types.RuleAction{
				{
					Kind:        types.RuleActionKindAdd,
					Description: "",
					Enabled:     false,
					Selectors:   nil,
					Add: &types.LogAttr{
						Key:   "aa",
						Value: "bb",
					},
				},
			},
		})
		assert.NoError(t, err)
		assert.NotZero(t, pr)

		format := types.ConfigFormatINI
		fetchPip, err := asUser.Pipeline(ctx, pip.ID, types.PipelineParams{
			RenderWithConfigSections: true,
			ConfigFormat:             &format,
		})
		assert.NoError(t, err)
		assert.NotZero(t, fetchPip)
		assert.Contains(t, fetchPip.Config.RawConfig, `[INPUT]
    dummy {"message":"dummy"}
    rate 1
    samples 0
    start_time_sec -1
    start_time_nsec -1
    Name dummy
    tag dummy.2a227ce1-6a50-4ad2-ba37-9e436c454116
[OUTPUT]
    host calyptia-vivo
    port 5489
    uri /console
    format json
    Name http
    Match *
[FILTER]
    name lua
    match *
    call processing_rule`)
	})
}

func TestClient_PipelineConfigHistory(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	t.Run("ok", func(t *testing.T) {
		_, err := asUser.UpdatePipeline(ctx, pipeline.ID, types.UpdatePipeline{
			RawConfig: ptr(testFbitConfigWithAddr3),
		})
		wantEqual(t, err, nil)

		got, err := asUser.PipelineConfigHistory(ctx, pipeline.ID, types.PipelineConfigHistoryParams{})
		wantEqual(t, err, nil)

		wantEqual(t, len(got.Items), 2) // Initial config should be already there by default.

		wantNoEqual(t, got.Items[0].ID, "")
		wantEqual(t, got.Items[0].RawConfig, testFbitConfigWithAddr3)
		wantNoTimeZero(t, got.Items[0].CreatedAt)

		wantNoEqual(t, got.Items[1].ID, "")
		wantEqual(t, got.Items[1].RawConfig, testFbitConfigWithAddr)
		wantNoTimeZero(t, got.Items[1].CreatedAt)
	})

	t.Run("pagination", func(t *testing.T) {
		for i := 0; i < 9; i++ {
			_, err := asUser.UpdatePipeline(ctx, pipeline.ID, types.UpdatePipeline{
				RawConfig: ptr(strings.ReplaceAll(testFbitConfigWithAddr, "24224", strconv.Itoa(24224+i))),
			})
			wantEqual(t, err, nil)
		}
		allConfigs, err := asUser.PipelineConfigHistory(ctx, pipeline.ID, types.PipelineConfigHistoryParams{})
		wantEqual(t, err, nil)
		page1, err := asUser.PipelineConfigHistory(ctx, pipeline.ID, types.PipelineConfigHistoryParams{Last: ptrUint(3)})
		wantEqual(t, err, nil)
		page2, err := asUser.PipelineConfigHistory(ctx, pipeline.ID, types.PipelineConfigHistoryParams{Last: ptrUint(3), Before: page1.EndCursor})
		wantEqual(t, err, nil)

		want := allConfigs.Items[3:6]
		wantEqual(t, page2.Items, want)
	})
}
