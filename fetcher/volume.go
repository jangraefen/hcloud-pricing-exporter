package fetcher

import (
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

var _ Fetcher = &volume{}

// NewVolume creates a new fetcher that will collect pricing information on volumes.
func NewVolume(pricing *PriceProvider, additionalLabels ...string) Fetcher {
	labels := append([]string{"location", "bytes"}, additionalLabels...)
	return &volume{newBase(pricing, "volume", labels...), additionalLabels}
}

type volume struct {
	*baseFetcher
	additionalLabels []string
}

func (volume volume) Run(client *hcloud.Client) error {
	volumes, _, err := client.Volume.List(ctx, hcloud.VolumeListOpts{})
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
