package client_test

import (
	"context"
	"testing"

	"github.com/calyptia/api/types"
)

func TestClient_CreateProject(t *testing.T) {
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		asUser := userClient(t)
		got, err := asUser.CreateProject(ctx, types.CreateProject{
			Name: "test-project",
		})
		wantEqual(t, err, nil)
		wantNoEqual(t, got.ID, "")
		wantNoTimeZero(t, got.CreatedAt)
		wantNoEqual(t, got.Token, "")
		wantNoEqual(t, got.Membership.ID, "")
		wantNoTimeZero(t, got.Membership.CreatedAt)
		wantEqual(t, got.Membership.Roles, []types.MembershipRole{types.MembershipRoleCreator})
	})
}

func TestClient_Projects(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	got, err := asUser.Projects(ctx, types.ProjectsParams{})
	wantEqual(t, err, nil)
	wantEqual(t, len(got), 1) // a default project must be created for the user as a side-effect.
	wantNoEqual(t, got[0].ID, "")
	wantEqual(t, got[0].Name, "default")
	wantEqual(t, got[0].MembersCount, uint64(1))
	wantNoTimeZero(t, got[0].CreatedAt)
	wantNoEqual(t, got[0].Membership, nil)
	wantNoEqual(t, got[0].Membership.ID, "")
	wantNoTimeZero(t, got[0].Membership.CreatedAt)
	wantEqual(t, got[0].Membership.Roles, []types.MembershipRole{types.MembershipRoleCreator})
}

func TestClient_Project(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	project := defaultProject(t, asUser)

	got, err := asUser.Project(ctx, project.ID)
	wantEqual(t, err, nil)
	wantEqual(t, got.ID, project.ID)
	wantEqual(t, got.Name, "default")
	wantEqual(t, got.MembersCount, uint64(1))
	wantEqual(t, got.CreatedAt, project.CreatedAt)
	wantNoEqual(t, got.Membership, nil)
	wantEqual(t, *got.Membership, *project.Membership)
}

func TestClient_UpdateProject(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	project := defaultProject(t, asUser)

	err := asUser.UpdateProject(ctx, project.ID, types.UpdateProject{Name: ptrStr("test-project-updated")})
	wantEqual(t, err, nil)
}
