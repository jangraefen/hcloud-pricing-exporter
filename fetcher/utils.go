package fetcher

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	sizeGB = 1 << (10 * 3)
	sizeTB = 1 << (10 * 4)
)

func daysInMonth() int {
	now := time.Now()

	switch now.Month() {
	case time.April, time.June, time.September, time.November:
		return 30
	case time.February:
		year := now.Year()
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			return 29
		}
		return 28
	default:
		return 31
	}
}

func pricingPerHour(monthlyPrice float64) float64 {
	return monthlyPrice / float64(daysInMonth()) / 24
}

func parseToGauge(gauge prometheus.Gauge, value string) {
	parsed, err := strconv.ParseFloat(value, 32)
	if err != nil {
		panic(err)
	}
	gauge.Set(parsed)
}
