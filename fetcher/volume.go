package fetcher

import (
	"strconv"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var _ Fetcher = &volume{}

// NewVolume creates a new fetcher that will collect pricing information on volumes.
func NewVolume(pricing *PriceProvider, additionalLabels ...string) Fetcher {
	return &volume{newBase(pricing, "volume", []string{"location", "bytes"}, additionalLabels...)}
}

type volume struct {
	*baseFetcher
}

func (volume volume) Run(client *hcloud.Client) error {
	volumes, err := client.Volume.All(ctx)
	if err != nil {
		return err
	}

	for _, v := range volumes {
		monthlyPrice := float64(v.Size) * volume.pricing.Volume()
		hourlyPrice := pricingPerHour(monthlyPrice)

		labels := append([]string{
			v.Name,
			v.Location.Name,
			strconv.Itoa(v.Size),
		},
			parseAdditionalLabels(volume.additionalLabels, v.Labels)...,
		)

		volume.hourly.WithLabelValues(labels...).Set(hourlyPrice)
		volume.monthly.WithLabelValues(labels...).Set(monthlyPrice)
	}

	return nil
}
