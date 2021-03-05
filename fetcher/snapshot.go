package fetcher

import (
	"math"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

const (
	imagePrice = 0.0119
)

var _ Fetcher = &snapshot{}

// NewSnapshot creates a new fetcher that will collect pricing information on server snapshots.
func NewSnapshot() Fetcher {
	return &snapshot{newBase("snapshot")}
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
			monthlyPrice := math.Ceil(float64(i.ImageSize)/sizeGB) * imagePrice
			hourlyPrice := pricingPerHour(monthlyPrice)

			snapshot.hourly.WithLabelValues(i.Name).Set(hourlyPrice)
			snapshot.monthly.WithLabelValues(i.Name).Set(monthlyPrice)
		}
	}

	return nil
}
