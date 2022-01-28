package client_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/calyptia/api/types"
)

func TestClient_AgentConfigHistory(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	withToken := withToken(t, asUser)

	registered, err := withToken.RegisterAgent(ctx, types.RegisterAgent{
		Name:      "test-agent",
		MachineID: "test-machine-id",
		Type:      types.AgentTypeFluentBit,
		Version:   "1.8.6",
		Edition:   types.AgentEditionCommunity,
		RawConfig: "test-raw-config",
	})
	wantEqual(t, err, nil)

	withToken.Client = http.DefaultClient

	// update of raw config requires agent token authorization.
	withToken.SetAgentToken(registered.Token)

	err = withToken.UpdateAgent(ctx, registered.ID, types.UpdateAgent{
		RawConfig: ptrStr("test-raw-config-updated"),
	})
	wantEqual(t, err, nil)

	got, err := asUser.AgentConfigHistory(ctx, registered.ID, types.AgentConfigHistoryParams{})
	wantEqual(t, err, nil)
	wantEqual(t, len(got), 2)

	wantNoEqual(t, got[0].ID, "")
	wantEqual(t, got[0].RawConfig, "test-raw-config-updated")
	wantNoTimeZero(t, got[0].CreatedAt)

	wantNoEqual(t, got[1].ID, "")
	wantEqual(t, got[1].RawConfig, "test-raw-config")
	wantNoTimeZero(t, got[1].CreatedAt)
}
