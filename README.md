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
```
# HELP wanguard_action_status Status of the response actions
# TYPE wanguard_action_status gauge
wanguard_action_status{action_name="test",action_type="Execute a command or script with dynamic parameters as arguments",response_branch="test"} 1
# HELP wanguard_announcements_total Total amount of announcements
# TYPE wanguard_announcements_total gauge
wanguard_announcements_total 5
# HELP wanguard_anomalies_total Total amount of anomalies
# TYPE wanguard_anomalies_total gauge
wanguard_anomalies_total 5
# HELP wanguard_component_status Status of the component
# TYPE wanguard_component_status gauge
wanguard_component_status{component_category="sensor",component_name="test"} 1
# HELP wanguard_license_days_remaining License days remaining
# TYPE wanguard_license_days_remaining gauge
wanguard_license_days_remaining 9999
# HELP wanguard_license_filters Licensed filters total
# TYPE wanguard_license_filters gauge
wanguard_license_filters 0
# HELP wanguard_license_filters_remaining Licensed filters available
# TYPE wanguard_license_filters_remaining gauge
wanguard_license_filters_remaining 0
# HELP wanguard_license_filters_used Licensed filters total
# TYPE wanguard_license_filters_used gauge
wanguard_license_filters_used 0
# HELP wanguard_license_sensors_remaining Licensed sensors remaining
# TYPE wanguard_license_sensors_remaining gauge
wanguard_license_sensors_remaining 0
# HELP wanguard_license_sensors_total Licensed sensors total
# TYPE wanguard_license_sensors_total gauge
wanguard_license_sensors_total 1
# HELP wanguard_license_sensors_used Licensed sensors used
# TYPE wanguard_license_sensors_used gauge
wanguard_license_sensors_used 1
# HELP wanguard_license_support_days_remaining Support license days remaining
# TYPE wanguard_license_support_days_remaining gauge
wanguard_license_support_days_remaining 9999
# HELP wanguard_sensor_bits_per_second_in Incoming bits per second
# TYPE wanguard_sensor_bits_per_second_in gauge
wanguard_sensor_bits_per_second_in{sensor_id="1",sensor_name="test"} 3.85019084e+08
# HELP wanguard_sensor_bits_per_second_out Outgoing bits per second
# TYPE wanguard_sensor_bits_per_second_out gauge
wanguard_sensor_bits_per_second_out{sensor_id="1",sensor_name="test"} 4.929323e+08
# HELP wanguard_sensor_cpu Sensors CPU usage
# TYPE wanguard_sensor_cpu gauge
wanguard_sensor_cpu{sensor_id="1",sensor_name="test"} 0.2
# HELP wanguard_sensor_dropped_in Total number of dropped packets in
# TYPE wanguard_sensor_dropped_in gauge
wanguard_sensor_dropped_in{sensor_id="1",sensor_name="test"} 0
# HELP wanguard_sensor_dropped_out Total number of dropped packets out
# TYPE wanguard_sensor_dropped_out gauge
wanguard_sensor_dropped_out{sensor_id="1",sensor_name="test"} 0
# HELP wanguard_sensor_external_ips Total number of external ip addresses
# TYPE wanguard_sensor_external_ips gauge
wanguard_sensor_external_ips{sensor_id="1",sensor_name="test"} 0
# HELP wanguard_sensor_internal_ips Total number of internal ip addresses
# TYPE wanguard_sensor_internal_ips gauge
wanguard_sensor_internal_ips{sensor_id="1",sensor_name="test"} 68
# HELP wanguard_sensor_load Sensors load
# TYPE wanguard_sensor_load gauge
wanguard_sensor_load{sensor_id="1",sensor_name="test"} 0.06
# HELP wanguard_sensor_packets_per_second_in Incoming packets per second
# TYPE wanguard_sensor_packets_per_second_in gauge
wanguard_sensor_packets_per_second_in{sensor_id="1",sensor_name="test"} 56934
# HELP wanguard_sensor_packets_per_second_out Incoming packets per second
# TYPE wanguard_sensor_packets_per_second_out gauge
wanguard_sensor_packets_per_second_out{sensor_id="1",sensor_name="test"} 58163
# HELP wanguard_sensor_ram Sensors ram usage
# TYPE wanguard_sensor_ram gauge
wanguard_sensor_ram{sensor_id="1",sensor_name="test"} 215
# HELP wanguard_sensor_usage_in Interface incoming traffic usage
# TYPE wanguard_sensor_usage_in gauge
wanguard_sensor_usage_in{sensor_id="1",sensor_name="test"} 4
# HELP wanguard_sensor_usage_out Interface outgoing traffic usage
# TYPE wanguard_sensor_usage_out gauge
wanguard_sensor_usage_out{sensor_id="1",sensor_name="test"} 5
# HELP wanguard_traffic_country_bits_per_second_in Bits per second in by country
# TYPE wanguard_traffic_country_bits_per_second_in gauge
wanguard_traffic_country_bits_per_second_in{country="Azerbaijan",country_code="AZ"} 5.39525e+06
wanguard_traffic_country_bits_per_second_in{country="Brazil",country_code="BR"} 907673
wanguard_traffic_country_bits_per_second_in{country="Canada",country_code="CA"} 1.161625e+06
wanguard_traffic_country_bits_per_second_in{country="Cyprus",country_code="CY"} 1.600192511e+09
wanguard_traffic_country_bits_per_second_in{country="Estonia",country_code="EE"} 2.493644e+06
# HELP wanguard_traffic_country_bits_per_second_out Bits per second out by country
# TYPE wanguard_traffic_country_bits_per_second_out gauge
wanguard_traffic_country_bits_per_second_out{country="Armenia",country_code="AM"} 4.987289e+06
wanguard_traffic_country_bits_per_second_out{country="Austria",country_code="AT"} 4.987289e+06
wanguard_traffic_country_bits_per_second_out{country="Bangladesh",country_code="BD"} 9.777971e+06
wanguard_traffic_country_bits_per_second_out{country="Belarus",country_code="BY"} 4.987289e+06
wanguard_traffic_country_bits_per_second_out{country="Canada",country_code="CA"} 4.987289e+06
# HELP wanguard_traffic_country_packets_per_second_in Packets per second in by country
# TYPE wanguard_traffic_country_packets_per_second_in gauge
wanguard_traffic_country_packets_per_second_in{country="Azerbaijan",country_code="AZ"} 1638
wanguard_traffic_country_packets_per_second_in{country="Canada",country_code="CA"} 818
wanguard_traffic_country_packets_per_second_in{country="Colombia",country_code="CO"} 614
wanguard_traffic_country_packets_per_second_in{country="Cyprus",country_code="CY"} 175513
wanguard_traffic_country_packets_per_second_in{country="France",country_code="FR"} 2048
# HELP wanguard_traffic_country_packets_per_second_out Packets per second out by country
# TYPE wanguard_traffic_country_packets_per_second_out gauge
wanguard_traffic_country_packets_per_second_out{country="Armenia",country_code="AM"} 409
wanguard_traffic_country_packets_per_second_out{country="Bangladesh",country_code="BD"} 819
wanguard_traffic_country_packets_per_second_out{country="Brazil",country_code="BR"} 409
wanguard_traffic_country_packets_per_second_out{country="Cyprus",country_code="CY"} 166707
wanguard_traffic_country_packets_per_second_out{country="Finland",country_code="FI"} 409
# HELP wanguard_traffic_ip_protocol_bits_per_second_in Bits per second in by IP protocol
# TYPE wanguard_traffic_ip_protocol_bits_per_second_in gauge
wanguard_traffic_ip_protocol_bits_per_second_in{ip_protocol="ICMP"} 484966
wanguard_traffic_ip_protocol_bits_per_second_in{ip_protocol="TCP"} 4.42463026e+08
wanguard_traffic_ip_protocol_bits_per_second_in{ip_protocol="UDP"} 1.61788887e+09
# HELP wanguard_traffic_ip_protocol_bits_per_second_out Bits per second out by IP protocol
# TYPE wanguard_traffic_ip_protocol_bits_per_second_out gauge
wanguard_traffic_ip_protocol_bits_per_second_out{ip_protocol="TCP"} 4.33071717e+08
wanguard_traffic_ip_protocol_bits_per_second_out{ip_protocol="UDP"} 1.646831206e+09
# HELP wanguard_traffic_ip_protocol_packets_per_second_in Packets per second in by IP protocol
# TYPE wanguard_traffic_ip_protocol_packets_per_second_in gauge
wanguard_traffic_ip_protocol_packets_per_second_in{ip_protocol="ICMP"} 818
wanguard_traffic_ip_protocol_packets_per_second_in{ip_protocol="TCP"} 94208
wanguard_traffic_ip_protocol_packets_per_second_in{ip_protocol="UDP"} 178380
# HELP wanguard_traffic_ip_protocol_packets_per_second_out Packets per second out by IP protocol
# TYPE wanguard_traffic_ip_protocol_packets_per_second_out gauge
wanguard_traffic_ip_protocol_packets_per_second_out{ip_protocol="TCP"} 41778
wanguard_traffic_ip_protocol_packets_per_second_out{ip_protocol="UDP"} 172441
# HELP wanguard_traffic_ip_version_bits_per_second_in Bits per second in by IP version
# TYPE wanguard_traffic_ip_version_bits_per_second_in gauge
wanguard_traffic_ip_version_bits_per_second_in{ip_version="IPv4"} 1.969233919e+09
wanguard_traffic_ip_version_bits_per_second_in{ip_version="IPv6"} 9.1602943e+07
# HELP wanguard_traffic_ip_version_bits_per_second_out Bits per second out by IP version
# TYPE wanguard_traffic_ip_version_bits_per_second_out gauge
wanguard_traffic_ip_version_bits_per_second_out{ip_version="IPv4"} 2.029151846e+09
wanguard_traffic_ip_version_bits_per_second_out{ip_version="IPv6"} 5.0751078e+07
# HELP wanguard_traffic_ip_version_packets_per_second_in Packets per second in by IP version
# TYPE wanguard_traffic_ip_version_packets_per_second_in gauge
wanguard_traffic_ip_version_packets_per_second_in{ip_version="IPv4"} 250265
wanguard_traffic_ip_version_packets_per_second_in{ip_version="IPv6"} 23141
# HELP wanguard_traffic_ip_version_packets_per_second_out Packets per second out by IP version
# TYPE wanguard_traffic_ip_version_packets_per_second_out gauge
wanguard_traffic_ip_version_packets_per_second_out{ip_version="IPv4"} 208895
wanguard_traffic_ip_version_packets_per_second_out{ip_version="IPv6"} 5324
# HELP wanguard_traffic_talkers_bits_per_second_in Bits per second in by IP address
# TYPE wanguard_traffic_talkers_bits_per_second_in gauge
wanguard_traffic_talkers_bits_per_second_in{ip_address="192.168.1.1"} 4.987289e+06
wanguard_traffic_talkers_bits_per_second_in{ip_address="192.168.1.2"} 1.119820185e+09
# HELP wanguard_traffic_talkers_bits_per_second_out Bits per second out by IP address
# TYPE wanguard_traffic_talkers_bits_per_second_out gauge
wanguard_traffic_talkers_bits_per_second_out{ip_address="192.168.1.1"} 4.96753048e+08
wanguard_traffic_talkers_bits_per_second_out{ip_address="192.168.1.2"} 1.4574551e+08
# HELP wanguard_traffic_talkers_packets_per_second_in Packets per second in by IP address
# TYPE wanguard_traffic_talkers_packets_per_second_in gauge
wanguard_traffic_talkers_packets_per_second_in{ip_address="192.168.1.1"} 2662
wanguard_traffic_talkers_packets_per_second_in{ip_address="192.168.1.2"} 123289
# HELP wanguard_traffic_talkers_packets_per_second_out Packets per second out by IP address
# TYPE wanguard_traffic_talkers_packets_per_second_out gauge
wanguard_traffic_talkers_packets_per_second_out{ip_address="192.168.1.1"} 42188
wanguard_traffic_talkers_packets_per_second_out{ip_address="192.168.1.2"} 22937
```
