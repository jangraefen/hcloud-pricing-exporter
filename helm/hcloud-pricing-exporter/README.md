# hcloud-pricing-exporter Helm Chart

[![Build Status](https://img.shields.io/github/workflow/status/jangraefen/hcloud-pricing-exporter/Build?logo=GitHub)](https://github.com/jangraefen/hcloud-pricing-exporter/actions?query=workflow:Build)
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
helm repo add hcloud-pricing-exporter https://jangraefen.github.io/hcloud-pricing-exporter
helm repo update
helm upgrade --install hcloud-pricing-exporter hcloud-pricing-exporter/hcloud-pricing-exporter
```

## Configuration

Parameter                      | Default                                | Description
------------------------------ | -------------------------------------- | -----------
`replicaCount`                 | `1`                                    |
`image.repository`             | `"jangraefen/hcloud-pricing-exporter"` |
`image.pullPolicy`             | `"IfNotPresent"`                       |
`image.tag`                    | `""`                                   |
`imagePullSecrets`             | `[]`                                   |
`nameOverride`                 | `""`                                   |
`fullnameOverride`             | `""`                                   |
`podAnnotations`               | `{}`                                   |
`service.type`                 | `"ClusterIP"`                          |
`service.port`                 | `8080`                                 |
`ingress.enabled`              | `false`                                |
`ingress.annotations`          | `{}`                                   |
`ingress.hosts`                | `[{"host": "chart-example.local"}]`    |
`ingress.tls`                  | `[]`                                   |
`secret.token`                 | `null`                                 | The API token to access your HCloud data.
`secret.create`                | `true`                                 | If you want to provision the secret for the API token yourself, set this to `false`.
`secret.reference.name`        | `null`                                 | The name of the secret that contains the API token to access your HCloud data.
`secret.reference.key`         | `null`                                 | The key of the secret that contains the API token to access your HCloud data.
`serviceMonitor.create`        | `false`                                | Enable this if you want to monitor the exporter with the Prometheus Operator.
`serviceMonitor.interval`      | `null`                                 |
`serviceMonitor.labels`        | `null`                                 |
`serviceMonitor.scrapeTimeout` | `null`                                 |
`resources`                    | `{}`                                   |
`nodeSelector`                 | `{}`                                   |
`tolerations`                  | `[]`                                   |
`affinity`                     | `{}`                                   |
