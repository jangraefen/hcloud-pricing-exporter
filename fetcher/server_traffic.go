package fetcher

import (
	"math"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

const (
	trafficPrice = 1.19
)

var _ Fetcher = &serverTraffic{}

func NewServerTraffic() Fetcher {
	return &serverTraffic{new("server_traffic", "location", "type")}
}

type serverTraffic struct {
	*baseFetcher
}

func (serverTraffic serverTraffic) Run(client *hcloud.Client) error {
	servers, _, err := client.Server.List(ctx, hcloud.ServerListOpts{})
	if err != nil {
		return err
	}

	for _, s := range servers {
		location := s.Datacenter.Location

		additionalTraffic := s.OutgoingTraffic - s.IncludedTraffic
		if additionalTraffic < 0 {
			serverTraffic.hourly.WithLabelValues(s.Name, location.Name, s.ServerType.Name).Set(0)
			serverTraffic.monthly.WithLabelValues(s.Name, location.Name, s.ServerType.Name).Set(0)
			break
		}

		monthlyPrice := math.Ceil(float64(additionalTraffic)/sizeTB) * trafficPrice
		hourlyPrice := pricingPerHour(monthlyPrice)

		serverTraffic.hourly.WithLabelValues(s.Name, location.Name, s.ServerType.Name).Set(hourlyPrice)
		serverTraffic.monthly.WithLabelValues(s.Name, location.Name, s.ServerType.Name).Set(monthlyPrice)
	}

	return nil
}
