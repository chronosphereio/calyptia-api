// Package msgpack05x implements the cmetrics format in cmetrics 0.5.x.
//
// It has the following structure:
// {
//
//	"meta": {
//		"cmetrics": {
//		},
//		"external": {
//		},
//		"processing": {
//		    "static_labels": [
//			[
//			    "pipeline_id",
//			    "84056c1e-c872-47ec-887c-56a32f9d0704"
//			]
//		    ]
//		}
//	    },
//	    "metrics": [
//		{
//		    "meta": {
//			"ver": 2,
//			"type": 0,
//			"opts": {
//			    "ns": "fluentbit",
//			    "ss": "",
//			    "name": "uptime",
//			    "desc": "Number of seconds that Fluent Bit has been running."
//			},
//			"labels": [
//			    "hostname"
//			],
//			"aggregation_type": 2
//		    },
//		    "values": [
//			{
//			    "ts": 1669299892279587459,
//			    "value": 60.0,
//			    "labels": [
//				"calyptia-health-check-71b1-5798c65dcf-h7rfq"
//			    ],
//			    "hash": 11716548876681578473
//			}
//		    ]
//		},
//		{
//		    "meta": {
//			"ver": 2,
//			"type": 0,
//			"opts": {
//			    "ns": "fluentbit",
//			    "ss": "input",
//			    "name": "bytes_total",
//			    "desc": "Number of input bytes."
//			},
//			"labels": [
//			    "name"
//			],
//			"aggregation_type": 2
//		    },
//		    "values": [
//			{
//			    "ts": 1669299832361890876,
//			    "value": 0.0,
//			    "labels": [
//				"fluentbit_metrics.0"
//			    ],
//			    "hash": 4185491748558833568
//			}
//		    ]
//		},
//	}
package msgpack05x

import (
	"fmt"
	"strings"
	"time"

	"github.com/influxdata/influxdb/models"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/calyptia/api/types"
)

type MetaHeader struct {
	CMetrics   struct{} `msgpack:"cmetrics"`
	External   struct{} `msgpack:"external"`
	Processing struct {
		StaticLabels [][]string `msgpack:"static_labels"`
	} `msgpack:"processing"`
}

type Meta struct {
	Ver  int `msgpack:"ver"`
	Type int `msgpack:"type"`
	Opts struct {
		NS   string `msgpack:"ns"`
		SS   string `msgpack:"ss"`
		Name string `msgpack:"name"`
		Desc string `msgpack:"desc"`
	} `msgpack:"opts"`
	Labels          []string `msgpack:"labels"`
	StaticLabels    []string `msgpack:"static_labels"`
	AggregationType int      `msgpack:"aggregation_type"`
}

type Value struct {
	TS     uint64   `msgpack:"ts"`
	Value  float64  `msgpack:"value"`
	Labels []string `msgpack:"labels"`
	Hash   uint64   `msgpack:"hash"`
}

type Values []Value

type Metric struct {
	Meta   Meta   `msgpack:"meta"`
	Values Values `msgpack:"values"`
}

type Metrics []Metric

type cmetrics5 struct {
	Meta    MetaHeader `msgpack:"meta"`
	Metrics Metrics    `msgpack:"metrics"`
}

func makeTags(fields, values []string, static [][]string) models.Tags {
	tags := make(models.Tags, 0)

	for idx, field := range fields {
		if idx < len(values) {
			tags = append(tags, models.Tag{
				Key:   []byte(field),
				Value: []byte(values[idx]),
			})
		}
	}

	for _, set := range static {
		for i := 0; i+1 < len(set); i += 2 {
			tags = append(tags, models.Tag{
				Key:   []byte(set[i]),
				Value: []byte(set[i+1]),
			})
		}
	}

	return tags
}

func makeName(meta Meta) string {
	if len(meta.Opts.NS) > 0 && len(meta.Opts.SS) > 0 {
		return strings.Join([]string{
			meta.Opts.NS,
			meta.Opts.SS,
		}, "_")
	}
	return meta.Opts.NS
}

func Decode(data []byte) ([]types.Metric, error) {
	var metaHeader cmetrics5

	err := msgpack.Unmarshal(data, &metaHeader)
	if err != nil {
		return []types.Metric{}, err
	}

	metrics := make([]types.Metric, 0)
	for _, measurement := range metaHeader.Metrics {
		for _, val := range measurement.Values {
			tags := makeTags(measurement.Meta.Labels, val.Labels, metaHeader.Meta.Processing.StaticLabels)
			point, err := models.NewPoint(
				makeName(measurement.Meta),
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

func Encode(metrics []types.Metric) ([]byte, error) {
	// type cmetrics5 struct {
	//	Meta    MetaHeader `msgpack:"meta"`
	//	Metrics Metrics    `msgpack:"metrics"`
	// }
	// type MetaHeader struct {
	//	CMetrics   struct{} `msgpack:"cmetrics"`
	//	External   struct{} `msgpack:"external"`
	//	Processing struct {
	//		StaticLabels [][]string `msgpack:"static_labels"`
	//	} `msgpack:"processing"`
	// }
	metaheader := cmetrics5{}
	metaheader.Meta = MetaHeader{}
	metaheader.Meta.Processing = struct {
		StaticLabels [][]string "msgpack:\"static_labels\""
	}{
		StaticLabels: make([][]string, 0),
	}

	for _, metric := range metrics {
		tags := metric.Point.Tags()
		for _, tag := range tags {
			metaheader.Meta.Processing.StaticLabels = append(metaheader.Meta.Processing.StaticLabels, []string{
				string(tag.Key),
				string(tag.Value),
			})
		}
	}

	metaheader.Metrics = make(Metrics, 0)
	for _, metric := range metrics {
		values := make(Values, 0)

		fields, _ := metric.Fields()
		for _, v := range fields {
			val, ok := v.(float64)
			if !ok {
				return []byte{}, fmt.Errorf("illegal metrics value: %s", v)
			}
			values = append(values, Value{
				TS:    uint64(time.Now().Unix()),
				Value: val,
			})
		}

		metaheader.Metrics = append(metaheader.Metrics, Metric{
			Meta: Meta{
				Opts: struct {
					NS   string "msgpack:\"ns\""
					SS   string "msgpack:\"ss\""
					Name string "msgpack:\"name\""
					Desc string "msgpack:\"desc\""
				}{
					NS:   "fluentbit_output",
					Name: "fluentbit_output",
				},
			},
			Values: values,
		})
	}

	return msgpack.Marshal(&metaheader)
}
