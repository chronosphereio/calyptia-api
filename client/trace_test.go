package client_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/calyptia/api/client"
	"github.com/calyptia/api/types"
)

func TestClient_CreateTraceSession(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	agg := setupAggregator(t, withToken(t, asUser))
	pip := setupPipeline(t, asUser, agg.ID)

	t.Run("plugin_not_found", func(t *testing.T) {
		_, err := asUser.CreateTraceSession(ctx, pip.ID, types.CreateTraceSession{
			Plugins:  []string{"nope"},
			Lifespan: types.Duration(time.Minute * 10),
		})
		wantErrMsg(t, err, "plugin not found")
	})

	t.Run("ok", func(t *testing.T) {
		_ = setupTraceSession(t, setupTraceSessionOpts{
			Client:     asUser,
			PipelineID: &pip.ID,
		})
	})
}

func TestClient_TraceSessions(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	sess := setupTraceSession(t, setupTraceSessionOpts{
		Client: asUser,
	})

	got, err := asUser.TraceSessions(ctx, sess.PipelineID, types.TraceSessionsParams{})
	wantEqual(t, err, nil)
	wantEqual(t, len(got.Items), 1)
	wantNoEqual(t, got.EndCursor, nil)
	wantEqual(t, got.Items[0].ID, sess.ID)
}

func TestClient_TraceSession(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	sess := setupTraceSession(t, setupTraceSessionOpts{
		Client: asUser,
	})

	got, err := asUser.TraceSession(ctx, sess.ID)
	wantEqual(t, err, nil)
	wantEqual(t, got.ID, sess.ID)
}

func TestClient_ActiveTraceSession(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	sess := setupTraceSession(t, setupTraceSessionOpts{
		Client: asUser,
	})

	got, err := asUser.ActiveTraceSession(ctx, sess.PipelineID)
	wantEqual(t, err, nil)
	wantEqual(t, got.ID, sess.ID)
}

func TestClient_UpdateTraceSession(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	sess := setupTraceSession(t, setupTraceSessionOpts{
		Client: asUser,
	})

	newLifespan := time.Minute * 33
	got, err := asUser.UpdateTraceSession(ctx, sess.ID, types.UpdateTraceSession{
		Lifespan: (*types.Duration)(&newLifespan),
	})
	wantEqual(t, err, nil)
	wantNoTimeZero(t, got.UpdatedAt)

	found, err := asUser.TraceSession(ctx, sess.ID)
	wantEqual(t, err, nil)
	wantEqual(t, found.Lifespan.AsDuration(), newLifespan)
}

func TestClient_TerminateActiveTraceSession(t *testing.T) {
	ctx := context.Background()

	asUser := userClient(t)
	sess := setupTraceSession(t, setupTraceSessionOpts{
		Client: asUser,
	})

	got, err := asUser.TerminateActiveTraceSession(ctx, sess.PipelineID)
	wantEqual(t, err, nil)
	wantEqual(t, got.ID, sess.ID)
	wantNoTimeZero(t, got.UpdatedAt)

	found, err := asUser.TraceSession(ctx, sess.ID)
	wantEqual(t, err, nil)
	wantEqual(t, found.Active(), false)
}

func TestClient_CreateTraceRecord(t *testing.T) {
	_ = setupTraceRecord(t, setupTraceRecordOpts{})
}

func TestClient_TraceRecords(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)

	sess := setupTraceSession(t, setupTraceSessionOpts{
		Client: asUser,
	})

	for i := 0; i < 3; i++ {
		_ = setupTraceRecord(t, setupTraceRecordOpts{
			Client:     asUser,
			PipelineID: &sess.PipelineID,
		})
	}

	got, err := asUser.TraceRecords(ctx, sess.ID, types.TraceRecordsParams{})
	wantEqual(t, err, nil)
	wantEqual(t, len(got.Items), 3)
	wantNoEqual(t, got.EndCursor, nil)
}

type setupTraceSessionOpts struct {
	Client     *client.Client
	PipelineID *string
	Plugins    *[]string
	Lifespan   *time.Duration
}

func setupTraceSession(t *testing.T, opts setupTraceSessionOpts) types.TraceSession {
	t.Helper()

	ctx := context.Background()
	var client *client.Client
	if opts.Client != nil {
		client = opts.Client
	} else {
		client = userClient(t)
	}

	var pipelineID string
	if opts.PipelineID != nil {
		pipelineID = *opts.PipelineID
	} else {
		agg := setupAggregator(t, withToken(t, client))
		pip := setupPipeline(t, client, agg.ID)
		pipelineID = pip.ID
	}

	in := types.CreateTraceSession{
		Plugins:  []string{"forward"}, // the test pipeline has a forward plugin.
		Lifespan: types.Duration(time.Minute * 10),
	}
	if opts.Plugins != nil {
		in.Plugins = *opts.Plugins
	}
	if opts.Lifespan != nil {
		in.Lifespan = types.Duration(*opts.Lifespan)
	}

	got, err := client.CreateTraceSession(ctx, pipelineID, in)
	wantEqual(t, err, nil)
	wantNoEqual(t, got.ID, "")
	wantNoTimeZero(t, got.CreatedAt)

	return types.TraceSession{
		ID:         got.ID,
		PipelineID: pipelineID,
		Plugins:    in.Plugins,
		Lifespan:   in.Lifespan,
		CreatedAt:  got.CreatedAt,
		UpdatedAt:  got.CreatedAt,
	}
}

type setupTraceRecordOpts struct {
	Client     *client.Client
	PipelineID *string
	Input      *types.CreateTraceRecord
}

//nolint:unparam // returning trace record might be used in the tests later.
func setupTraceRecord(t *testing.T, opts setupTraceRecordOpts) types.TraceRecord {
	t.Helper()

	ctx := context.Background()

	var client *client.Client
	if opts.Client != nil {
		client = opts.Client
	} else {
		client = userClient(t)
	}

	var pipelineID string
	if opts.PipelineID != nil {
		pipelineID = *opts.PipelineID
	} else {
		sess := setupTraceSession(t, setupTraceSessionOpts{Client: client})
		pipelineID = sess.PipelineID
	}

	now := time.Now().UTC().Truncate(time.Microsecond)
	nowUnixStr := strconv.FormatInt(now.Unix(), 10)

	var in types.CreateTraceRecord
	if opts.Input != nil {
		in = *opts.Input
	} else {
		in = types.CreateTraceRecord{
			Kind:           types.TraceRecordKindInput,
			TraceID:        "dummy.0",
			StartTime:      now.Unix(),
			EndTime:        now.Add(time.Second).Unix(),
			PluginInstance: "dummy",
			PluginAlias:    "mydummy",
			ReturnCode:     0,
			Records: []byte(`[
			{
				"timestamp": ` + nowUnixStr + `,
				"record": {
					"dummy": "dummy_0",
                    "powered_by": "calyptia",
                    "data": {
                      "key_name": "foo",
                      "key_cnt": "1"
                    }
				}
			}
		]`),
		}
	}

	got, err := client.CreateTraceRecord(ctx, pipelineID, in)
	wantEqual(t, err, nil)
	wantNoEqual(t, got.ID, "")
	wantNoTimeZero(t, got.CreatedAt)

	return types.TraceRecord{
		ID:             got.ID,
		SessionID:      got.SessionID, // TODO: this is empty right now.
		Kind:           in.Kind,
		TraceID:        in.TraceID,
		StartTime:      time.Unix(in.StartTime, 0),
		EndTime:        time.Unix(in.EndTime, 0),
		PluginInstance: in.PluginInstance,
		PluginAlias:    in.PluginAlias,
		ReturnCode:     in.ReturnCode,
		Records:        in.Records,
		CreatedAt:      got.CreatedAt,
	}
}
