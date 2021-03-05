package fetcher

import (
	"context"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ctx = context.Background()
)

type Fetcher interface {
	GetHourly() prometheus.Collector
	GetMonthly() prometheus.Collector
	Run(*hcloud.Client) error
}

type baseFetcher struct {
	hourly  *prometheus.GaugeVec
	monthly *prometheus.GaugeVec
}

func (fetcher baseFetcher) GetHourly() prometheus.Collector {
	return fetcher.hourly
}

func (fetcher baseFetcher) GetMonthly() prometheus.Collector {
	return fetcher.monthly
}

func new(resource string, additionalLabels ...string) *baseFetcher {
	labels := []string{"name"}
	labels = append(labels, additionalLabels...)

	hourlyGaugeOpts := prometheus.GaugeOpts{
		Namespace: "hcloud",
		Subsystem: "pricing",
		Name:      fmt.Sprintf("%s_hourly", resource),
		Help:      fmt.Sprintf("The cost of the resource %s per hour", resource),
	}
	monthlyGaugeOpts := prometheus.GaugeOpts{
		Namespace: "hcloud",
		Subsystem: "pricing",
		Name:      fmt.Sprintf("%s_monthly", resource),
		Help:      fmt.Sprintf("The cost of the resource %s per month", resource),
	}

	return &baseFetcher{
		hourly:  prometheus.NewGaugeVec(hourlyGaugeOpts, labels),
		monthly: prometheus.NewGaugeVec(monthlyGaugeOpts, labels),
	}
}
