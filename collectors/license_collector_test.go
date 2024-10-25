package collectors

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	wgc "github.com/tomvil/wanguard_exporter/client"
)

func TestMain(m *testing.M) {
	mux := http.NewServeMux()
	mux.HandleFunc("/wanguard-api/v1/license_manager", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(licenseManagerPayload()))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	os.Setenv("TEST_SERVER_URL", server.URL)
	exitCode := m.Run()
	os.Unsetenv("TEST_SERVER_URL")

	os.Exit(exitCode)
}

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
		"wanguard_license_days_remaining",
		"wanguard_license_support_days_remaining")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func licenseManagerPayload() string {
	return `{
  "software_version": "8.3-21",
  "licensed_sensors": 1,
  "licensed_sensors_used": 1,
  "licensed_sensors_remaining": 0,
  "licensed_sensor_interfaces": "∞",
  "licensed_dpdk_engines": 0,
  "licensed_dpdk_engines_used": 0,
  "licensed_dpdk_engines_remaining": 0,
  "licensed_filters": 1,
  "licensed_filters_used": 1,
  "licensed_filters_remaining": 0,
  "licensed_on": "2024-10-24 12:15:29",
  "license_expiry_date_remaining": "365 days",
  "support_expiry_date_remaining": "365 days"
}`
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

	# HELP wanguard_license_days_remaining License days remaining
	# TYPE wanguard_license_days_remaining gauge
	wanguard_license_days_remaining 365

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

	# HELP wanguard_license_support_days_remaining Support license days remaining
	# TYPE wanguard_license_support_days_remaining gauge
	wanguard_license_support_days_remaining 365
`
}
