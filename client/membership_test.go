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
	wantEqual(t, len(got.Items), 1)
	wantNoEqual(t, got.Items[0].ID, "")
	wantEqual(t, got.Items[0].Roles, []types.MembershipRole{types.MembershipRoleCreator})
	wantNoTimeZero(t, got.Items[0].CreatedAt)
	wantNoEqual(t, got.Items[0].User, nil)
}
