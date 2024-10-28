package collectors

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	mux := http.NewServeMux()

	mux.HandleFunc("/wanguard-api/v1/license_manager", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(licenseManagerPayload()))
	})

	mux.HandleFunc("/wanguard-api/v1/firewall_rules", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("count") == "true" {
			w.Write([]byte(`{"count": "1"}`))
		}

		if r.URL.Query().Get("fields") == "firewall_rule_id,attack_id,source_prefix,destination_prefix,ip_protocol,from,until,pkts/s,bits/s,max_pkts/s,max_bits/s,pkts,bits" {
			w.Write([]byte(firewallRulesPayload()))
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	os.Setenv("TEST_SERVER_URL", server.URL)
	exitCode := m.Run()
	os.Unsetenv("TEST_SERVER_URL")

	os.Exit(exitCode)
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

func firewallRulesPayload() string {
	return `[
  {
    "firewall_rule_id": "1",
    "attack_id": "1",
    "source_prefix": "any",
    "destination_prefix": "any",
    "ip_protocol": "tcp",
    "from": {
      "iso_8601": "2024-10-28 06:37:02",
      "unixtime": "1730097422"
    },
    "until": {
      "iso_8601": "",
      "unixtime": ""
    },
    "pkts/s": "0",
    "bits/s": "0",
    "max_pkts/s": "0",
    "max_bits/s": "0",
    "pkts": "0",
    "bits": "0",
    "href": "/wanguard-api/v1/firewall_rules/1"
  }
]`
}
