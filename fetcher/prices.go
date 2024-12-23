package fetcher

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// PriceProvider provides easy access to current HCloud prices.
type PriceProvider struct {
	Client      *hcloud.Client
	pricing     *hcloud.Pricing
	pricingLock sync.RWMutex
}

// FloatingIP returns the current price for a floating IP per month.
func (provider *PriceProvider) FloatingIP(ipType hcloud.FloatingIPType, location string) float64 {
	provider.pricingLock.RLock()
	defer provider.pricingLock.RUnlock()

	for _, byType := range provider.getPricing().FloatingIPs {
		if byType.Type == ipType {
			for _, pricing := range byType.Pricings {
				if pricing.Location.Name == location {
					return parsePrice(pricing.Monthly.Gross)
				}
			}
		}
	}

	// If the pricing can not be determined by the type and location, we just return 0.00
	return 0.0
}

// PrimaryIP returns the current price for a primary IP per hour and month.
func (provider *PriceProvider) PrimaryIP(ipType hcloud.PrimaryIPType, location string) (hourly, monthly float64, err error) {
	provider.pricingLock.RLock()
	defer provider.pricingLock.RUnlock()

	// v6 pricing is not defined by the API
	if string(ipType) == "ipv6" {
		return 0, 0, nil
	}

	for _, byType := range provider.getPricing().PrimaryIPs {
		if byType.Type == string(ipType) {
			for _, pricing := range byType.Pricings {
				if pricing.Location == location {
					return parsePrice(pricing.Hourly.Gross), parsePrice(pricing.Monthly.Gross), nil
				}
			}
		}
	}

	return 0, 0, fmt.Errorf("no primary IP pricing found for location %s", location)
}

// Image returns the current price for an image per GB per month.
func (provider *PriceProvider) Image() float64 {
	provider.pricingLock.RLock()
	defer provider.pricingLock.RUnlock()

	return parsePrice(provider.getPricing().Image.PerGBMonth.Gross)
}

// ServerTraffic returns the current price for a TB of extra traffic per month.
func (provider *PriceProvider) ServerTraffic(serverType *hcloud.ServerType, location string) (gross float64, err error) {
	provider.pricingLock.RLock()
	defer provider.pricingLock.RUnlock()

	for _, byType := range provider.getPricing().ServerTypes {
		if byType.ServerType.ID == serverType.ID {
			for _, price := range byType.Pricings {
				if price.Location.Name == location {
					return parsePrice(price.PerTBTraffic.Gross), nil
				}
			}
		}
	}

	return 0, fmt.Errorf("no traffic pricing found for server type %s and location %s", serverType.Name, location)
}

// LoadBalancerTraffic returns the current price for a TB of extra traffic per month.
func (provider *PriceProvider) LoadBalancerTraffic(loadBalancerType *hcloud.LoadBalancerType, location string) (gross float64, err error) {
	provider.pricingLock.RLock()
	defer provider.pricingLock.RUnlock()

	for _, byType := range provider.getPricing().LoadBalancerTypes {
		if byType.LoadBalancerType.ID == loadBalancerType.ID {
			for _, price := range byType.Pricings {
				if price.Location.Name == location {
					return parsePrice(price.PerTBTraffic.Gross), nil
				}
			}
		}
	}

	return 0, fmt.Errorf("no traffic pricing found for load balancer type %s and location %s", loadBalancerType.Name, location)
}

// ServerBackup returns the percentage of base price increase for server backups per month.
func (provider *PriceProvider) ServerBackup() float64 {
	provider.pricingLock.RLock()
	defer provider.pricingLock.RUnlock()

	return parsePrice(provider.getPricing().ServerBackup.Percentage)
}

// Volume returns the current price for a volume per GB per month.
func (provider *PriceProvider) Volume() float64 {
	provider.pricingLock.RLock()
	defer provider.pricingLock.RUnlock()

	return parsePrice(provider.getPricing().Volume.PerGBMonthly.Gross)
}

// Sync forces the provider to re-fetch prices from the HCloud API.
func (provider *PriceProvider) Sync() {
	provider.pricingLock.Lock()
	defer provider.pricingLock.Unlock()

	provider.pricing = nil
}

func (provider *PriceProvider) getPricing() *hcloud.Pricing {
	if provider.pricing == nil {
		pricing, _, err := provider.Client.Pricing.Get(context.Background())
		if err != nil {
			panic(err)
		}

		provider.pricing = &pricing
	}

	return provider.pricing
}

func parsePrice(rawPrice string) float64 {
	if price, err := strconv.ParseFloat(rawPrice, 32); err == nil {
		return price
	}

	return 0
}
