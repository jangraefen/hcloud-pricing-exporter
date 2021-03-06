package fetcher

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
)

var _ Fetcher = &snapshot{}

// NewSnapshot creates a new fetcher that will collect pricing information on server snapshots.
func NewSnapshot(pricing *PriceProvider) Fetcher {
	return &snapshot{newBase(pricing, "snapshot")}
}

type snapshot struct {
	*baseFetcher
}

func (snapshot snapshot) Run(client *hcloud.Client) error {
	images, _, err := client.Image.List(ctx, hcloud.ImageListOpts{})
	if err != nil {
		return err
	}

	for _, i := range images {
		if i.Type == "snapshot" {
			monthlyPrice := float64(i.ImageSize) * snapshot.pricing.Image()
			hourlyPrice := pricingPerHour(monthlyPrice)

			snapshot.hourly.WithLabelValues(i.Name).Set(hourlyPrice)
			snapshot.monthly.WithLabelValues(i.Name).Set(monthlyPrice)
		}
	}

	return nil
}
