package client_test

import (
	"context"
	"testing"

	"github.com/calyptia/api/client"
	"github.com/calyptia/api/types"
)

const zeroUUID = "00000000-0000-0000-0000-000000000000"

func TestClient_CreateIngestCheck(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	project := defaultProject(t, asUser)
	coreInstance := setupAggregator(t, withToken(t, asUser))
	configSection := setupConfigSection(t, project.ID, withToken(t, asUser))

	got, err := asUser.CreateIngestCheck(ctx, coreInstance.ID,
		types.CreateIngestCheck{
			ConfigSectionID: configSection.ID,
			Status:          types.CheckStatusOK,
			Retries:         3,
		},
	)
	wantEqual(t, err, nil)
	wantNoEqual(t, got.ID, "")
	wantNoTimeZero(t, got.CreatedAt)
}

func TestClient_IngestChecks(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	project := defaultProject(t, asUser)
	coreInstance := setupAggregator(t, withToken(t, asUser))
	configSection := setupConfigSection(t, project.ID, withToken(t, asUser))

	t.Run("ok", func(t *testing.T) {
		check, err := asUser.CreateIngestCheck(ctx, coreInstance.ID,
			types.CreateIngestCheck{
				ConfigSectionID: configSection.ID,
				Status:          types.CheckStatusOK,
				Retries:         3,
			},
		)
		wantEqual(t, err, nil)

		got, err := asUser.IngestChecks(ctx, coreInstance.ID, types.IngestChecksParams{})
		wantEqual(t, err, nil)
		wantEqual(t, len(got.Items), 1)

		wantEqual(t, got.Items[0].ID, check.ID)
		wantEqual(t, got.Items[0].CreatedAt, check.CreatedAt)
		wantEqual(t, got.Items[0].UpdatedAt, check.CreatedAt)
	})

	t.Run("pagination", func(t *testing.T) {
		for i := 0; i < 9; i++ {
			configSection := setupConfigSection(t, project.ID, withToken(t, asUser))

			_, err := asUser.CreateIngestCheck(ctx, coreInstance.ID,
				types.CreateIngestCheck{
					ConfigSectionID: configSection.ID,
					Status:          types.CheckStatusOK,
					Retries:         3,
				},
			)
			wantEqual(t, err, nil)
		}

		allIngestChecks, err := asUser.IngestChecks(ctx, coreInstance.ID, types.IngestChecksParams{})
		wantEqual(t, err, nil)
		page1, err := asUser.IngestChecks(ctx, coreInstance.ID, types.IngestChecksParams{Last: ptrUint(3)})
		wantEqual(t, err, nil)
		page2, err := asUser.IngestChecks(ctx, coreInstance.ID, types.IngestChecksParams{Last: ptrUint(3), Before: page1.EndCursor})
		wantEqual(t, err, nil)

		want := allIngestChecks.Items[3:6]
		wantEqual(t, page2.Items, want)
	})
}

func TestClient_IngestCheck(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	project := defaultProject(t, asUser)
	coreInstance := setupAggregator(t, withToken(t, asUser))
	configSection := setupConfigSection(t, project.ID, withToken(t, asUser))

	check, err := asUser.CreateIngestCheck(ctx, coreInstance.ID,
		types.CreateIngestCheck{
			ConfigSectionID: configSection.ID,
			Status:          types.CheckStatusOK,
			Retries:         3,
		},
	)
	wantEqual(t, err, nil)

	got, err := asUser.IngestCheck(ctx, check.ID)
	wantEqual(t, err, nil)
	wantEqual(t, got.ID, check.ID)
	wantEqual(t, got.CreatedAt, check.CreatedAt)
	wantEqual(t, got.UpdatedAt, check.CreatedAt)
}

func TestClient_UpdateIngestCheck(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	project := defaultProject(t, asUser)
	coreInstance := setupAggregator(t, withToken(t, asUser))
	configSection := setupConfigSection(t, project.ID, withToken(t, asUser))

	t.Run("ok", func(t *testing.T) {
		check, err := asUser.CreateIngestCheck(ctx, coreInstance.ID,
			types.CreateIngestCheck{
				ConfigSectionID: configSection.ID,
				Status:          types.CheckStatusOK,
				Retries:         3,
			},
		)
		wantEqual(t, err, nil)

		err = asUser.UpdateIngestCheck(ctx, check.ID,
			types.UpdateIngestCheck{
				Status: ptr(types.CheckStatusFailed),
			},
		)
		wantEqual(t, err, nil)
	})

	t.Run("invalid config section", func(t *testing.T) {
		configSection := setupConfigSection(t, project.ID, withToken(t, asUser))

		check, err := asUser.CreateIngestCheck(ctx, coreInstance.ID,
			types.CreateIngestCheck{
				ConfigSectionID: configSection.ID,
				Status:          types.CheckStatusOK,
				Retries:         3,
			},
		)
		wantEqual(t, err, nil)

		err = asUser.UpdateIngestCheck(ctx, check.ID,
			types.UpdateIngestCheck{
				ConfigSectionID: ptr(zeroUUID),
			},
		)
		wantNoEqual(t, err, nil)
	})
}

func TestClient_DeleteIngestCheck(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	project := defaultProject(t, asUser)
	coreInstance := setupAggregator(t, withToken(t, asUser))
	configSection := setupConfigSection(t, project.ID, withToken(t, asUser))

	check, err := asUser.CreateIngestCheck(ctx, coreInstance.ID,
		types.CreateIngestCheck{
			ConfigSectionID: configSection.ID,
			Status:          types.CheckStatusOK,
			Retries:         3,
		},
	)
	wantEqual(t, err, nil)
	err = asUser.DeleteIngestCheck(ctx, check.ID)
	wantEqual(t, err, nil)
}

func setupConfigSection(t *testing.T, projectID string, withToken *client.Client) types.CreatedConfigSection {
	t.Helper()
	ctx := context.Background()

	configSection, err := withToken.CreateConfigSection(ctx, projectID, types.CreateConfigSection{
		Kind: types.SectionKindOutput,
		Properties: types.Pairs{
			{
				Key:   "name",
				Value: ptr("es"),
			},
		},
	})
	wantEqual(t, err, nil)
	return configSection
}
