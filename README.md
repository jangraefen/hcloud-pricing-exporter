# hcloud-pricing-exporter

[![Build Status](https://img.shields.io/github/workflow/status/jangraefen/hcloud-pricing-exporter/Build?logo=GitHub)](https://github.com/jangraefen/hcloud-pricing-exporter/actions?query=workflow:Build)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/jangraefen/hcloud-pricing-exporter)](https://pkg.go.dev/mod/github.com/jangraefen/hcloud-pricing-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/jangraefen/hcloud-pricing-exporter)](https://goreportcard.com/report/github.com/jangraefen/hcloud-pricing-exporter)
[![Docker Pulls](https://img.shields.io/docker/pulls/jangraefen/hcloud-pricing-exporter)](https://hub.docker.com/r/jangraefen/hcloud-pricing-exporter)

A Prometheus exporter that connects to your HCloud account and collects data on your current expenses. The aim is to
make cost of cloud infrastructure more transparent and manageable, especially for private projects.

Please note that no gurantees on correctness are made and any financial decisions should be always be based on the
billing and cost functions provided by HCloud itself. Some hourly costs are estimations based on monthly costs, if the
HCloud API does not provide an hourly expense.

## Exported metrics

- `hcloud_pricing_floatingip_hourly{name, location}` _(Estimated based on the monthly price)_
- `hcloud_pricing_floatingip_monthly{name, location}`
- `hcloud_pricing_loadbalancer_hourly{name, location, type}`
- `hcloud_pricing_loadbalancer_monthly{name, location, type}`
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
