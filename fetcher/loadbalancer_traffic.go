package fetcher

import (
	"math"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

var _ Fetcher = &loadbalancerTraffic{}

// NewLoadbalancerTraffic creates a new fetcher that will collect pricing information on load balancer traffic.
func NewLoadbalancerTraffic(pricing *PriceProvider, additionalLabels ...string) Fetcher {
	labels := append([]string{"location", "type"}, additionalLabels...)
	return &loadbalancerTraffic{newBase(pricing, "loadbalancer_traffic", labels...), additionalLabels}
}

type loadbalancerTraffic struct {
	*baseFetcher
	additionalLabels []string
}

func (loadbalancerTraffic loadbalancerTraffic) Run(client *hcloud.Client) error {
	loadBalancers, _, err := client.LoadBalancer.List(ctx, hcloud.LoadBalancerListOpts{})
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
			parseAdditionalLabels(loadbalancerTraffic.additionalLabels, lb.Labels)...,
		)

		additionalTraffic := int(lb.OutgoingTraffic) - int(lb.IncludedTraffic)
		if additionalTraffic < 0 {
			loadbalancerTraffic.hourly.WithLabelValues(labels...).Set(0)
			loadbalancerTraffic.monthly.WithLabelValues(labels...).Set(0)
			break
		}

		monthlyPrice := math.Ceil(float64(additionalTraffic)/sizeTB) * loadbalancerTraffic.pricing.Traffic()
		hourlyPrice := pricingPerHour(monthlyPrice)

		loadbalancerTraffic.hourly.WithLabelValues(labels...).Set(hourlyPrice)
		loadbalancerTraffic.monthly.WithLabelValues(labels...).Set(monthlyPrice)
	}

	return nil
}
