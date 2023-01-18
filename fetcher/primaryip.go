package fetcher

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
)

var _ Fetcher = &floatingIP{}

// NewPrimaryIP creates a new fetcher that will collect pricing information on primary IPs.
func NewPrimaryIP(pricing *PriceProvider, additionalLabels ...string) Fetcher {
	labels := append([]string{"datacenter", "type"}, additionalLabels...)
	return &primaryIP{newBase(pricing, "primaryip", labels...), additionalLabels}
}

type primaryIP struct {
	*baseFetcher
	additionalLabels []string
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

		labels := append([]string{
			p.Name,
			datacenter.Name,
			string(p.Type),
		},
			parseAdditionalLabels(primaryIP.additionalLabels, p.Labels)...,
		)

		primaryIP.hourly.WithLabelValues(labels...).Set(hourlyPrice)
		primaryIP.monthly.WithLabelValues(labels...).Set(monthlyPrice)
	}

	return nil
}
