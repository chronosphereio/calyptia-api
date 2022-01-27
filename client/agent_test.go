//go:build integration
// +build integration

package client_test

import (
	"context"
	"encoding/json"
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
	wantEqual(t, len(got), 1)
	wantEqual(t, got[0].ID, registered.ID)
	wantEqual(t, got[0].Name, registered.Name)
	wantEqual(t, got[0].Token, registered.Token)
	wantEqual(t, got[0].MachineID, "test-machine-id")
	wantEqual(t, got[0].Type, types.AgentTypeFluentBit)
	wantEqual(t, got[0].Version, "1.8.6")
	wantEqual(t, got[0].Edition, types.AgentEditionCommunity)
	wantEqual(t, got[0].Flags, []string{"test-flag"})
	wantEqual(t, got[0].RawConfig, "test-raw-config")
	wantEqual(t, got[0].Metadata, &rawMetadata)
	wantNoTimeZero(t, got[0].CreatedAt)

	// skipping metrics tests here.
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
