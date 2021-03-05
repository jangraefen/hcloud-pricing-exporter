package fetcher

import (
	"context"
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
func (provider *PriceProvider) FloatingIP() float64 {
	provider.pricingLock.RLock()
	defer provider.pricingLock.RUnlock()

	return parsePrice(provider.getPricing().FloatingIP.Monthly.Gross)
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
	return 0.0476
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
