package collectors

import (
	"os"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	wgc "github.com/tomvil/wanguard_exporter/client"
)

func TestSensorsCollector(t *testing.T) {
	wgcClient := wgc.NewClient(os.Getenv("TEST_SERVER_URL"), "u", "p")
	SensorsCollector := NewSensorsCollector(wgcClient)

	metricsCount := testutil.CollectAndCount(SensorsCollector)
	if metricsCount != 13 {
		t.Errorf("Expected 13 metrics, got %d", metricsCount)
	}

	expectedMetrics := sensorsExpectedMetrics()
	err := testutil.CollectAndCompare(SensorsCollector, strings.NewReader(expectedMetrics),
		"wanguard_sensor_internal_ips",
		"wanguard_sensor_external_ips",
		"wanguard_sensor_packets_per_second_in",
		"wanguard_sensor_packets_per_second_out",
		"wanguard_sensor_bits_per_second_in",
		"wanguard_sensor_bits_per_second_out",
		"wanguard_sensor_dropped_in",
		"wanguard_sensor_dropped_out",
		"wanguard_sensor_usage_in",
		"wanguard_sensor_usage_out",
		"wanguard_sensor_load",
		"wanguard_sensor_cpu",
		"wanguard_sensor_ram")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func sensorsExpectedMetrics() string {
	return `
	# HELP wanguard_sensor_bits_per_second_in Incoming bits per second
	# TYPE wanguard_sensor_bits_per_second_in gauge
	wanguard_sensor_bits_per_second_in{sensor_id="1",sensor_name="Interface 1"} 1000

	# HELP wanguard_sensor_bits_per_second_out Outgoing bits per second
	# TYPE wanguard_sensor_bits_per_second_out gauge
	wanguard_sensor_bits_per_second_out{sensor_id="1",sensor_name="Interface 1"} 1000

	# HELP wanguard_sensor_cpu Sensors CPU usage
	# TYPE wanguard_sensor_cpu gauge
	wanguard_sensor_cpu{sensor_id="1",sensor_name="Interface 1"} 0

	# HELP wanguard_sensor_dropped_in Total number of dropped packets in
	# TYPE wanguard_sensor_dropped_in gauge
	wanguard_sensor_dropped_in{sensor_id="1",sensor_name="Interface 1"} 0

	# HELP wanguard_sensor_dropped_out Total number of dropped packets out
	# TYPE wanguard_sensor_dropped_out gauge
	wanguard_sensor_dropped_out{sensor_id="1",sensor_name="Interface 1"} 0

	# HELP wanguard_sensor_external_ips Total number of external ip addresses
	# TYPE wanguard_sensor_external_ips gauge
	wanguard_sensor_external_ips{sensor_id="1",sensor_name="Interface 1"} 0

	# HELP wanguard_sensor_internal_ips Total number of internal ip addresses
	# TYPE wanguard_sensor_internal_ips gauge
	wanguard_sensor_internal_ips{sensor_id="1",sensor_name="Interface 1"} 1

	# HELP wanguard_sensor_load Sensors load
	# TYPE wanguard_sensor_load gauge
	wanguard_sensor_load{sensor_id="1",sensor_name="Interface 1"} 0

	# HELP wanguard_sensor_packets_per_second_in Incoming packets per second
	# TYPE wanguard_sensor_packets_per_second_in gauge
	wanguard_sensor_packets_per_second_in{sensor_id="1",sensor_name="Interface 1"} 100

	# HELP wanguard_sensor_packets_per_second_out Incoming packets per second
	# TYPE wanguard_sensor_packets_per_second_out gauge
	wanguard_sensor_packets_per_second_out{sensor_id="1",sensor_name="Interface 1"} 100

	# HELP wanguard_sensor_ram Sensors ram usage
	# TYPE wanguard_sensor_ram gauge
	wanguard_sensor_ram{sensor_id="1",sensor_name="Interface 1"} 128

	# HELP wanguard_sensor_usage_in Interface incoming traffic usage
	# TYPE wanguard_sensor_usage_in gauge
	wanguard_sensor_usage_in{sensor_id="1",sensor_name="Interface 1"} 1

	# HELP wanguard_sensor_usage_out Interface outgoing traffic usage
	# TYPE wanguard_sensor_usage_out gauge
	wanguard_sensor_usage_out{sensor_id="1",sensor_name="Interface 1"} 1
`
}
