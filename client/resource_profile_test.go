//go:build integration
// +build integration

package client_test

import (
	"context"
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

	got, err := asUser.ResourceProfiles(ctx, aggregator.ID, types.ResourceProfilesParams{})
	wantEqual(t, err, nil)
	wantEqual(t, len(got), 3) // By default, there are 3 resource profiles.

	// Sort resource profiles by name
	// since they are created at the same time
	// and the order may switch.
	sort.Slice(got, func(i, j int) bool {
		return got[i].Name < got[j].Name
	})

	// TODO: check on the expected profile spec values.

	wantNoEqual(t, got[0].ID, "")
	wantEqual(t, got[0].Name, types.ResourceProfileBestEffortLowResource)
	wantNoTimeZero(t, got[0].CreatedAt)
	wantNoTimeZero(t, got[0].UpdatedAt)

	wantNoEqual(t, got[1].ID, "")
	wantEqual(t, got[1].Name, types.ResourceProfileHighPerformanceGuaranteedDelivery)
	wantNoTimeZero(t, got[1].CreatedAt)
	wantNoTimeZero(t, got[1].UpdatedAt)

	wantNoEqual(t, got[2].ID, "")
	wantEqual(t, got[2].Name, types.ResourceProfileHighPerformanceOptimalThroughput)
	wantNoTimeZero(t, got[2].CreatedAt)
	wantNoTimeZero(t, got[2].UpdatedAt)
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
