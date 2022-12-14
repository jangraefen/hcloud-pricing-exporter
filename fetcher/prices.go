package fetcher

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/hetznercloud/hcloud-go/hcloud"
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

	// If the pricing can not be determined by the type and location, we fall back to the old pricing
	return parsePrice(provider.getPricing().FloatingIP.Monthly.Gross)
}

// PrimaryIP returns the current price for a primary IP per hour and month.
func (provider *PriceProvider) PrimaryIP(ipType hcloud.PrimaryIPType, datacenter string) (hourly, monthly float64, err error) {
	provider.pricingLock.RLock()
	defer provider.pricingLock.RUnlock()

	for _, byType := range provider.getPricing().PrimaryIPs {
		if byType.Type == string(ipType) {
			for _, pricing := range byType.Pricings {
				if pricing.Datacenter == datacenter {
					return parsePrice(pricing.Hourly.Gross), parsePrice(pricing.Monthly.Gross), nil
				}
			}
		}
	}

	return 0, 0, fmt.Errorf("no primary IP pricing found for datacenter %s", datacenter)
}

// Image returns the current price for an image per GB per month.
func (provider *PriceProvider) Image() float64 {
	provider.pricingLock.RLock()
	defer provider.pricingLock.RUnlock()

	return parsePrice(provider.getPricing().Image.PerGBMonth.Gross)
}

// Traffic returns the current price for a TB of extra traffic per month.
func (provider *PriceProvider) Traffic() float64 {
	provider.pricingLock.RLock()
	defer provider.pricingLock.RUnlock()

	return parsePrice(provider.getPricing().Traffic.PerTB.Gross)
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
