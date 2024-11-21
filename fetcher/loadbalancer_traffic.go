package fetcher

import (
	"math"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var _ Fetcher = &loadbalancerTraffic{}

// NewLoadbalancerTraffic creates a new fetcher that will collect pricing information on load balancer traffic.
func NewLoadbalancerTraffic(pricing *PriceProvider, additionalLabels ...string) Fetcher {
	return &loadbalancerTraffic{newBase(pricing, "loadbalancer_traffic", []string{"location", "type"}, additionalLabels...)}
}

type loadbalancerTraffic struct {
	*baseFetcher
}

func (loadbalancerTraffic loadbalancerTraffic) Run(client *hcloud.Client) error {
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
			parseAdditionalLabels(loadbalancerTraffic.additionalLabels, lb.Labels)...,
		)

		//nolint:gosec
		additionalTraffic := int64(lb.OutgoingTraffic) - int64(lb.IncludedTraffic)
		if additionalTraffic < 0 {
			loadbalancerTraffic.hourly.WithLabelValues(labels...).Set(0)
			loadbalancerTraffic.monthly.WithLabelValues(labels...).Set(0)
			break
		}

		lbTrafficPrice, err := loadbalancerTraffic.pricing.LoadBalancerTraffic(lb.LoadBalancerType, location.Name)
		if err != nil {
			return err
		}

		monthlyPrice := math.Ceil(float64(additionalTraffic)/sizeTB) * lbTrafficPrice
		hourlyPrice := pricingPerHour(monthlyPrice)

		loadbalancerTraffic.hourly.WithLabelValues(labels...).Set(hourlyPrice)
		loadbalancerTraffic.monthly.WithLabelValues(labels...).Set(monthlyPrice)
	}

	return nil
}
