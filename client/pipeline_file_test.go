package client_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/calyptia/api/types"
)

func TestClient_CreatePipelineFile(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	got, err := asUser.CreatePipelineFile(ctx, pipeline.ID, types.CreatePipelineFile{
		Name:      "testfile",
		Contents:  []byte("test-contents"),
		Encrypted: true,
	})
	wantEqual(t, err, nil)
	wantNoEqual(t, got.ID, "")
	wantNoTimeZero(t, got.CreatedAt)
}

func TestClient_PipelineFiles(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	t.Run("ok", func(t *testing.T) {
		file, err := asUser.CreatePipelineFile(ctx, pipeline.ID, types.CreatePipelineFile{
			Name:      "testfile",
			Contents:  []byte("test-contents"),
			Encrypted: true,
		})
		wantEqual(t, err, nil)

		got, err := asUser.PipelineFiles(ctx, pipeline.ID, types.PipelineFilesParams{})
		wantEqual(t, err, nil)
		wantEqual(t, len(got.Items), 2) // Aditional "parsers" file should be created by default with each pipeline.

		wantEqual(t, got.Items[0].ID, file.ID)
		wantEqual(t, got.Items[0].Name, "testfile")
		wantNoEqual(t, got.Items[0].Contents, []byte("test-contents"))
		wantEqual(t, got.Items[0].Encrypted, true)
		wantEqual(t, got.Items[0].CreatedAt, file.CreatedAt)
		wantEqual(t, got.Items[0].UpdatedAt, file.CreatedAt)

		wantEqual(t, got.Items[1].Name, "parsers")
	})
	t.Run("pagination", func(t *testing.T) {
		for i := 0; i < 9; i++ {
			_, err := asUser.CreatePipelineFile(ctx, pipeline.ID, types.CreatePipelineFile{
				Name:      "testfile" + strconv.Itoa(i),
				Contents:  []byte("test-contents"),
				Encrypted: true,
			})
			wantEqual(t, err, nil)
		}

		allPipelineFiles, err := asUser.PipelineFiles(ctx, pipeline.ID, types.PipelineFilesParams{})
		wantEqual(t, err, nil)
		page1, err := asUser.PipelineFiles(ctx, pipeline.ID, types.PipelineFilesParams{Last: ptrUint64(3)})
		wantEqual(t, err, nil)
		page2, err := asUser.PipelineFiles(ctx, pipeline.ID, types.PipelineFilesParams{Last: ptrUint64(3), Before: page1.EndCursor})
		wantEqual(t, err, nil)

		want := allPipelineFiles.Items[3:6]
		wantEqual(t, page2.Items, want)
	})
}

func TestClient_PipelineFile(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	file, err := asUser.CreatePipelineFile(ctx, pipeline.ID, types.CreatePipelineFile{
		Name:      "testfile",
		Contents:  []byte("test-contents"),
		Encrypted: true,
	})
	wantEqual(t, err, nil)

	got, err := asUser.PipelineFile(ctx, file.ID)
	wantEqual(t, err, nil)
	wantEqual(t, got.ID, file.ID)
	wantEqual(t, got.Name, "testfile")
	wantNoEqual(t, got.Contents, []byte("test-contents"))
	wantEqual(t, got.Encrypted, true)
	wantEqual(t, got.CreatedAt, file.CreatedAt)
	wantEqual(t, got.UpdatedAt, file.CreatedAt)
}

func TestClient_UpdatePipelineFile(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	file, err := asUser.CreatePipelineFile(ctx, pipeline.ID, types.CreatePipelineFile{
		Name:      "testfile",
		Contents:  []byte("test-contents"),
		Encrypted: true,
	})
	wantEqual(t, err, nil)

	err = asUser.UpdatePipelineFile(ctx, file.ID, types.UpdatePipelineFile{
		Name:      ptrStr("testfileupdated"),
		Contents:  ptrBytes([]byte("test-contents-updated")),
		Encrypted: ptrBool(true),
	})
	wantEqual(t, err, nil)
}

func TestClient_DeletePipelineFile(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	file, err := asUser.CreatePipelineFile(ctx, pipeline.ID, types.CreatePipelineFile{
		Name:      "testfile",
		Contents:  []byte("test-contents"),
		Encrypted: true,
	})
	wantEqual(t, err, nil)

	err = asUser.DeletePipelineFile(ctx, file.ID)
	wantEqual(t, err, nil)
}

func Test_decryptPipelineFile(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	aggregator := setupAggregator(t, withToken(t, asUser))
	pipeline := setupPipeline(t, asUser, aggregator.ID)

	file, err := asUser.CreatePipelineFile(ctx, pipeline.ID, types.CreatePipelineFile{
		Name:      "testfile",
		Contents:  []byte("test-contents"),
		Encrypted: true,
	})
	wantEqual(t, err, nil)

	got, err := asUser.PipelineFile(ctx, file.ID)
	wantEqual(t, err, nil)

	decrypted := decrypt(t, got.Contents, aggregator.PrivateRSAKey)
	wantEqual(t, decrypted, []byte("test-contents"))
}
