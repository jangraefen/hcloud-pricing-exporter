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

// Fetcher defines a common interface for types that fetch pricing data from the HCloud API.
type Fetcher interface {
	// GetHourly returns the prometheus collector that collects pricing data for hourly expenses.
	GetHourly() prometheus.Collector
	// GetMonthly returns the prometheus collector that collects pricing data for monthly expenses.
	GetMonthly() prometheus.Collector
	// Run executes a new data fetching cycle and updates the prometheus exposed collectors.
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

func newBase(resource string, additionalLabels ...string) *baseFetcher {
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
