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
	floatingIP    = fetcher.NewFloatingIP()
	loadBalancer  = fetcher.NewLoadbalancer()
	server        = fetcher.NewServer()
	serverBackup  = fetcher.NewServerBackup()
	serverTraffic = fetcher.NewServerTraffic()
	snapshot      = fetcher.NewSnapshot()
	volume        = fetcher.NewVolume()
)

func toScheduler(client *hcloud.Client, f func(*hcloud.Client) error) func() {
	return func() {
		if err := f(client); err != nil {
			panic(err)
		}
	}
}

func main() {
	var (
		hcloudAPIToken string
		port           uint
		fetchInterval  time.Duration
	)

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

	client := hcloud.NewClient(hcloud.WithToken(hcloudAPIToken))
	scheduler.RunTaskAtInterval(toScheduler(client, floatingIP.Run), fetchInterval, 0)
	scheduler.RunTaskAtInterval(toScheduler(client, loadBalancer.Run), fetchInterval, 0)
	scheduler.RunTaskAtInterval(toScheduler(client, server.Run), fetchInterval, 0)
	scheduler.RunTaskAtInterval(toScheduler(client, serverBackup.Run), fetchInterval, 0)
	scheduler.RunTaskAtInterval(toScheduler(client, serverTraffic.Run), fetchInterval, 0)
	scheduler.RunTaskAtInterval(toScheduler(client, snapshot.Run), fetchInterval, 0)
	scheduler.RunTaskAtInterval(toScheduler(client, volume.Run), fetchInterval, 0)

	registry := prometheus.NewRegistry()
	registry.MustRegister(
		floatingIP.GetHourly(),
		floatingIP.GetMonthly(),
		loadBalancer.GetHourly(),
		loadBalancer.GetMonthly(),
		server.GetHourly(),
		server.GetMonthly(),
		serverBackup.GetHourly(),
		serverBackup.GetMonthly(),
		serverTraffic.GetHourly(),
		serverTraffic.GetMonthly(),
		snapshot.GetHourly(),
		snapshot.GetMonthly(),
		volume.GetHourly(),
		volume.GetMonthly(),
	)

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
