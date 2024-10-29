# wanguard_exporter
This is prometheus exporter for Andrisoft WANGuard API.

## Install
```go install github.com/tomvil/wanguard_exporter@latest```

## Configuration flags
Name     | Description | Default
---------|-------------|---------
version | Print information about exporter version |
web.listen-address | Address on which to expose metrics | :9868
web.metrics-path | Path under which to expose metrics | /metrics
api.address | WANGuard API Address | 127.0.0.1:81
api.username | WANGuard API Username | admin
api.password | WANGuard API Password |
licenseCollectorEnabled | Export license metrics | true
announcementsCollectorEnabled | Export announcements metrics | true
anomaliesCollectorEnabled | Export anomalies metrics | true
componentsCollectorEnabled | Export components metrics | true
actionsCollectorEnabled | Export actions metrics | true
sensorsCollectorEnabled | Export sensors metrics | true
trafficCollectorEnabled | Export traffic metrics | true
firewallRulesCollectorEnabled | Export firewall rules metrics | true

## Configuration environment variables
Name     | Description
---------|-------------
WANGUARD_PASSWORD | WANGuard API Password

This will be used automatically if `api.password` flag is not set.


## Usage
`./wanguard_exporter -api.address="127.0.0.1:81" -api.username="admin" api.password="password"`

## Metrics
The latest, accurate metrics for each collector can be viewed by accessing the corresponding test file located at `collectors/name_collector_test.go`. In this file, the `nameExpectedMetrics()` function provides users with a structured way to examine the expected metrics associated with each collector.
- Change `name` with actual collectors name.
