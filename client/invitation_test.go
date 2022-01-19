package client_test

import (
	"context"
	"testing"

	"github.com/calyptia/api/types"
)

func TestClient_CreateInvitation(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	project := defaultProject(t, asUser)

	err := asUser.CreateInvitation(ctx, project.ID, types.CreateInvitation{
		Email:       randStr(t) + "@example.com",
		RedirectURI: testBaseURL + "/callback",
	})
	wantEqual(t, err, nil)
}
