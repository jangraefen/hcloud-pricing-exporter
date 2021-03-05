package fetcher

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
)

const (
	floatingIPPrice = float64(1.19)
)

var _ Fetcher = &floatingIP{}

func NewFloatingIP() Fetcher {
	return &floatingIP{new("floatingip", "location")}
}

type floatingIP struct {
	*baseFetcher
}

func (floatingIP floatingIP) Run(client *hcloud.Client) error {
	floatingIPs, _, err := client.FloatingIP.List(ctx, hcloud.FloatingIPListOpts{})
	if err != nil {
		return err
	}

	for _, f := range floatingIPs {
		location := f.HomeLocation

		monthlyPrice := floatingIPPrice
		hourlyPrice := pricingPerHour(monthlyPrice)

		floatingIP.hourly.WithLabelValues(f.Name, location.Name).Set(hourlyPrice)
		floatingIP.monthly.WithLabelValues(f.Name, location.Name).Set(monthlyPrice)
	}

	return nil
}
