package fetcher

import (
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

var _ Fetcher = &volume{}

// NewVolume creates a new fetcher that will collect pricing information on volumes.
func NewVolume(pricing *PriceProvider) Fetcher {
	return &volume{newBase(pricing, "volume", "location", "bytes")}
}

type volume struct {
	*baseFetcher
}

func (volume volume) Run(client *hcloud.Client) error {
	volumes, _, err := client.Volume.List(ctx, hcloud.VolumeListOpts{})
	if err != nil {
		return err
	}

	for _, v := range volumes {
		monthlyPrice := float64(v.Size) * volume.pricing.Volume()
		hourlyPrice := pricingPerHour(monthlyPrice)

		volume.hourly.WithLabelValues(v.Name, v.Location.Name, strconv.Itoa(v.Size)).Set(hourlyPrice)
		volume.monthly.WithLabelValues(v.Name, v.Location.Name, strconv.Itoa(v.Size)).Set(monthlyPrice)
	}

	return nil
}
