package client_test

import (
	"context"
	"fmt"
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
	t.Run("ok", func(t *testing.T) {
		err = withToken.UpdateAgent(ctx, registered.ID, types.UpdateAgent{
			RawConfig: ptrStr("test-raw-config-updated"),
		})
		wantEqual(t, err, nil)

		got, err := asUser.AgentConfigHistory(ctx, registered.ID, types.AgentConfigHistoryParams{})
		wantEqual(t, err, nil)
		wantEqual(t, len(got.Items), 2)

		wantNoEqual(t, got.Items[0].ID, "")
		wantEqual(t, got.Items[0].RawConfig, "test-raw-config-updated")
		wantNoTimeZero(t, got.Items[0].CreatedAt)

		wantNoEqual(t, got.Items[1].ID, "")
		wantEqual(t, got.Items[1].RawConfig, "test-raw-config")
		wantNoTimeZero(t, got.Items[1].CreatedAt)
	})
	t.Run("pagination", func(t *testing.T) {
		for i := 0; i < 9; i++ {
			err = withToken.UpdateAgent(ctx, registered.ID, types.UpdateAgent{
				RawConfig: ptrStr(fmt.Sprintf("test-raw-config-updated-%d", i)),
			})
			wantEqual(t, err, nil)
		}
		allConfigs, err := asUser.AgentConfigHistory(ctx, registered.ID, types.AgentConfigHistoryParams{})
		wantEqual(t, err, nil)
		page1, err := asUser.AgentConfigHistory(ctx, registered.ID, types.AgentConfigHistoryParams{Last: ptrUint(3)})
		wantEqual(t, err, nil)
		page2, err := asUser.AgentConfigHistory(ctx, registered.ID, types.AgentConfigHistoryParams{Last: ptrUint(3), Before: page1.EndCursor})
		wantEqual(t, err, nil)

		want := allConfigs.Items[3:6]
		wantEqual(t, page2.Items, want)
	})
}
