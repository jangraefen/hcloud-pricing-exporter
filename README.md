# hcloud-pricing-exporter

[![Build Status](https://img.shields.io/github/actions/workflow/status/jangraefen/hcloud-pricing-exporter/build.yaml?branch=main&logo=GitHub)](https://github.com/jangraefen/hcloud-pricing-exporter/actions?query=workflow:Build)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/jangraefen/hcloud-pricing-exporter)](https://pkg.go.dev/mod/github.com/jangraefen/hcloud-pricing-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/jangraefen/hcloud-pricing-exporter)](https://goreportcard.com/report/github.com/jangraefen/hcloud-pricing-exporter)
[![Docker Pulls](https://img.shields.io/docker/pulls/jangraefen/hcloud-pricing-exporter)](https://hub.docker.com/r/jangraefen/hcloud-pricing-exporter)
[![Artifact HUB](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/hcloud-pricing-exporter)](https://artifacthub.io/packages/search?repo=hcloud-pricing-exporter)

A Prometheus exporter that connects to your HCloud account and collects data on your current expenses. The aim is to
make cost of cloud infrastructure more transparent and manageable, especially for private projects.

Please note that no guarantees on correctness are made and any financial decisions should be always be based on the
billing and cost functions provided by HCloud itself. Some hourly costs are estimations based on monthly costs, if the
HCloud API does not provide an hourly expense.

## Deployment

To run the exporter from the CLI you need to run the following commands:

```shell
# Just run it with the default settings
./hcloud-pricing-exporter -hcloud-token <TOKEN>

# Get the token from an ENV variable
export HCLOUD_TOKEN=<TOKEN>
./hcloud-pricing-exporter

# Run the exporter on a different port with another fetch interval
./hcloud-pricing-exporter -port 1234 -fetch-interval 45m
```

Alternatively, the exporter can be run by using the provided docker image:

```shell
docker run jangraefen/hcloud-pricing-exporter:latest -e HCLOUD_TOKEN=<TOKEN> -p 8080:8080
```

If you want to deploy the exporter to a Kubernetes environment, you can use the provided helm chart. Just perform the
following commands:

```shell
helm repo add hcloud-pricing-exporter https://jangraefen.github.io/hcloud-pricing-exporter
helm repo update
helm upgrade --install hcloud-pricing-exporter hcloud-pricing-exporter/hcloud-pricing-exporter --version {VERSION}
```

## Exported metrics

- `hcloud_pricing_floatingip_hourly{name, location, type}` _(Estimated based on the monthly price)_
- `hcloud_pricing_floatingip_monthly{name, location, type}`
- `hcloud_pricing_loadbalancer_hourly{name, location, type}`
- `hcloud_pricing_loadbalancer_monthly{name, location, type}`
- `hcloud_pricing_primaryip_hourly{name, datacenter, type}`
- `hcloud_pricing_primaryip_monthly{name, datacenter, type}`
- `hcloud_pricing_server_hourly{name, location, type}`
- `hcloud_pricing_server_monthly{name, location, type}`
- `hcloud_pricing_server_backups_hourly{name, location, type}`
- `hcloud_pricing_server_backups_monthly{name, location, type}`
- `hcloud_pricing_server_traffic_hourly{name, location, type}` _(Estimated based on the monthly price)_
- `hcloud_pricing_server_traffic_monthly{name, location, type}`
- `hcloud_pricing_snapshot_hourly{name}` _(Estimated based on the monthly price)_
- `hcloud_pricing_snapshot_monthly{name}`
- `hcloud_pricing_volume_hourly{name, location, bytes}` _(Estimated based on the monthly price)_
- `hcloud_pricing_volume_monthly{name, location, bytes}`
