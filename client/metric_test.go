package client_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/calyptia/api/types"
)

func TestClient_AggregatorPipelinesMetrics(t *testing.T) {
	ctx := context.Background()
	asUser := userClient(t)

	project := defaultProject(t, asUser)
	token := defaultToken(t, asUser)

	coreInstance, err := setupCoreInstance(dockerPool, asUser.BaseURL, token)
	wantEqual(t, err, nil)
	t.Cleanup(func() {
		_ = coreInstance.Close()
	})

	withToken := withToken(t, asUser)

	var aggregators types.Aggregators

	err = retry(func() error {
		aggregators, err = withToken.Aggregators(ctx, project.ID, types.AggregatorsParams{})
		if err != nil || len(aggregators.Items) == 0 {
			return fmt.Errorf("empty aggregators")
		}
		return nil
	})

	wantEqual(t, err, nil)
	wantNoEqual(t, len(aggregators.Items), 0)

	agg := aggregators.Items[0]

	var pipelines types.Pipelines

	err = retry(func() error {
		pipelines, err = withToken.Pipelines(ctx, agg.ID, types.PipelinesParams{})
		if err != nil || len(pipelines.Items) == 0 {
			return fmt.Errorf("empty pipelines")
		}
		return nil
	})

	wantEqual(t, err, nil)
	wantNoEqual(t, len(pipelines.Items), 0)

	metricsParams := types.MetricsParams{
		Start:    -24 * time.Hour,
		Interval: 1 * time.Second,
	}

	t.Run("ok", func(t *testing.T) {
		var metrics types.PipelinesMetrics

		err = retry(func() error {
			var pipelineIDs []string

			for _, pipeline := range pipelines.Items {
				pipelineIDs = append(pipelineIDs, pipeline.ID)
			}

			metrics, err = withToken.AggregatorPipelinesMetrics(ctx, agg.ID, types.PipelinesMetricsParams{
				MetricsParams: metricsParams,
				PipelineIDs:   pipelineIDs,
			})

			if err != nil || len(metrics) == 0 {
				return fmt.Errorf("no metrics found")
			}

			return nil
		})

		wantEqual(t, err, nil)
		wantNoEqual(t, metrics, types.PipelinesMetrics{})
	})

	t.Run("zero_pipeline_ids", func(t *testing.T) {
		err = retry(func() error {
			_, err = withToken.AggregatorPipelinesMetrics(ctx, agg.ID, types.PipelinesMetricsParams{
				MetricsParams: metricsParams,
				PipelineIDs:   []string{},
			})
			return err
		})
		wantNoEqual(t, err, nil)
	})

	t.Run("pipeline_not_found", func(t *testing.T) {
		notFoundPipelineID := "7b45e316-8fa8-4a5f-8736-d0dee22c52ad"
		var metrics types.PipelinesMetrics
		err = retry(func() error {
			var pipelineIDs []string

			for _, pipeline := range pipelines.Items {
				pipelineIDs = append(pipelineIDs, pipeline.ID)
			}

			pipelineIDs = append(pipelineIDs, notFoundPipelineID)
			metrics, err = withToken.AggregatorPipelinesMetrics(ctx, agg.ID, types.PipelinesMetricsParams{
				MetricsParams: metricsParams,
				PipelineIDs:   pipelineIDs,
			})

			if err != nil && len(metrics) == 0 {
				return fmt.Errorf("no metrics found")
			}

			return err
		})

		wantEqual(t, err, nil)
		wantEqual(t, metrics[notFoundPipelineID].Error, "pipeline not found")
	})
}
