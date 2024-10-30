package collectors

import (
	"os"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	wgc "github.com/tomvil/wanguard_exporter/client"
)

func TestFirewallRulesCollector(t *testing.T) {
	wgcClient := wgc.NewClient(os.Getenv("TEST_SERVER_URL"), "u", "p")
	FirewallRulesCollector := NewFirewallRulesCollector(wgcClient)

	metricsCount := testutil.CollectAndCount(FirewallRulesCollector)
	if metricsCount != 2 {
		t.Errorf("Expected 2 metrics, got %d", metricsCount)
	}

	lintErrors, err := testutil.CollectAndLint(FirewallRulesCollector)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	for _, lintError := range lintErrors {
		t.Errorf("metric %v has lint error: %v", lintError.Metric, lintError.Text)
	}

	expectedMetrics := firewallRulesExpectedMetrics()
	err = testutil.CollectAndCompare(FirewallRulesCollector, strings.NewReader(expectedMetrics),
		"wanguard_firewall_rules_active",
		"wanguard_firewall_rules_activated")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func firewallRulesExpectedMetrics() string {
	return `
	# HELP wanguard_firewall_rules_active Active firewall rules at the moment
	# TYPE wanguard_firewall_rules_active gauge
	wanguard_firewall_rules_active{attack_id="1",bits="0",bits_s="0",destination_prefix="any",from="2024-10-28 06:37:02",ip_protocol="tcp",max_bits_s="0",max_pkts_s="0",pkts="0",pkts_s="0",rule_id="1",source_prefix="any",until=""} 1

	# HELP wanguard_firewall_rules_activated Number of activated firewall rules
	# TYPE wanguard_firewall_rules_activated gauge
	wanguard_firewall_rules_activated 1
`
}
