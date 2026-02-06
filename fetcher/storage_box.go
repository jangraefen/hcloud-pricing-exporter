package fetcher

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var _ Fetcher = &storageBox{}

// NewStorageBox creates a new fetcher that will collect pricing information on storage boxes.
func NewStorageBox(pricing *PriceProvider, additionalLabels ...string) Fetcher {
	return &storageBox{newBase(pricing, "storage_box", []string{"location", "type"}, additionalLabels...)}
}

type storageBox struct {
	*baseFetcher
}

func (sb storageBox) Run(client *hcloud.Client) error {
	storageBoxes, err := client.StorageBox.All(ctx)
	if err != nil {
		return err
	}

	for _, s := range storageBoxes {
		location := s.Location.Name

		labels := append([]string{
			s.Name,
			location,
			s.StorageBoxType.Name,
		},
			parseAdditionalLabels(sb.additionalLabels, s.Labels)...,
		)

		hourly, monthly, err := sb.pricing.StorageBox(s.StorageBoxType, location)
		if err != nil {
			return err
		}

		sb.hourly.WithLabelValues(labels...).Set(hourly)
		sb.monthly.WithLabelValues(labels...).Set(monthly)
	}

	return nil
}
