//go:build integration
// +build integration

package client_test

import (
	"context"
	"testing"

	"github.com/calyptia/api/types"
)

func TestClient_ValidateConfig(t *testing.T) {
	if testFluentbitConfigValidatorAPIKey == "" || testFluentdConfigValidatorAPIKey == "" {
		t.Skip("TODO: setup fluentbit/d config validator API keys")
	}

	ctx := context.Background()
	asUser := userClient(t)
	withToken := withToken(t, asUser)

	// Yes, config validation requires project token auth.
	got, err := withToken.ValidateConfig(ctx, types.AgentTypeFluentBit, types.ValidatingConfig{
		Configs: []types.ValidatingConfigEntry{
			{
				Command: "INPUT",
				Name:    "tail",
			},
			{
				Command: "FILTER",
				Name:    "record_modifier",
			},
			{
				Command: "OUTPUT",
				Name:    "stdout",
			},
		},
	})
	wantEqual(t, err, nil)
	wantEqual(t, len(got.Errors.Runtime), 0)
	wantEqual(t, len(got.Errors.Input), 0)
	wantEqual(t, len(got.Errors.Filter), 0)
	wantEqual(t, len(got.Errors.Output), 0)
}
