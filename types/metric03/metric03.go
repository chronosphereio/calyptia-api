// Package metric03 holds the structure for the old mskpack cmetrics version 0.3.x.
package metric03

import (
	"fmt"

	"github.com/calyptia/api/types"
)

type CreateMetrics []Metric

type Metric struct {
	Meta   Meta    `json:"meta" msgpack:"meta"`
	Values []Value `json:"values" msgpack:"values"`
}

type Meta struct {
	LabelDictionary []string   `json:"label_dictionary" msgpack:"label_dictionary"` // ex: ["foo", "bar", "my-label", "my-value", "name", "dummy.0"]
	Labels          []int      `json:"labels" msgpack:"labels"`                     // Tells which labels are dynamic. Ex: [4]
	Opts            Opts       `json:"opts" msgpack:"opts"`
	StaticLabels    []int      `json:"static_labels" msgpack:"static_labels"` // Tells which labels are static. Ex: [0, 1, 2, 3]
	Type            MetricType `json:"type" msgpack:"type"`
	// Ver             int        `json:"ver" msgpack:"ver"`
}

type Opts struct {
	Desc      string `json:"desc" msgpack:"desc"`
	Name      string `json:"name" msgpack:"name"` // ex: "uptime", "bytes_total"
	Namespace string `json:"ns" msgpack:"ns"`     // ex: "fluentbit"
	Subsystem string `json:"ss" msgpack:"ss"`     // ex: "input"
}

type MetricType int

type Value struct {
	// Hash   int64   `json:"hash" msgpack:"hash"`
	Labels []int   `json:"labels" msgpack:"labels"` // Tells the dynamic label value. Ex: [5]
	TS     int64   `json:"ts" msgpack:"ts"`         // nanoseconds
	Value  float64 `json:"value" msgpack:"value"`
}

// Upgrade transforms the old mspack cmetrics version 0.3.x to the current version.
func (in CreateMetrics) Upgrade() (types.CreateMetrics, error) {
	var out types.CreateMetrics
	for _, oldmetric := range in {
		// set static labels only once, since in the old metric type the static labels are repeated on each metric.
		if out.Meta.Processing.StaticLabels == nil {
			if len(oldmetric.Meta.StaticLabels)%2 != 0 {
				return out, fmt.Errorf("invalid static labels: not even")
			}

			for i := 0; i < len(oldmetric.Meta.StaticLabels); i += 2 {
				if len(oldmetric.Meta.LabelDictionary) <= i+1 {
					return out, fmt.Errorf("invalid static label index overflow: %d", i+1)
				}

				label := oldmetric.Meta.LabelDictionary[i]
				value := oldmetric.Meta.LabelDictionary[i+1]

				out.Meta.Processing.StaticLabels = append(out.Meta.Processing.StaticLabels, [2]string{label, value})
			}
		}

		meta := types.MetricMeta{
			// AggregationType: types.MetricAggregationTypeUnspecified,
			Opts: types.MetricMetaOpts{
				Desc:      oldmetric.Meta.Opts.Desc,
				Name:      oldmetric.Meta.Opts.Name,
				Namespace: oldmetric.Meta.Opts.Namespace,
				Subsystem: oldmetric.Meta.Opts.Subsystem,
			},
			Type: types.MetricType(oldmetric.Meta.Type),
			// Ver:  oldmetric.Meta.Ver,
		}
		for _, labelIndex := range oldmetric.Meta.Labels {
			if len(oldmetric.Meta.LabelDictionary) <= labelIndex {
				return out, fmt.Errorf("invalid label index overflow: %d", labelIndex)
			}

			meta.Labels = append(meta.Labels, oldmetric.Meta.LabelDictionary[labelIndex])
		}

		var values []types.MetricValue
		for _, oldvalue := range oldmetric.Values {
			if len(meta.Labels) != len(oldvalue.Labels) {
				return out, fmt.Errorf("invalid label length: %d != %d", len(meta.Labels), len(oldvalue.Labels))
			}

			value := types.MetricValue{
				// Hash:  oldvalue.Hash,
				TS:    oldvalue.TS,
				Value: oldvalue.Value,
			}

			for _, labelValueIndex := range oldvalue.Labels {
				if len(oldmetric.Meta.LabelDictionary) <= labelValueIndex {
					return out, fmt.Errorf("invalid label value index overflow: %d", labelValueIndex)
				}

				value.Labels = append(value.Labels, oldmetric.Meta.LabelDictionary[labelValueIndex])
			}

			values = append(values, value)
		}

		metric := types.Metric{
			Meta:   meta,
			Values: values,
		}

		out.Metrics = append(out.Metrics, metric)
	}

	return out, nil
}
