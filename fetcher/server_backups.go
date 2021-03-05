package fetcher

import (
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

const (
	backupPriceMultiplier = 0.2
)

var _ Fetcher = &server{}

// NewServerBackup creates a new fetcher that will collect pricing information on server backups.
func NewServerBackup() Fetcher {
	return &serverBackup{newBase("server_backup", "location", "type")}
}

type serverBackup struct {
	*baseFetcher
}

func (serverBackup serverBackup) Run(client *hcloud.Client) error {
	servers, _, err := client.Server.List(ctx, hcloud.ServerListOpts{})
	if err != nil {
		return err
	}

	for _, s := range servers {
		if s.BackupWindow != "" {
			location := s.Datacenter.Location

			serverPrice, err := findServerPricing(location, s.ServerType.Pricings)
			if err != nil {
				return err
			}

			hourlyPrice := toBackupPrice(serverPrice.Hourly.Gross)
			monthlyPrice := toBackupPrice(serverPrice.Monthly.Gross)

			serverBackup.hourly.WithLabelValues(s.Name, location.Name, s.ServerType.Name).Set(hourlyPrice)
			serverBackup.monthly.WithLabelValues(s.Name, location.Name, s.ServerType.Name).Set(monthlyPrice)
		}
	}

	return nil
}

func toBackupPrice(rawServerPrice string) float64 {
	serverPrice, err := strconv.Atoi(rawServerPrice)
	if err != nil {
		return 0
	}

	return float64(serverPrice) * backupPriceMultiplier
}
