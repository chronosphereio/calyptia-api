//go:build integration
// +build integration

package client_test

import (
	"context"
	"testing"

	"github.com/calyptia/api/types"
)

func TestClient_CreateInvitation(t *testing.T) {
	if testSMTPHost == "" || testSMTPPort == "" {
		t.Skip("TODO: setup SMTP config")
	}

	ctx := context.Background()

	asUser := userClient(t)
	project := defaultProject(t, asUser)

	err := asUser.CreateInvitation(ctx, project.ID, types.CreateInvitation{
		Email:       randStr(t) + "@example.com",
		RedirectURI: "http://cloud-api-testing.localhost/callback",
	})
	wantEqual(t, err, nil)
}
