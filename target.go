package main

import (
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// Target is the internal state about a target
type Target struct {
	sync.Mutex
	spec         TargetSpec
	registry     *prometheus.Registry
	measurements *Measurements
}

func NewTarget(ts TargetSpec) *Target {
	t := Target{
		spec:     ts,
		registry: prometheus.NewRegistry()}

	log.Println("new target:", ts.host)

	//prometheus.WrapRegistererWith(prometheus.Labels{"zone": "aaa"}, t.registry).MustRegister(t)
	t.registry.MustRegister(&t)

	return &t
}

var (
	fpingSentCountDesc = prometheus.NewDesc(
		"fping_sent_count",
		"Number of sent probes",
		nil, nil,
	)

	fpingLostCountDesc = prometheus.NewDesc(
		"fping_lost_count",
		"Number of lost probes",
		nil, nil,
	)

	fpingRTTSumDesc = prometheus.NewDesc(
		"fping_rtt_sum",
		"Sum of measured latencies",
		nil, nil,
	)

	fpingRTTCountDesc = prometheus.NewDesc(
		"fping_rtt_count",
		"Number of measured latencies (successful probes)",
		nil, nil,
	)

	fpingRTTDesc = prometheus.NewDesc(
		"fping_rtt",
		"Summary of measured latencies",
		[]string{"quantile"}, nil,
	)
)

func (t *Target) AddMeasurements(m Measurements) {
	t.Lock()
	t.measurements = &m
	t.Unlock()
}

func (t *Target) Collect(ch chan<- prometheus.Metric) {
	t.Lock()
	defer t.Unlock()

	if t.measurements == nil {
		return
	}

	// fping_sent_count
	ch <- prometheus.MustNewConstMetric(
		fpingSentCountDesc,
		prometheus.GaugeValue,
		float64(t.measurements.GetSentCount()),
	)

	// fping_lost_count
	ch <- prometheus.MustNewConstMetric(
		fpingLostCountDesc,
		prometheus.GaugeValue,
		float64(t.measurements.GetLostCount()),
	)

	// sum
	ch <- prometheus.MustNewConstMetric(
		fpingRTTSumDesc,
		prometheus.GaugeValue,
		t.measurements.GetRTTSum(),
	)

	// count
	ch <- prometheus.MustNewConstMetric(
		fpingRTTCountDesc,
		prometheus.GaugeValue,
		float64(t.measurements.GetRTTCount()),
	)

	count := float64(len(t.measurements.rtt))
	for i := range t.measurements.rtt {
		quantile := fmt.Sprintf("%.3f", (float64(i)+1.0)/count)
		if t.measurements.lost[i] {
			ch <- prometheus.MustNewConstMetric(
				fpingRTTDesc,
				prometheus.GaugeValue,
				math.NaN(),
				quantile,
			)
		} else {
			ch <- prometheus.MustNewConstMetric(
				fpingRTTDesc,
				prometheus.GaugeValue,
				t.measurements.rtt[i],
				quantile,
			)
		}
	}
}

func (t *Target) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(t, ch)
}
