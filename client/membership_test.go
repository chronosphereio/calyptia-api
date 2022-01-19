package client_test

import (
	"context"
	"testing"

	"github.com/calyptia/api/types"
)

func TestClient_Members(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	project := defaultProject(t, asUser)

	got, err := asUser.Members(ctx, project.ID, types.MembersParams{})
	wantEqual(t, err, nil)
	wantEqual(t, len(got), 1)
	wantNoEqual(t, got[0].ID, "")
	wantEqual(t, got[0].Roles, []types.MembershipRole{types.MembershipRoleCreator})
	wantNoTimeZero(t, got[0].CreatedAt)
	wantNoEqual(t, got[0].User, nil)
}
