package fetcher

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
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

func parseAdditionalLabels(additionalLabels []string, labels map[string]string) (result []string) {
	for _, al := range additionalLabels {
		result = append(result, findLabel(labels, al))
	}
	return result
}

func findLabel(labels map[string]string, label string) string {
	for k, v := range labels {
		if k == label {
			return v
		}
	}
	return ""
}
