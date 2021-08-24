package fetcher

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
)

var _ Fetcher = &floatingIP{}

// NewFloatingIP creates a new fetcher that will collect pricing information on floating IPs.
func NewFloatingIP(pricing *PriceProvider) Fetcher {
	return &floatingIP{newBase(pricing, "floatingip", "location")}
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

		monthlyPrice := floatingIP.pricing.FloatingIP(f.Type, location.Name)
		hourlyPrice := pricingPerHour(monthlyPrice)

		floatingIP.hourly.WithLabelValues(f.Name, location.Name).Set(hourlyPrice)
		floatingIP.monthly.WithLabelValues(f.Name, location.Name).Set(monthlyPrice)
	}

	return nil
}
