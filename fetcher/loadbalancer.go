package fetcher

import (
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

var _ Fetcher = &loadBalancer{}

// NewLoadbalancer creates a new fetcher that will collect pricing information on load balancers.
func NewLoadbalancer() Fetcher {
	return &loadBalancer{newBase("loadbalancer", "location", "type")}
}

type loadBalancer struct {
	*baseFetcher
}

func (loadBalancer loadBalancer) Run(client *hcloud.Client) error {
	loadBalancers, _, err := client.LoadBalancer.List(ctx, hcloud.LoadBalancerListOpts{})
	if err != nil {
		return err
	}

	for _, lb := range loadBalancers {
		location := lb.Location

		pricing, err := findLBPricing(location, lb.LoadBalancerType.Pricings)
		if err != nil {
			return err
		}

		parseToGauge(loadBalancer.hourly.WithLabelValues(lb.Name, location.Name, lb.LoadBalancerType.Name), pricing.Hourly.Gross)
		parseToGauge(loadBalancer.monthly.WithLabelValues(lb.Name, location.Name, lb.LoadBalancerType.Name), pricing.Monthly.Gross)
	}

	return nil
}

func findLBPricing(location *hcloud.Location, pricings []hcloud.LoadBalancerTypeLocationPricing) (*hcloud.LoadBalancerTypeLocationPricing, error) {
	for _, pricing := range pricings {
		if pricing.Location.ID == location.ID {
			return &pricing, nil
		}
	}

	return nil, fmt.Errorf("no pricing found for location %s", location.Name)
}
