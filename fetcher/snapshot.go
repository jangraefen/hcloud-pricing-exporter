package fetcher

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var _ Fetcher = &snapshot{}

// NewSnapshot creates a new fetcher that will collect pricing information on server snapshots.
func NewSnapshot(pricing *PriceProvider, additionalLabels ...string) Fetcher {
	return &snapshot{newBase(pricing, "snapshot", nil, additionalLabels...)}
}

type snapshot struct {
	*baseFetcher
}

func (snapshot snapshot) Run(client *hcloud.Client) error {
	images, err := client.Image.All(ctx)
	if err != nil {
		return err
	}

	for _, i := range images {
		if i.Type == "snapshot" {
			monthlyPrice := float64(i.ImageSize) * snapshot.pricing.Image()
			hourlyPrice := pricingPerHour(monthlyPrice)

			labels := append([]string{
				i.Name,
			},
				parseAdditionalLabels(snapshot.additionalLabels, i.Labels)...,
			)

			snapshot.hourly.WithLabelValues(labels...).Set(hourlyPrice)
			snapshot.monthly.WithLabelValues(labels...).Set(monthlyPrice)
		}
	}

	return nil
}
