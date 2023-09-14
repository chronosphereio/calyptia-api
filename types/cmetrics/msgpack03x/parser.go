// Package msgpack03x is a parser for the cmetrics 0.3.x format. The
// format has the following structure:
// [
//
//	{
//		"meta": {
//			"ver": 2,
//			"type": 0,
//			"opts": {
//				"ns": "fluentbit",
//				"ss": "",
//				"name": "uptime",
//				"desc": "Number of seconds that Fluent Bit has been running."
//			},
//			"label_dictionary": [
//				"hostname",
//				"Phillips-MBP.lan"
//			],
//			"static_labels": [],
//			"labels": [0]
//		},
//		"values": [{
//			"ts": 1669390294568988000,
//			"value": 5,
//			"labels": [1],
//			"hash": 8754447662303377000
//		}]
//	},
//	{
//		"meta": {
//			"ver": 2,
//			"type": 0,
//			"opts": {
//				"ns": "fluentbit",
//				"ss": "input",
//				"name": "bytes_total",
//				"desc": "Number of input bytes."
//			},
//			"label_dictionary": [
//				"name",
//				"dummy.0"
//			],
//			"static_labels": [],
//			"labels": [0]
//		},
//		"values": [{
//			"ts": 1669390294568876000,
//			"value": 130,
//			"labels": [1],
//			"hash": 6783070066932087000
//		}]
//	}
//
// ]
package msgpack03x

import (
	"strings"
	"time"

	"github.com/influxdata/influxdb/models"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/calyptia/api/types"
)

type Metav3 struct {
	Ver  int `msgpack:"ver"`
	Type int `msgpack:"type"`
	Opts struct {
		NS   string `msgpack:"ns"`
		SS   string `msgpack:"ss"`
		Name string `msgpack:"name"`
		Desc string `msgpack:"desc"`
	} `msgpack:"opts"`
	LabelDictionary []string `msgpack:"label_dictionary"`
	StaticLabels    []string `msgpack:"static_label"`
	Labels          []int64  `msgpack:"labels"`
}

type ValueV3 struct {
	TS     uint64  `msgpack:"ts"`
	Value  float64 `msgpack:"value"`
	Labels []int64 `msgpack:"labels"`
	Hash   uint64  `msgpack:"hash"`
}

type ValuesV3 []ValueV3

type MetricV3 struct {
	Meta   Metav3   `msgpack:"meta"`
	Values ValuesV3 `msgpack:"values"`
}

type MetricsV3 []MetricV3

type cmetrics3 MetricsV3

// The current version fo the parser simply assumes all the labels
// from label_dictionary are used and ignores the labels properties.
// From what I can tell these might have been an attempt to use a
// dictionary to compress the cmetrics stream.
func makeTagsV3(tagDictionary []string) models.Tags {
	tags := make(models.Tags, 0)

	for idx := 0; idx < len(tagDictionary)-1; idx += 2 {
		tags = append(tags, models.Tag{
			Key:   []byte(tagDictionary[idx]),
			Value: []byte(tagDictionary[idx+1]),
		})
	}
	return tags
}

func makeNameV3(meta Metav3) string {
	if len(meta.Opts.NS) > 0 && len(meta.Opts.SS) > 0 {
		return strings.Join([]string{
			meta.Opts.NS,
			meta.Opts.SS,
		}, "_")
	}
	return meta.Opts.NS
}

func Decode(data []byte) ([]types.Metric, error) {
	var metaHeader cmetrics3

	err := msgpack.Unmarshal(data, &metaHeader)
	if err != nil {
		return []types.Metric{}, err
	}

	metrics := make([]types.Metric, 0)
	for _, measurement := range metaHeader {
		for _, val := range measurement.Values {
			tags := makeTagsV3(measurement.Meta.LabelDictionary)
			point, err := models.NewPoint(
				makeNameV3(measurement.Meta),
				tags,
				models.Fields{
					measurement.Meta.Opts.Name: val.Value,
				},
				time.Unix(0, int64(val.TS)))
			if err != nil {
				return []types.Metric{}, err
			}
			metrics = append(metrics, types.Metric{
				Measurement: measurement.Meta.Opts.Name,
				Point:       point,
			})
		}
	}

	return metrics, nil
}
