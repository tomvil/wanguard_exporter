package collectors

import (
	"os"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	wgc "github.com/tomvil/wanguard_exporter/client"
)

func TestLicenseCollector(t *testing.T) {
	wgcClient := wgc.NewClient(os.Getenv("TEST_SERVER_URL"), "u", "p")
	licenseCollector := NewLicenseCollector(wgcClient)

	metricsCount := testutil.CollectAndCount(licenseCollector)
	if metricsCount != 12 {
		t.Errorf("Expected 12 metrics, got %d", metricsCount)
	}

	expectedMetrics := licenseManagerExpectedMetrics()
	err := testutil.CollectAndCompare(licenseCollector, strings.NewReader(expectedMetrics),
		"wanguard_license_software_version",
		"wanguard_license_sensors_total",
		"wanguard_license_sensors_used",
		"wanguard_license_sensors_remaining",
		"wanguard_license_dpdk_engines_total",
		"wanguard_license_dpdk_engines_used",
		"wanguard_license_dpdk_engines_remaining",
		"wanguard_license_filters",
		"wanguard_license_filters_used",
		"wanguard_license_filters_remaining",
		"wanguard_license_seconds_remaining",
		"wanguard_license_support_seconds_remaining")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func licenseManagerExpectedMetrics() string {
	return `
	# HELP wanguard_license_dpdk_engines_total Licensed DPDK engines total
	# TYPE wanguard_license_dpdk_engines_total gauge
	wanguard_license_dpdk_engines_total 0

	# HELP wanguard_license_dpdk_engines_remaining Licensed DPDK engines remaining
	# TYPE wanguard_license_dpdk_engines_remaining gauge
	wanguard_license_dpdk_engines_remaining 0

	# HELP wanguard_license_dpdk_engines_used Licensed DPDK engines used
	# TYPE wanguard_license_dpdk_engines_used gauge
	wanguard_license_dpdk_engines_used 0

	# HELP wanguard_license_filters Licensed filters total
	# TYPE wanguard_license_filters gauge
	wanguard_license_filters 1

	# HELP wanguard_license_filters_remaining Licensed filters available
	# TYPE wanguard_license_filters_remaining gauge
	wanguard_license_filters_remaining 0

	# HELP wanguard_license_filters_used Licensed filters total
	# TYPE wanguard_license_filters_used gauge
	wanguard_license_filters_used 1

	# HELP wanguard_license_seconds_remaining License seconds remaining
	# TYPE wanguard_license_seconds_remaining gauge
	wanguard_license_seconds_remaining 86400

	# HELP wanguard_license_sensors_total Licensed sensors total
	# TYPE wanguard_license_sensors_total gauge
	wanguard_license_sensors_total 1

	# HELP wanguard_license_sensors_remaining Licensed sensors remaining
	# TYPE wanguard_license_sensors_remaining gauge
	wanguard_license_sensors_remaining 0

	# HELP wanguard_license_sensors_used Licensed sensors used
	# TYPE wanguard_license_sensors_used gauge
	wanguard_license_sensors_used 1

	# HELP wanguard_license_software_version Software version
	# TYPE wanguard_license_software_version gauge
	wanguard_license_software_version{software_version="8.3-21"} 1

	# HELP wanguard_license_support_seconds_remaining Support license seconds remaining
	# TYPE wanguard_license_support_seconds_remaining gauge
	wanguard_license_support_seconds_remaining 86400
`
}
