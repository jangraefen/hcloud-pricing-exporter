package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/jangraefen/hcloud-pricing-exporter/fetcher"
	"github.com/jtaczanowski/go-scheduler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultPort          = 8080
	defaultFetchInterval = 1 * time.Minute
)

var (
	hcloudAPIToken string
	port           uint
	fetchInterval  time.Duration
)

func handleFlags() {
	flag.StringVar(&hcloudAPIToken, "hcloud-token", "", "the token to authenticate against the HCloud API")
	flag.UintVar(&port, "port", defaultPort, "the port that the exporter exposes its data on")
	flag.DurationVar(&fetchInterval, "fetch-interval", defaultFetchInterval, "the interval between data fetching cycles")
	flag.Parse()

	if hcloudAPIToken == "" {
		if envHCloudAPIToken, present := os.LookupEnv("HCLOUD_TOKEN"); present {
			hcloudAPIToken = envHCloudAPIToken
		}
	}
	if hcloudAPIToken == "" {
		panic(fmt.Errorf("no API token for HCloud specified, but required"))
	}
}

func main() {
	handleFlags()

	client := hcloud.NewClient(hcloud.WithToken(hcloudAPIToken))
	priceRepository := &fetcher.PriceProvider{Client: client}

	fetchers := fetcher.Fetchers{
		fetcher.NewFloatingIP(priceRepository),
		fetcher.NewPrimaryIP(priceRepository),
		fetcher.NewLoadbalancer(priceRepository),
		fetcher.NewLoadbalancerTraffic(priceRepository),
		fetcher.NewServer(priceRepository),
		fetcher.NewServerBackup(priceRepository),
		fetcher.NewServerTraffic(priceRepository),
		fetcher.NewSnapshot(priceRepository),
		fetcher.NewVolume(priceRepository),
	}

	fetchers.MustRun(client)
	scheduler.RunTaskAtInterval(func() { fetchers.MustRun(client) }, fetchInterval, 0)
	scheduler.RunTaskAtInterval(priceRepository.Sync, 10*fetchInterval, 10*fetchInterval)

	registry := prometheus.NewRegistry()
	fetchers.RegisterCollectors(registry)

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
