package fetcher

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
)

var _ Fetcher = &floatingIP{}

// NewPrimaryIP creates a new fetcher that will collect pricing information on primary IPs.
func NewPrimaryIP(pricing *PriceProvider) Fetcher {
	return &primaryIP{newBase(pricing, "primaryip", "datacenter", "type")}
}

type primaryIP struct {
	*baseFetcher
}

func (primaryIP primaryIP) Run(client *hcloud.Client) error {
	primaryIPs, _, err := client.PrimaryIP.List(ctx, hcloud.PrimaryIPListOpts{})
	if err != nil {
		return err
	}

	for _, p := range primaryIPs {
		datacenter := p.Datacenter

		hourlyPrice, monthlyPrice, err := primaryIP.pricing.PrimaryIP(p.Type, datacenter.Location.Name)
		if err != nil {
			return err
		}

		primaryIP.hourly.WithLabelValues(p.Name, datacenter.Name, string(p.Type)).Set(hourlyPrice)
		primaryIP.monthly.WithLabelValues(p.Name, datacenter.Name, string(p.Type)).Set(monthlyPrice)
	}

	return nil
}
