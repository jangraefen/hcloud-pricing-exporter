package fetcher

import (
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

var _ Fetcher = &server{}

func NewServer() Fetcher {
	return &server{new("server", "location")}
}

type server struct {
	*baseFetcher
}

func (server server) Run(client *hcloud.Client) error {
	servers, _, err := client.Server.List(ctx, hcloud.ServerListOpts{})
	if err != nil {
		return err
	}

	for _, s := range servers {
		location := s.Datacenter.Location

		pricing, err := findServerPricing(location, s.ServerType.Pricings)
		if err != nil {
			return err
		}

		parseToGauge(server.hourly.WithLabelValues(s.Name, location.Name, s.ServerType.Name), pricing.Hourly.Gross)
		parseToGauge(server.monthly.WithLabelValues(s.Name, location.Name, s.ServerType.Name), pricing.Monthly.Gross)
	}

	return nil
}

func findServerPricing(location *hcloud.Location, pricings []hcloud.ServerTypeLocationPricing) (*hcloud.ServerTypeLocationPricing, error) {
	for _, pricing := range pricings {
		if pricing.Location.ID == location.ID {
			return &pricing, nil
		}
	}

	return nil, fmt.Errorf("no pricing found for location %s", location.Name)
}
