package client_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/calyptia/api/types"
)

func TestClient_CreateEnvironment(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctx := context.Background()
		asUser := userClient(t)
		withToken := withToken(t, asUser)
		project := defaultProject(t, asUser)
		got, err := withToken.CreateEnvironment(ctx, project.ID, types.CreateEnvironment{
			Name: "production",
		})
		wantEqual(t, err, nil)
		wantNoTimeZero(t, got.CreatedAt)
	})
	t.Run("duplicated", func(t *testing.T) {
		ctx := context.Background()
		asUser := userClient(t)
		withToken := withToken(t, asUser)
		project := defaultProject(t, asUser)
		_, err := withToken.CreateEnvironment(ctx, project.ID, types.CreateEnvironment{
			Name: "production",
		})
		if err != nil {
			t.Error(err)
		}
		_, err = withToken.CreateEnvironment(ctx, project.ID, types.CreateEnvironment{
			Name: "production",
		})
		wantEqual(t, err.Error(), "environment name taken")
	})
	t.Run("invalid name", func(t *testing.T) {
		ctx := context.Background()
		asUser := userClient(t)
		withToken := withToken(t, asUser)
		project := defaultProject(t, asUser)
		_, err := withToken.CreateEnvironment(ctx, project.ID, types.CreateEnvironment{
			Name: "invalid**",
		})
		wantEqual(t, err.Error(), "invalid environment name")
	})
}

func TestClient_Environments(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	withToken := withToken(t, asUser)
	project := defaultProject(t, asUser)
	created, err := withToken.CreateEnvironment(ctx, project.ID, types.CreateEnvironment{
		Name: "test-environment",
	})
	wantEqual(t, err, nil)

	got, err := asUser.Environments(ctx, project.ID, types.EnvironmentsParams{})
	wantEqual(t, err, nil)
	wantEqual(t, len(got.Items), 2)
	wantEqual(t, got.Items[0].ID, created.ID)
	wantNoTimeZero(t, got.Items[0].CreatedAt)

	t.Run("pagination", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			_, err := withToken.CreateEnvironment(ctx, project.ID, types.CreateEnvironment{
				Name: fmt.Sprintf("test-environment-%d", i),
			})
			wantEqual(t, err, nil)
		}

		page1, err := asUser.Environments(ctx, project.ID, types.EnvironmentsParams{
			Last: ptrUint(3),
		})
		wantEqual(t, err, nil)
		wantEqual(t, len(page1.Items), 3)
		wantNoEqual(t, page1.EndCursor, (*string)(nil))

		page2, err := asUser.Environments(ctx, project.ID, types.EnvironmentsParams{
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
}

func TestClient_UpdateEnvironment(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctx := context.Background()

		asUser := userClient(t)
		withToken := withToken(t, asUser)

		project := defaultProject(t, asUser)
		created, err := withToken.CreateEnvironment(ctx, project.ID, types.CreateEnvironment{
			Name: "test-environment",
		})
		wantEqual(t, err, nil)

		err = asUser.UpdateEnvironment(ctx, created.ID, types.UpdateEnvironment{
			Name: ptr("test-environment-updated"),
		})
		wantEqual(t, err, nil)
	})
	t.Run("already existent name", func(t *testing.T) {
		ctx := context.Background()

		asUser := userClient(t)
		withToken := withToken(t, asUser)

		project := defaultProject(t, asUser)
		_, err := withToken.CreateEnvironment(ctx, project.ID, types.CreateEnvironment{
			Name: "existent-name",
		})
		wantEqual(t, err, nil)
		created, err := withToken.CreateEnvironment(ctx, project.ID, types.CreateEnvironment{
			Name: "test-environment",
		})
		wantEqual(t, err, nil)

		err = asUser.UpdateEnvironment(ctx, created.ID, types.UpdateEnvironment{
			Name: ptr("existent-name"),
		})
		wantEqual(t, err.Error(), "environment name taken")
	})
}

func TestClient_DeleteEnvironment(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctx := context.Background()

		asUser := userClient(t)
		withToken := withToken(t, asUser)
		project := defaultProject(t, asUser)
		created, err := withToken.CreateEnvironment(ctx, project.ID, types.CreateEnvironment{
			Name: "test-environment",
		})
		wantEqual(t, err, nil)

		err = asUser.DeleteEnvironment(ctx, created.ID)
		wantEqual(t, err, nil)
	})
	t.Run("deleting non-existing environment", func(t *testing.T) {
		ctx := context.Background()

		asUser := userClient(t)
		withToken := withToken(t, asUser)
		project := defaultProject(t, asUser)
		_, err := withToken.CreateEnvironment(ctx, project.ID, types.CreateEnvironment{
			Name: "test-environment",
		})
		wantEqual(t, err, nil)

		err = asUser.DeleteEnvironment(ctx, randUUID(t))
		wantEqual(t, err.Error(), "permission denied")
	})
}

func TestClient_EnvironmentAssociation(t *testing.T) {
	t.Run("should create agent with default environment", func(t *testing.T) {
		ctx := context.Background()

		asUser := userClient(t)
		withToken := withToken(t, asUser)

		agent, err := withToken.RegisterAgent(ctx, types.RegisterAgent{
			Name:      "test-agent",
			MachineID: "test-machine-id",
			Type:      types.AgentTypeFluentBit,
			Version:   "1.8.6",
			Edition:   types.AgentEditionCommunity,
			RawConfig: "test-raw-config",
		})
		wantEqual(t, err, nil)
		got, err := withToken.Agent(ctx, agent.ID)
		wantEqual(t, err, nil)
		wantEqual(t, got.EnvironmentName, "default")
	})
	t.Run("should associate agent with given environment", func(t *testing.T) {
		ctx := context.Background()

		asUser := userClient(t)
		withToken := withToken(t, asUser)
		project := defaultProject(t, asUser)
		envName := "test-environment"
		env, err := withToken.CreateEnvironment(ctx, project.ID, types.CreateEnvironment{
			Name: envName,
		})
		wantEqual(t, err, nil)
		agent, err := withToken.RegisterAgent(ctx, types.RegisterAgent{
			Name:          "test-agent",
			MachineID:     "test-machine-id",
			Type:          types.AgentTypeFluentBit,
			Version:       "1.8.6",
			Edition:       types.AgentEditionCommunity,
			RawConfig:     "test-raw-config",
			EnvironmentID: env.ID,
		})
		wantEqual(t, err, nil)
		got, err := withToken.Agent(ctx, agent.ID)
		wantEqual(t, err, nil)
		wantEqual(t, got.EnvironmentName, envName)
	})
}
