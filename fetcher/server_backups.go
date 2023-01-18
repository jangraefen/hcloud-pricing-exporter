package fetcher

import (
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

var _ Fetcher = &serverBackup{}

// NewServerBackup creates a new fetcher that will collect pricing information on server backups.
func NewServerBackup(pricing *PriceProvider, additionalLabels ...string) Fetcher {
	labels := append([]string{"location", "type"}, additionalLabels...)
	return &serverBackup{newBase(pricing, "server_backup", labels...), additionalLabels}
}

type serverBackup struct {
	*baseFetcher
	additionalLabels []string
}

func (serverBackup serverBackup) Run(client *hcloud.Client) error {
	servers, _, err := client.Server.List(ctx, hcloud.ServerListOpts{})
	if err != nil {
		return err
	}

	for _, s := range servers {
		location := s.Datacenter.Location

		labels := append([]string{
			s.Name,
			location.Name,
			s.ServerType.Name,
		},
			parseAdditionalLabels(serverBackup.additionalLabels, s.Labels)...,
		)

		if s.BackupWindow != "" {
			serverPrice, err := findServerPricing(location, s.ServerType.Pricings)
			if err != nil {
				return err
			}

			hourlyPrice := serverBackup.toBackupPrice(serverPrice.Hourly.Gross)
			monthlyPrice := serverBackup.toBackupPrice(serverPrice.Monthly.Gross)

			serverBackup.hourly.WithLabelValues(labels...).Set(hourlyPrice)
			serverBackup.monthly.WithLabelValues(labels...).Set(monthlyPrice)
		} else {
			serverBackup.hourly.WithLabelValues(labels...).Set(0)
			serverBackup.monthly.WithLabelValues(labels...).Set(0)
		}
	}

	return nil
}

func (serverBackup serverBackup) toBackupPrice(rawServerPrice string) float64 {
	serverPrice, err := strconv.ParseFloat(rawServerPrice, 32)
	if err != nil {
		return 0
	}

	return serverPrice * (serverBackup.pricing.ServerBackup() / 100)
}
