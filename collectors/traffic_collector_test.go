package collectors

import (
	"os"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	wgc "github.com/tomvil/wanguard_exporter/client"
)

func TestTrafficCollector(t *testing.T) {
	wgcClient := wgc.NewClient(os.Getenv("TEST_SERVER_URL"), "u", "p")
	TrafficCollector := NewTrafficCollector(wgcClient)

	metricsCount := testutil.CollectAndCount(TrafficCollector)
	if metricsCount != 16 {
		t.Errorf("Expected 16 metrics, got %d", metricsCount)
	}

	lintErrors, err := testutil.CollectAndLint(TrafficCollector)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	for _, lintError := range lintErrors {
		t.Errorf("metric %v has lint error: %v", lintError.Metric, lintError.Text)
	}

	expectedMetrics := trafficExpectedMetrics()
	err = testutil.CollectAndCompare(TrafficCollector, strings.NewReader(expectedMetrics),
		"wanguard_traffic_country_packets_per_second_in",
		"wanguard_traffic_country_packets_per_second_out",
		"wanguard_traffic_country_bytes_per_second_in",
		"wanguard_traffic_country_bytes_per_second_out",
		"wanguard_traffic_ip_version_packets_per_second_in",
		"wanguard_traffic_ip_version_packets_per_second_out",
		"wanguard_traffic_ip_version_bytes_per_second_in",
		"wanguard_traffic_ip_version_bytes_per_second_out",
		"wanguard_traffic_ip_protocol_packets_per_second_in",
		"wanguard_traffic_ip_protocol_packets_per_second_out",
		"wanguard_traffic_ip_protocol_bytes_per_second_in",
		"wanguard_traffic_ip_protocol_bytes_per_second_out",
		"wanguard_traffic_talkers_packets_per_second_in",
		"wanguard_traffic_talkers_packets_per_second_out",
		"wanguard_traffic_talkers_bytes_per_second_in",
		"wanguard_traffic_talkers_bytes_per_second_out")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func trafficExpectedMetrics() string {
	return `
	# HELP wanguard_traffic_country_bytes_per_second_in bytes per second in by country
	# TYPE wanguard_traffic_country_bytes_per_second_in gauge
	wanguard_traffic_country_bytes_per_second_in{country="United States",country_code="US"} 200

	# HELP wanguard_traffic_country_bytes_per_second_out bytes per second out by country
	# TYPE wanguard_traffic_country_bytes_per_second_out gauge
	wanguard_traffic_country_bytes_per_second_out{country="United States",country_code="US"} 200

	# HELP wanguard_traffic_country_packets_per_second_in Packets per second in by country
	# TYPE wanguard_traffic_country_packets_per_second_in gauge
	wanguard_traffic_country_packets_per_second_in{country="United States",country_code="US"} 200

	# HELP wanguard_traffic_country_packets_per_second_out Packets per second out by country
	# TYPE wanguard_traffic_country_packets_per_second_out gauge
	wanguard_traffic_country_packets_per_second_out{country="United States",country_code="US"} 200

	# HELP wanguard_traffic_ip_protocol_bytes_per_second_in bytes per second in by IP protocol
	# TYPE wanguard_traffic_ip_protocol_bytes_per_second_in gauge
	wanguard_traffic_ip_protocol_bytes_per_second_in{ip_protocol="TCP"} 100

	# HELP wanguard_traffic_ip_protocol_bytes_per_second_out bytes per second out by IP protocol
	# TYPE wanguard_traffic_ip_protocol_bytes_per_second_out gauge
	wanguard_traffic_ip_protocol_bytes_per_second_out{ip_protocol="TCP"} 100

	# HELP wanguard_traffic_ip_protocol_packets_per_second_in Packets per second in by IP protocol
	# TYPE wanguard_traffic_ip_protocol_packets_per_second_in gauge
	wanguard_traffic_ip_protocol_packets_per_second_in{ip_protocol="TCP"} 100

	# HELP wanguard_traffic_ip_protocol_packets_per_second_out Packets per second out by IP protocol
	# TYPE wanguard_traffic_ip_protocol_packets_per_second_out gauge
	wanguard_traffic_ip_protocol_packets_per_second_out{ip_protocol="TCP"} 100

	# HELP wanguard_traffic_ip_version_bytes_per_second_in bytes per second in by IP version
	# TYPE wanguard_traffic_ip_version_bytes_per_second_in gauge
	wanguard_traffic_ip_version_bytes_per_second_in{ip_version="IPv4"} 100

	# HELP wanguard_traffic_ip_version_bytes_per_second_out bytes per second out by IP version
	# TYPE wanguard_traffic_ip_version_bytes_per_second_out gauge
	wanguard_traffic_ip_version_bytes_per_second_out{ip_version="IPv4"} 100

	# HELP wanguard_traffic_ip_version_packets_per_second_in Packets per second in by IP version
	# TYPE wanguard_traffic_ip_version_packets_per_second_in gauge
	wanguard_traffic_ip_version_packets_per_second_in{ip_version="IPv4"} 100

	# HELP wanguard_traffic_ip_version_packets_per_second_out Packets per second out by IP version
	# TYPE wanguard_traffic_ip_version_packets_per_second_out gauge
	wanguard_traffic_ip_version_packets_per_second_out{ip_version="IPv4"} 100

	# HELP wanguard_traffic_talkers_bytes_per_second_in bytes per second in by IP address
	# TYPE wanguard_traffic_talkers_bytes_per_second_in gauge
	wanguard_traffic_talkers_bytes_per_second_in{ip_address="10.10.10.10"} 100

	# HELP wanguard_traffic_talkers_bytes_per_second_out bytes per second out by IP address
	# TYPE wanguard_traffic_talkers_bytes_per_second_out gauge
	wanguard_traffic_talkers_bytes_per_second_out{ip_address="10.10.10.10"} 100

	# HELP wanguard_traffic_talkers_packets_per_second_in Packets per second in by IP address
	# TYPE wanguard_traffic_talkers_packets_per_second_in gauge
	wanguard_traffic_talkers_packets_per_second_in{ip_address="10.10.10.10"} 100

	# HELP wanguard_traffic_talkers_packets_per_second_out Packets per second out by IP address
	# TYPE wanguard_traffic_talkers_packets_per_second_out gauge
	wanguard_traffic_talkers_packets_per_second_out{ip_address="10.10.10.10"} 100
`
}
