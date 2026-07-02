package fetcher

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var _ Fetcher = &floatingIP{}

// NewPrimaryIP creates a new fetcher that will collect pricing information on primary IPs.
func NewPrimaryIP(pricing *PriceProvider, additionalLabels ...string) Fetcher {
	return &primaryIP{newBase(pricing, "primaryip", []string{"datacenter", "type"}, additionalLabels...)}
}

type primaryIP struct {
	*baseFetcher
}

func (primaryIP primaryIP) Run(client *hcloud.Client) error {
	primaryIPs, err := client.PrimaryIP.All(ctx)
	if err != nil {
		return err
	}

	for _, p := range primaryIPs {
		location := p.Location

		hourlyPrice, monthlyPrice, err := primaryIP.pricing.PrimaryIP(p.Type, location.Name)
		if err != nil {
			return err
		}

		labels := append(
			[]string{
				p.Name,
				location.Name, // Using location name as datacenter/location label since datacenter is deprecated
				string(p.Type),
			},
			parseAdditionalLabels(primaryIP.additionalLabels, p.Labels)...,
		)

		primaryIP.hourly.WithLabelValues(labels...).Set(hourlyPrice)
		primaryIP.monthly.WithLabelValues(labels...).Set(monthlyPrice)
	}

	return nil
}
