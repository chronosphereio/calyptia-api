package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/calyptia/api/types"
)

func TestClient_RegisterAgent(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	withToken := withToken(t, asUser)

	metadata, err := json.Marshal(map[string]interface{}{
		"test-key": "test-value",
	})
	if err != nil {
		t.Error(err)
		return
	}

	got, err := withToken.RegisterAgent(ctx, types.RegisterAgent{
		Name:      "test-agent",
		MachineID: "test-machine-id",
		Type:      types.AgentTypeFluentBit,
		Version:   "1.8.6",
		Edition:   types.AgentEditionCommunity,
		Flags:     []string{"test-flag"},
		RawConfig: "test-raw-config",
		Metadata:  (*json.RawMessage)(&metadata),
	})
	wantEqual(t, err, nil)
	wantNoEqual(t, got.ID, "")
	wantNoEqual(t, got.Token, "")
	wantEqual(t, got.Name, "test-agent")
	wantNoTimeZero(t, got.CreatedAt)
}

func TestClient_Agents(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	withToken := withToken(t, asUser)

	jsonMetadata, err := json.Marshal(map[string]interface{}{
		"test-key": "test-value",
	})
	if err != nil {
		t.Error(err)
		return
	}

	rawMetadata := json.RawMessage(jsonMetadata)

	registered, err := withToken.RegisterAgent(ctx, types.RegisterAgent{
		Name:      "test-agent",
		MachineID: "test-machine-id",
		Type:      types.AgentTypeFluentBit,
		Version:   "1.8.6",
		Edition:   types.AgentEditionCommunity,
		Flags:     []string{"test-flag"},
		RawConfig: "test-raw-config",
		Metadata:  &rawMetadata,
	})
	wantEqual(t, err, nil)

	project := defaultProject(t, asUser)
	got, err := asUser.Agents(ctx, project.ID, types.AgentsParams{})
	wantEqual(t, err, nil)
	wantEqual(t, len(got.Items), 1)
	wantEqual(t, got.Items[0].ID, registered.ID)
	wantEqual(t, got.Items[0].Name, registered.Name)
	wantEqual(t, got.Items[0].Token, registered.Token)
	wantEqual(t, got.Items[0].MachineID, "test-machine-id")
	wantEqual(t, got.Items[0].Type, types.AgentTypeFluentBit)
	wantEqual(t, got.Items[0].Version, "1.8.6")
	wantEqual(t, got.Items[0].Edition, types.AgentEditionCommunity)
	wantEqual(t, got.Items[0].Flags, []string{"test-flag"})
	wantEqual(t, got.Items[0].RawConfig, "test-raw-config")
	wantEqual(t, got.Items[0].Metadata, &rawMetadata)
	wantNoTimeZero(t, got.Items[0].CreatedAt)

	// skipping metrics tests here.

	t.Run("pagination", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			_, err := withToken.RegisterAgent(ctx, types.RegisterAgent{
				Name:      fmt.Sprintf("test-agent-%d", i),
				MachineID: fmt.Sprintf("test-machine-id-%d", i), // unique machine id otherwise it gets upserted.
				Type:      types.AgentTypeFluentBit,
				Version:   "v1.8.6",
				Edition:   types.AgentEditionCommunity,
			})
			wantEqual(t, err, nil)
		}

		page1, err := asUser.Agents(ctx, project.ID, types.AgentsParams{
			Last: ptrUint64(3),
		})
		wantEqual(t, err, nil)
		wantEqual(t, len(page1.Items), 3)
		wantNoEqual(t, page1.EndCursor, (*string)(nil))

		page2, err := asUser.Agents(ctx, project.ID, types.AgentsParams{
			Last:   ptrUint64(3),
			Before: page1.EndCursor,
		})
		wantEqual(t, err, nil)
		wantEqual(t, len(page2.Items), 3)
		wantNoEqual(t, page2.EndCursor, (*string)(nil))

		wantNoEqual(t, page1.Items, page2.Items)
		wantNoEqual(t, *page1.EndCursor, *page2.EndCursor)
		wantEqual(t, page1.Items[2].CreatedAt.After(page2.Items[0].CreatedAt), true)
	})
}

func TestClient_Agent(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	withToken := withToken(t, asUser)

	jsonMetadata, err := json.Marshal(map[string]interface{}{
		"test-key": "test-value",
	})
	if err != nil {
		t.Error(err)
		return
	}

	rawMetadata := json.RawMessage(jsonMetadata)

	registered, err := withToken.RegisterAgent(ctx, types.RegisterAgent{
		Name:      "test-agent",
		MachineID: "test-machine-id",
		Type:      types.AgentTypeFluentBit,
		Version:   "1.8.6",
		Edition:   types.AgentEditionCommunity,
		Flags:     []string{"test-flag"},
		RawConfig: "test-raw-config",
		Metadata:  &rawMetadata,
	})
	wantEqual(t, err, nil)

	got, err := asUser.Agent(ctx, registered.ID)
	wantEqual(t, err, nil)
	wantEqual(t, got.ID, registered.ID)
	wantEqual(t, got.Name, registered.Name)
	wantEqual(t, got.Token, registered.Token)
	wantEqual(t, got.MachineID, "test-machine-id")
	wantEqual(t, got.Type, types.AgentTypeFluentBit)
	wantEqual(t, got.Version, "1.8.6")
	wantEqual(t, got.Edition, types.AgentEditionCommunity)
	wantEqual(t, got.Flags, []string{"test-flag"})
	wantEqual(t, got.RawConfig, "test-raw-config")
	wantEqual(t, got.Metadata, &rawMetadata)
	wantNoTimeZero(t, got.CreatedAt)

	// skipping metrics tests here.
}

func TestClient_UpdateAgent(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	withToken := withToken(t, asUser)

	registered, err := withToken.RegisterAgent(ctx, types.RegisterAgent{
		Name:      "test-agent",
		MachineID: "test-machine-id",
		Type:      types.AgentTypeFluentBit,
		Version:   "1.8.6",
		Edition:   types.AgentEditionCommunity,
	})
	wantEqual(t, err, nil)

	err = asUser.UpdateAgent(ctx, registered.ID, types.UpdateAgent{
		// Only name is allowed to be updated by users.
		Name: ptrStr("test-agent-updated"),
	})
	wantEqual(t, err, nil)
}

func TestClient_DeleteAgent(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	withToken := withToken(t, asUser)

	registered, err := withToken.RegisterAgent(ctx, types.RegisterAgent{
		Name:      "test-agent",
		MachineID: "test-machine-id",
		Type:      types.AgentTypeFluentBit,
		Version:   "1.8.6",
		Edition:   types.AgentEditionCommunity,
	})
	wantEqual(t, err, nil)

	err = asUser.DeleteAgent(ctx, registered.ID)
	wantEqual(t, err, nil)
}
