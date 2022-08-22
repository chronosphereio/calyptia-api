package client_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/calyptia/api/types"
)

func TestClient_CreateToken(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	project := defaultProject(t, asUser)

	got, err := asUser.CreateToken(ctx, project.ID, types.CreateToken{
		Name: "test-token",
	})
	wantEqual(t, err, nil)
	wantNoEqual(t, got.ID, "")
	wantEqual(t, got.Name, "test-token")
	wantNoEqual(t, got.Token, "")
	wantNoTimeZero(t, got.CreatedAt)
}

func TestClient_Tokens(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	project := defaultProject(t, asUser)

	t.Run("ok", func(t *testing.T) {
		token, err := asUser.CreateToken(ctx, project.ID, types.CreateToken{
			Name: "test-token",
		})
		wantEqual(t, err, nil)

		got, err := asUser.Tokens(ctx, project.ID, types.TokensParams{})
		wantEqual(t, err, nil)

		wantEqual(t, len(got.Items), 2) // Aditional "default" token should be created by default with each project.

		wantEqual(t, got.Items[0], token)
		wantEqual(t, got.Items[1].Name, "default")
	})
	t.Run("pagination", func(t *testing.T) {
		for i := 0; i < 9; i++ {
			_, err := asUser.CreateToken(ctx, project.ID, types.CreateToken{
				Name: "test-token" + strconv.Itoa(i),
			})
			wantEqual(t, err, nil)
		}

		allTokens, err := asUser.Tokens(ctx, project.ID, types.TokensParams{})
		wantEqual(t, err, nil)
		page1, err := asUser.Tokens(ctx, project.ID, types.TokensParams{Last: ptrUint(3)})
		wantEqual(t, err, nil)
		page2, err := asUser.Tokens(ctx, project.ID, types.TokensParams{Last: ptrUint(3), Before: page1.EndCursor})
		wantEqual(t, err, nil)

		want := allTokens.Items[3:6]
		wantEqual(t, page2.Items, want)
	})
}

func TestClient_Token(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	project := defaultProject(t, asUser)

	token, err := asUser.CreateToken(ctx, project.ID, types.CreateToken{
		Name: "test-token",
	})
	wantEqual(t, err, nil)

	got, err := asUser.Token(ctx, token.ID)
	wantEqual(t, err, nil)
	wantEqual(t, got, token)
}

func TestClient_UpdateToken(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	project := defaultProject(t, asUser)

	token, err := asUser.CreateToken(ctx, project.ID, types.CreateToken{
		Name: "test-token",
	})
	wantEqual(t, err, nil)

	err = asUser.UpdateToken(ctx, token.ID, types.UpdateToken{
		Name: ptrStr("test-token-updated"),
	})
	wantEqual(t, err, nil)
}

func TestClient_DeleteToken(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	project := defaultProject(t, asUser)

	token, err := asUser.CreateToken(ctx, project.ID, types.CreateToken{
		Name: "test-token",
	})
	wantEqual(t, err, nil)

	err = asUser.DeleteToken(ctx, token.ID)
	wantEqual(t, err, nil)
}
