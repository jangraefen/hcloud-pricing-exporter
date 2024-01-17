package fetcher

import (
	"fmt"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var _ Fetcher = &loadBalancer{}

// NewLoadbalancer creates a new fetcher that will collect pricing information on load balancers.
func NewLoadbalancer(pricing *PriceProvider, additionalLabels ...string) Fetcher {
	return &loadBalancer{newBase(pricing, "loadbalancer", []string{"location", "type"}, additionalLabels...)}
}

type loadBalancer struct {
	*baseFetcher
}

func (loadBalancer loadBalancer) Run(client *hcloud.Client) error {
	loadBalancers, err := client.LoadBalancer.All(ctx)
	if err != nil {
		return err
	}

	for _, lb := range loadBalancers {
		location := lb.Location

		labels := append([]string{
			lb.Name,
			location.Name,
			lb.LoadBalancerType.Name,
		},
			parseAdditionalLabels(loadBalancer.additionalLabels, lb.Labels)...,
		)

		pricing, err := findLBPricing(location, lb.LoadBalancerType.Pricings)
		if err != nil {
			return err
		}

		parseToGauge(loadBalancer.hourly.WithLabelValues(labels...), pricing.Hourly.Gross)
		parseToGauge(loadBalancer.monthly.WithLabelValues(labels...), pricing.Monthly.Gross)
	}

	return nil
}

func findLBPricing(location *hcloud.Location, pricings []hcloud.LoadBalancerTypeLocationPricing) (*hcloud.LoadBalancerTypeLocationPricing, error) {
	for _, pricing := range pricings {
		if pricing.Location.Name == location.Name {
			return &pricing, nil
		}
	}

	return nil, fmt.Errorf("no load balancer pricing found for location %s", location.Name)
}
