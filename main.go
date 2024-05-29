package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
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
	hcloudAPIToken       string
	port                 uint
	fetchInterval        time.Duration
	additionalLabelsFlag string
	additionalLabels     []string
)

func handleFlags() {
	flag.StringVar(&hcloudAPIToken, "hcloud-token", "", "the token to authenticate against the HCloud API")
	flag.UintVar(&port, "port", defaultPort, "the port that the exporter exposes its data on")
	flag.DurationVar(&fetchInterval, "fetch-interval", defaultFetchInterval, "the interval between data fetching cycles")
	flag.StringVar(&additionalLabelsFlag, "additional-labels", "", "comma separated additional labels to parse for all metrics, e.g: 'service,environment,owner'")
	flag.Parse()

	if hcloudAPIToken == "" {
		if envHCloudAPIToken, present := os.LookupEnv("HCLOUD_TOKEN"); present {
			hcloudAPIToken = envHCloudAPIToken
		}
	}
	if hcloudAPIToken == "" {
		panic(fmt.Errorf("no API token for HCloud specified, but required"))
	}
	if strings.HasPrefix(hcloudAPIToken, "file:") {
		hcloudAPITokenBytes, err := os.ReadFile(strings.TrimPrefix(hcloudAPIToken, "file:"))
		if err != nil {
			panic(fmt.Errorf("failed to read HCLOUD_TOKEN from file: %s", err.Error()))
		}
		hcloudAPIToken = strings.TrimSpace(string(hcloudAPITokenBytes))
	}
	if len(hcloudAPIToken) != 64 {
		panic(fmt.Errorf("invalid API token for HCloud specified, must be 64 characters long"))
	}

	additionalLabelsFlag = strings.TrimSpace(strings.ReplaceAll(additionalLabelsFlag, " ", ""))
	additionalLabelsSlice := strings.Split(additionalLabelsFlag, ",")
	if len(additionalLabelsSlice) > 0 && additionalLabelsSlice[0] != "" {
		additionalLabels = additionalLabelsSlice
	}
}

func main() {
	handleFlags()

	client := hcloud.NewClient(hcloud.WithToken(hcloudAPIToken))
	priceRepository := &fetcher.PriceProvider{Client: client}

	fetchers := fetcher.Fetchers{
		fetcher.NewFloatingIP(priceRepository, additionalLabels...),
		fetcher.NewPrimaryIP(priceRepository, additionalLabels...),
		fetcher.NewLoadbalancer(priceRepository, additionalLabels...),
		fetcher.NewLoadbalancerTraffic(priceRepository, additionalLabels...),
		fetcher.NewServer(priceRepository, additionalLabels...),
		fetcher.NewServerBackup(priceRepository, additionalLabels...),
		fetcher.NewServerTraffic(priceRepository, additionalLabels...),
		fetcher.NewSnapshot(priceRepository, additionalLabels...),
		fetcher.NewVolume(priceRepository, additionalLabels...),
	}

	fetchers.MustRun(client)
	scheduler.RunTaskAtInterval(func() { fetchers.MustRun(client) }, fetchInterval, 0)
	scheduler.RunTaskAtInterval(priceRepository.Sync, 10*fetchInterval, 10*fetchInterval)

	registry := prometheus.NewRegistry()
	fetchers.RegisterCollectors(registry)

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	log.Printf("Listening on: http://0.0.0.0:%d\n", port)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
