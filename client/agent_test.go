package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/calyptia/api/types"
)

const (
	TagOne   = "tagone"
	TagTwo   = "tagtwo"
	TagThree = "tagthree"
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
		Flags:     []string{"FLB_HAVE_IN_STORAGE_BACKLOG"},
		RawConfig: "test-raw-config",
		Metadata:  (*json.RawMessage)(&metadata),
	})

	wantEqual(t, err, nil)
	wantNoEqual(t, got.ID, "")
	wantNoEqual(t, got.Token, "")
	wantEqual(t, got.Name, "test-agent")
	wantNoTimeZero(t, got.CreatedAt)

	t.Run("invalid agent flags", func(t *testing.T) {
		_, err := withToken.RegisterAgent(ctx, types.RegisterAgent{
			Name:      "test-agent",
			MachineID: "test-machine-id",
			Type:      types.AgentTypeFluentBit,
			Version:   "1.8.6",
			Edition:   types.AgentEditionCommunity,
			Flags:     []string{"\u0000"},
			RawConfig: "test-raw-config",
			Metadata:  (*json.RawMessage)(&metadata),
		})
		wantNoEqual(t, err, nil)
	})

	t.Run("multiple valid agent flags", func(t *testing.T) {
		_, err := withToken.RegisterAgent(ctx, types.RegisterAgent{
			Name:      "test-agent",
			MachineID: "test-machine-id",
			Type:      types.AgentTypeFluentBit,
			Version:   "1.8.6",
			Edition:   types.AgentEditionCommunity,
			Flags: []string{
				"FLB_HAVE_IN_STORAGE_BACKLOG",
				"FLB_HAVE_PARSER",
				"FLB_HAVE_RECORD_ACCESSOR",
				"FLB_HAVE_STREAM_PROCESSOR",
				"JSMN_PARENT_LINKS",
			},
			RawConfig: "test-raw-config",
			Metadata:  (*json.RawMessage)(&metadata),
		})

		wantEqual(t, err, nil)
		wantNoEqual(t, got.ID, "")
		wantNoEqual(t, got.Token, "")
		wantEqual(t, got.Name, "test-agent")
		wantNoTimeZero(t, got.CreatedAt)
	})
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
		Flags:     []string{"FLB_HAVE_IN_STORAGE_BACKLOG"},
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
	wantEqual(t, got.Items[0].Flags, []string{"FLB_HAVE_IN_STORAGE_BACKLOG"})
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
			Last: ptrUint(3),
		})
		wantEqual(t, err, nil)
		wantEqual(t, len(page1.Items), 3)
		wantNoEqual(t, page1.EndCursor, (*string)(nil))

		page2, err := asUser.Agents(ctx, project.ID, types.AgentsParams{
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
		for i := 10; i < 20; i++ {
			agent := types.RegisterAgent{
				Name:      fmt.Sprintf("test-agent-%d", i),
				MachineID: fmt.Sprintf("test-machine-id-%d", i), // unique machine id otherwise it gets upserted.
				Type:      types.AgentTypeFluentBit,
				Version:   "v1.8.6",
				Edition:   types.AgentEditionCommunity,
			}
			if i >= 15 {
				agent.Tags = append(agent.Tags, TagOne, TagThree)
			} else {
				agent.Tags = append(agent.Tags, TagTwo, TagThree)
			}
			_, err := withToken.RegisterAgent(ctx, agent)
			wantEqual(t, err, nil)
		}

		opts := types.AgentsParams{}
		s := TagOne
		opts.Tags = &s
		tag1, err := asUser.Agents(ctx, project.ID, opts)
		wantEqual(t, err, nil)
		wantEqual(t, len(tag1.Items), 5)
		wantEqual(t, tag1.Items[0].Tags, []string{TagOne, TagThree})

		s2 := TagTwo
		opts.Tags = &s2
		tag2, err := asUser.Agents(ctx, project.ID, opts)
		wantEqual(t, err, nil)
		wantEqual(t, len(tag2.Items), 5)
		wantEqual(t, tag2.Items[0].Tags, []string{TagTwo, TagThree})

		s3 := TagThree
		opts.Tags = &s3
		tag3, err := asUser.Agents(ctx, project.ID, opts)
		wantEqual(t, err, nil)
		wantEqual(t, len(tag3.Items), 10)
		wantEqual(t, tag3.Items[0].Tags, []string{TagOne, TagThree})

		s4 := fmt.Sprintf("%s AND %s", TagOne, TagTwo)
		opts.Tags = &s4
		tag4, err := asUser.Agents(ctx, project.ID, opts)
		wantEqual(t, err, nil)
		wantEqual(t, len(tag4.Items), 0)
	})

	t.Run("multiple env", func(t *testing.T) {
		envOne, err := withToken.CreateEnvironment(
			ctx, project.ID, types.CreateEnvironment{Name: "one"})
		wantEqual(t, err, nil)
		wantNoEqual(t, envOne, nil)

		envTwo, err := withToken.CreateEnvironment(
			ctx, project.ID, types.CreateEnvironment{Name: "two"})
		wantEqual(t, err, nil)
		wantNoEqual(t, envTwo, nil)

		aggregatorOne, err := withToken.RegisterAgent(ctx, types.RegisterAgent{
			Name:          "core-instance",
			Version:       "v1.9.5",
			Edition:       types.AgentEditionCommunity,
			EnvironmentID: envOne.ID,
			MachineID:     "core-instance-one",
			Type:          types.AgentTypeFluentBit,
		})
		wantEqual(t, err, nil)
		wantNoEqual(t, aggregatorOne, nil)

		aggregatorTwo, err := withToken.RegisterAgent(ctx, types.RegisterAgent{
			Name:          "core-instance",
			Version:       "v1.9.5",
			Edition:       types.AgentEditionCommunity,
			EnvironmentID: envTwo.ID,
			MachineID:     "core-instance-two",
			Type:          types.AgentTypeFluentBit,
		})
		wantEqual(t, err, nil)
		wantNoEqual(t, aggregatorTwo, nil)

		envOneAggregators, err := asUser.Agents(ctx, project.ID, types.AgentsParams{
			EnvironmentID: ptrStr(envOne.ID),
			Last:          ptrUint(0),
		})
		wantEqual(t, err, nil)
		wantEqual(t, len(envOneAggregators.Items), 1)

		envTwoAggregators, err := asUser.Agents(ctx, project.ID, types.AgentsParams{
			EnvironmentID: ptrStr(envTwo.ID),
			Last:          ptrUint(0),
		})
		wantEqual(t, err, nil)
		wantEqual(t, len(envTwoAggregators.Items), 1)

		bothEnvsAggregators, err := asUser.Agents(ctx, project.ID, types.AgentsParams{
			Last: ptrUint(0),
		})
		wantEqual(t, err, nil)
		wantNoEqual(t, len(bothEnvsAggregators.Items), 1)
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
		Flags:     []string{"FLB_HAVE_IN_STORAGE_BACKLOG"},
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
	wantEqual(t, got.Flags, []string{"FLB_HAVE_IN_STORAGE_BACKLOG"})
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

func TestClient_DeleteAgents(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	project := defaultProject(t, asUser)
	withToken := withToken(t, asUser)

	agent1, err := withToken.RegisterAgent(ctx, types.RegisterAgent{
		Name:      "test-agent-1",
		MachineID: "test-machine-id-1",
		Type:      types.AgentTypeFluentBit,
		Version:   "1.8.6",
		Edition:   types.AgentEditionCommunity,
	})
	wantEqual(t, err, nil)

	agent2, err := withToken.RegisterAgent(ctx, types.RegisterAgent{
		Name:      "test-agent-2",
		MachineID: "test-machine-id-2",
		Type:      types.AgentTypeFluentBit,
		Version:   "1.8.6",
		Edition:   types.AgentEditionCommunity,
	})
	wantEqual(t, err, nil)

	defer func() {
		err := asUser.DeleteAgent(ctx, agent2.ID)
		wantEqual(t, err, nil)
	}()

	err = asUser.DeleteAgents(ctx, project.ID, agent1.ID)
	wantEqual(t, err, nil)

	_, err = asUser.Agent(ctx, agent1.ID)
	wantErrMsg(t, err, "agent not found")

	_, err = asUser.Agent(ctx, agent2.ID)
	wantEqual(t, err, nil)
}
