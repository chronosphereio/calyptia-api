package client_test

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/calyptia/api/types"
)

func TestClient_CreateResourceProfile(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))

	got, err := asUser.CreateResourceProfile(ctx, aggregator.ID, types.CreateResourceProfile{
		Name:                   "test-resource-profile",
		StorageMaxChunksUp:     3,
		StorageSyncFull:        true,
		StorageBacklogMemLimit: "40Mi",
		StorageVolumeSize:      "30Mi",
		StorageMaxChunksPause:  true,
		CPUBufferWorkers:       2,
		CPULimit:               "30Mi",
		CPURequest:             "20Mi",
		MemoryLimit:            "35Mi",
		MemoryRequest:          "25Mi",
	})
	wantEqual(t, err, nil)
	wantNoEqual(t, got.ID, "")
	wantNoTimeZero(t, got.CreatedAt)
}

func TestClient_ResourceProfiles(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))

	t.Run("ok", func(t *testing.T) {
		got, err := asUser.ResourceProfiles(ctx, aggregator.ID, types.ResourceProfilesParams{})
		wantEqual(t, err, nil)
		wantEqual(t, len(got.Items), 3) // By default, there are 3 resource profiles.

		// Sort resource profiles by name
		// since they are created at the same time
		// and the order may switch.
		sort.Slice(got.Items, func(i, j int) bool {
			return got.Items[i].Name < got.Items[j].Name
		})

		// TODO: check on the expected profile spec values.

		wantNoEqual(t, got.Items[0].ID, "")
		wantEqual(t, got.Items[0].Name, types.ResourceProfileBestEffortLowResource)
		wantNoTimeZero(t, got.Items[0].CreatedAt)
		wantNoTimeZero(t, got.Items[0].UpdatedAt)

		wantNoEqual(t, got.Items[1].ID, "")
		wantEqual(t, got.Items[1].Name, types.ResourceProfileHighPerformanceGuaranteedDelivery)
		wantNoTimeZero(t, got.Items[1].CreatedAt)
		wantNoTimeZero(t, got.Items[1].UpdatedAt)

		wantNoEqual(t, got.Items[2].ID, "")
		wantEqual(t, got.Items[2].Name, types.ResourceProfileHighPerformanceOptimalThroughput)
		wantNoTimeZero(t, got.Items[2].CreatedAt)
		wantNoTimeZero(t, got.Items[2].UpdatedAt)
	})
	t.Run("pagination", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			_, err := asUser.CreateResourceProfile(ctx, aggregator.ID, types.CreateResourceProfile{
				Name:                   fmt.Sprintf("test-resource-profile-%d", i),
				StorageMaxChunksUp:     3,
				StorageSyncFull:        true,
				StorageBacklogMemLimit: "40Mi",
				StorageVolumeSize:      "30Mi",
				StorageMaxChunksPause:  true,
				CPUBufferWorkers:       2,
				CPULimit:               "30Mi",
				CPURequest:             "20Mi",
				MemoryLimit:            "35Mi",
				MemoryRequest:          "25Mi",
			})
			wantEqual(t, err, nil)
		}
		allResourceProfile, err := asUser.ResourceProfiles(ctx, aggregator.ID, types.ResourceProfilesParams{})
		wantEqual(t, err, nil)
		page1, err := asUser.ResourceProfiles(ctx, aggregator.ID, types.ResourceProfilesParams{Last: ptrUint64(3)})
		wantEqual(t, err, nil)
		page2, err := asUser.ResourceProfiles(ctx, aggregator.ID, types.ResourceProfilesParams{Last: ptrUint64(3), Before: page1.EndCursor})
		wantEqual(t, err, nil)

		want := allResourceProfile.Items[3:6]
		wantEqual(t, page2.Items, want)
	})
}

func TestClient_ResourceProfile(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))

	resourceProfile, err := asUser.CreateResourceProfile(ctx, aggregator.ID, types.CreateResourceProfile{
		Name:                   "test-resource-profile",
		StorageMaxChunksUp:     3,
		StorageSyncFull:        true,
		StorageBacklogMemLimit: "40Mi",
		StorageVolumeSize:      "30Mi",
		StorageMaxChunksPause:  true,
		CPUBufferWorkers:       2,
		CPULimit:               "30Mi",
		CPURequest:             "20Mi",
		MemoryLimit:            "35Mi",
		MemoryRequest:          "25Mi",
	})
	wantEqual(t, err, nil)

	got, err := asUser.ResourceProfile(ctx, resourceProfile.ID)
	wantEqual(t, err, nil)
	wantEqual(t, got.ID, resourceProfile.ID)
	wantEqual(t, got.Name, "test-resource-profile")
	wantEqual(t, got.StorageMaxChunksUp, 3)
	wantEqual(t, got.StorageSyncFull, true)
	wantEqual(t, got.StorageBacklogMemLimit, "40Mi")
	wantEqual(t, got.StorageVolumeSize, "30Mi")
	wantEqual(t, got.StorageMaxChunksPause, true)
	wantEqual(t, got.CPUBufferWorkers, 2)
	wantEqual(t, got.CPULimit, "30Mi")
	wantEqual(t, got.CPURequest, "20Mi")
	wantEqual(t, got.MemoryLimit, "35Mi")
	wantEqual(t, got.MemoryRequest, "25Mi")
	wantEqual(t, got.CreatedAt, resourceProfile.CreatedAt)
	wantEqual(t, got.UpdatedAt, resourceProfile.CreatedAt)
}

func TestClient_UpdateResourceProfile(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))

	resourceProfile, err := asUser.CreateResourceProfile(ctx, aggregator.ID, types.CreateResourceProfile{
		Name:                   "test-resource-profile",
		StorageMaxChunksUp:     3,
		StorageSyncFull:        true,
		StorageBacklogMemLimit: "40Mi",
		StorageVolumeSize:      "30Mi",
		StorageMaxChunksPause:  true,
		CPUBufferWorkers:       2,
		CPULimit:               "30Mi",
		CPURequest:             "20Mi",
		MemoryLimit:            "35Mi",
		MemoryRequest:          "25Mi",
	})
	wantEqual(t, err, nil)

	err = asUser.UpdateResourceProfile(ctx, resourceProfile.ID, types.UpdateResourceProfile{
		Name:                   ptrStr("test-resource-profile-updated"),
		StorageMaxChunksUp:     ptrUint(4),
		StorageSyncFull:        ptrBool(false),
		StorageMaxChunksPause:  ptrBool(false),
		StorageBacklogMemLimit: ptrStr("51Mi"),
		StorageVolumeSize:      ptrStr("41Mi"),
		CPUBufferWorkers:       ptrUint(3),
		CPULimit:               ptrStr("41Mi"),
		CPURequest:             ptrStr("31Mi"),
		MemoryLimit:            ptrStr("46Mi"),
		MemoryRequest:          ptrStr("36Mi"),
	})
	wantEqual(t, err, nil)
}

func TestClient_DeleteResourceProfile(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))

	resourceProfile, err := asUser.CreateResourceProfile(ctx, aggregator.ID, types.CreateResourceProfile{
		Name:                   "test-resource-profile",
		StorageMaxChunksUp:     3,
		StorageSyncFull:        true,
		StorageBacklogMemLimit: "40Mi",
		StorageVolumeSize:      "30Mi",
		StorageMaxChunksPause:  true,
		CPUBufferWorkers:       2,
		CPULimit:               "30Mi",
		CPURequest:             "20Mi",
		MemoryLimit:            "35Mi",
		MemoryRequest:          "25Mi",
	})
	wantEqual(t, err, nil)

	err = asUser.DeleteResourceProfile(ctx, resourceProfile.ID)
	wantEqual(t, err, nil)
}
