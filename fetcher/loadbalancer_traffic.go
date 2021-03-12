package fetcher

import (
	"math"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

var _ Fetcher = &loadbalancerTraffic{}

// NewLoadbalancerTraffic creates a new fetcher that will collect pricing information on load balancer traffic.
func NewLoadbalancerTraffic(pricing *PriceProvider) Fetcher {
	return &loadbalancerTraffic{newBase(pricing, "loadbalancer_traffic", "location", "type")}
}

type loadbalancerTraffic struct {
	*baseFetcher
}

func (loadbalancerTraffic loadbalancerTraffic) Run(client *hcloud.Client) error {
	loadBalancers, _, err := client.LoadBalancer.List(ctx, hcloud.LoadBalancerListOpts{})
	if err != nil {
		return err
	}

	for _, lb := range loadBalancers {
		location := lb.Location

		additionalTraffic := int(lb.OutgoingTraffic) - int(lb.IncludedTraffic)
		if additionalTraffic < 0 {
			loadbalancerTraffic.hourly.WithLabelValues(lb.Name, location.Name, lb.LoadBalancerType.Name).Set(0)
			loadbalancerTraffic.monthly.WithLabelValues(lb.Name, location.Name, lb.LoadBalancerType.Name).Set(0)
			break
		}

		monthlyPrice := math.Ceil(float64(additionalTraffic)/sizeTB) * loadbalancerTraffic.pricing.Traffic()
		hourlyPrice := pricingPerHour(monthlyPrice)

		loadbalancerTraffic.hourly.WithLabelValues(lb.Name, location.Name, lb.LoadBalancerType.Name).Set(hourlyPrice)
		loadbalancerTraffic.monthly.WithLabelValues(lb.Name, location.Name, lb.LoadBalancerType.Name).Set(monthlyPrice)
	}

	return nil
}
