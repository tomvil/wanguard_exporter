package collectors

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/prometheus/common/log"
)

func TestMain(m *testing.M) {
	mux := http.NewServeMux()

	mux.HandleFunc("/wanguard-api/v1/license_manager", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(licenseManagerPayload())); err != nil {
			log.Errorln(err.Error())
		}
	})

	mux.HandleFunc("/wanguard-api/v1/firewall_rules", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("count") == "true" {
			if _, err := w.Write([]byte(`{"count": "1"}`)); err != nil {
				log.Errorln(err.Error())
			}
		}

		if r.URL.Query().Get("fields") == "firewall_rule_id,attack_id,source_prefix,destination_prefix,ip_protocol,from,until,pkts/s,bits/s,max_pkts/s,max_bits/s,pkts,bits" {
			if _, err := w.Write([]byte(firewallRulesPayload())); err != nil {
				log.Errorln(err.Error())
			}
		}
	})

	mux.HandleFunc("/wanguard-api/v1/bgp_connectors", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(bgpConnectorsPayload())); err != nil {
			log.Errorln(err.Error())
		}
	})

	mux.HandleFunc("/wanguard-api/v1/bgp_connectors/1/status", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(`{"status": "Active"}`)); err != nil {
			log.Errorln(err.Error())
		}
	})

	mux.HandleFunc("/wanguard-api/v1/filters", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(filtersPayload())); err != nil {
			log.Errorln(err.Error())
		}
	})

	mux.HandleFunc("/wanguard-api/v1/packet_filters/1/status", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(`{"status": "Active"}`)); err != nil {
			log.Errorln(err.Error())
		}
	})

	mux.HandleFunc("/wanguard-api/v1/sensors", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(sensorsPayload())); err != nil {
			log.Errorln(err.Error())
		}
	})

	mux.HandleFunc("/wanguard-api/v1/flow_sensors/1/status", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(`{"status": "Active"}`)); err != nil {
			log.Errorln(err.Error())
		}
	})

	mux.HandleFunc("/wanguard-api/v1/bgp_announcements", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("status") == "Finished" && r.URL.Query().Get("count") == "true" {
			if _, err := w.Write([]byte(`{"count": "1"}`)); err != nil {
				log.Errorln(err.Error())
			}
		}

		if r.URL.Query().Get("status") == "Active" && r.URL.Query().Get("fields") != "" {
			if _, err := w.Write([]byte(announcementsPayload())); err != nil {
				log.Errorln(err.Error())
			}
		}
	})

	mux.HandleFunc("/wanguard-api/v1/anomalies", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("status") == "Finished" && r.URL.Query().Get("count") == "true" {
			if _, err := w.Write([]byte(`{"count": "1"}`)); err != nil {
				log.Errorln(err.Error())
			}
		}

		if r.URL.Query().Get("status") == "Active" && r.URL.Query().Get("fields") != "" {
			if _, err := w.Write([]byte(anomaliesPayload())); err != nil {
				log.Errorln(err.Error())
			}
		}
	})

	mux.HandleFunc("/wanguard-api/v1/responses", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(responsesPayload())); err != nil {
			log.Errorln(err.Error())
		}
	})

	mux.HandleFunc("/wanguard-api/v1/responses/1/actions", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(actionsPayload())); err != nil {
			log.Errorln(err.Error())
		}
	})

	mux.HandleFunc("/wanguard-api/v1/responses/1/actions/1/status", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(`{"status": "Active"}`)); err != nil {
			log.Errorln(err.Error())
		}
	})

	mux.HandleFunc("/wanguard-api/v1/sensor_live_stats", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(sensorLiveStatsPayload())); err != nil {
			log.Errorln(err.Error())
		}
	})

	mux.HandleFunc("/wanguard-api/v1/sensor_live_tops", func(w http.ResponseWriter, r *http.Request) {
		topType := r.URL.Query().Get("top_type")
		unit := r.URL.Query().Get("unit")
		direction := r.URL.Query().Get("direction")

		if topType == "Countries" && (unit == "Packets" || unit == "Bits") && (direction == "Inbound" || direction == "Outbound") {
			if _, err := w.Write([]byte(countriesTopPayload())); err != nil {
				log.Errorln(err.Error())
			}
		}

		if topType == "IP Versions" && (unit == "Packets" || unit == "Bits") && (direction == "Inbound" || direction == "Outbound") {
			if _, err := w.Write([]byte(ipVersionsTopPayload())); err != nil {
				log.Errorln(err.Error())
			}
		}

		if topType == "IP Protocols" && (unit == "Packets" || unit == "Bits") && (direction == "Inbound" || direction == "Outbound") {
			if _, err := w.Write([]byte(ipProtocolsTopPayload())); err != nil {
				log.Errorln(err.Error())
			}
		}

		if topType == "Talkers" && (unit == "Packets" || unit == "Bits") && (direction == "Inbound" || direction == "Outbound") {
			if _, err := w.Write([]byte(talkersTopPayload())); err != nil {
				log.Errorln(err.Error())
			}
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
  "license_expiry_date_remaining": "1 days",
  "support_expiry_date_remaining": "1 days"
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

func bgpConnectorsPayload() string {
	return `[
  {
    "bgp_connector_id": "1",
    "bgp_connector_name": "BGP Connector 1",
    "href": "/wanguard-api/v1/bgp_connectors/1"
  }
]`
}

func filtersPayload() string {
	return `[
  {
    "packet_filter_id": "1",
    "filter_name": "Packet Filter 1",
    "href": "/wanguard-api/v1/packet_filters/1"
  }
]`
}

func sensorsPayload() string {
	return `[
  {
    "flow_sensor_id": "1",
    "sensor_name": "Flow Sensor 1",
    "href": "/wanguard-api/v1/flow_sensors/1"
  }
]`
}

func announcementsPayload() string {
	return `[
  {
    "bgp_announcement_id": "1",
    "prefix": "10.10.10.10/32",
    "from": {
      "iso_8601": "2024-10-23 09:31:01",
      "unixtime": "1729675861"
    },
    "until": {
      "iso_8601": "",
      "unixtime": ""
    },
    "href": "/wanguard-api/v1/bgp_announcements/1"
  }
]`
}

func anomaliesPayload() string {
	return `[
  {
    "anomaly_id": "1",
    "prefix": "10.10.10.10/32",
    "anomaly": "ICMP pkts/s > 1",
    "duration": "60",
    "pkts/s": "17500",
    "bits/s": "9014400",
    "packets": "320020500",
    "bits": "169576384000",
    "href": "/wanguard-api/v1/anomalies/1"
  }
]`
}

func responsesPayload() string {
	return `[
  {
    "response_id": "1",
    "response_name": "Response 1",
    "href": "/wanguard-api/v1/responses/1"
  }
]`
}

func actionsPayload() string {
	return `[
  {
    "action_id": "1",
    "status": {
      "href": "/wanguard-api/v1/responses/1/actions/1/status"
    },
    "action_name": "Action 1",
    "action_type": "Send a custom Syslog message",
    "response_branch": "When an anomaly is detected",
    "href": "/wanguard-api/v1/responses/1/actions/1"
  }
]`
}

func sensorLiveStatsPayload() string {
	return `[
  {
    "status": "Active",
    "sensor": {
      "sensor_interface_name": "Interface 1",
      "sensor_interface_id": "1",
      "sensor_interface_color": "#439F0A",
      "href": "/wanguard-api/v1/flow_sensors/1/interfaces/1"
    },
    "internal_ips": "1",
    "external_ips": "0",
    "packets/s_in": "100",
    "packets/s_out": "100",
    "bits/s_in": "1000",
    "bits/s_out": "1000",
    "dropped_in": "0",
    "dropped_out": "0",
    "usage_in": "1",
    "usage_out": "1",
    "load": "0.00",
    "cpu%": "0.00",
    "ram": 128,
    "start_time": {
      "iso_8601": "2024-10-28 00:00:01",
      "unixtime": "1730107187"
    },
    "server": {
      "server_id": "1",
      "server_name": "Server 1",
      "href": "/wanguard-api/v1/servers/1"
    }
  }
]`
}

func countriesTopPayload() string {
	return `{
  "top": {
    "1": {
      "country": "United States",
      "value": 200,
      "percent": 10
    }
  }
}`
}

func ipVersionsTopPayload() string {
	return `{
  "top": {
    "1": {
      "ip_version": 4,
      "description": "IPv4",
      "value": 100,
      "percent": 50
    }
  }
}`
}

func ipProtocolsTopPayload() string {
	return `{
  "top": {
    "1": {
      "ip_protocol": 6,
      "description": "Transmission Control",
      "value": 100,
      "percent": 50
    }
  }
}`
}

func talkersTopPayload() string {
	return `{
  "top": {
    "1": {
      "ip_address": "10.10.10.10",
      "ip_group_name": "External Zone",
      "value": 100,
      "percent": 10
    }
  }
}`
}
