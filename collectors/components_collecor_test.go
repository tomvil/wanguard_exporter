package collectors

import (
	"os"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	wgc "github.com/tomvil/wanguard_exporter/client"
)

func TestComponentsCollector(t *testing.T) {
	wgcClient := wgc.NewClient(os.Getenv("TEST_SERVER_URL"), "u", "p")
	ComponentsCollector := NewComponentsCollector(wgcClient)

	metricsCount := testutil.CollectAndCount(ComponentsCollector)
	if metricsCount != 3 {
		t.Errorf("Expected 3 metric, got %d", metricsCount)
	}

	lintErrors, err := testutil.CollectAndLint(ComponentsCollector)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	for _, lintError := range lintErrors {
		t.Errorf("metric %v has lint error: %v", lintError.Metric, lintError.Text)
	}

	expectedMetrics := componentsExpectedMetrics()
	err = testutil.CollectAndCompare(ComponentsCollector, strings.NewReader(expectedMetrics),
		"wanguard_component_status")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func componentsExpectedMetrics() string {
	return `
	# HELP wanguard_component_status Status of the component
	# TYPE wanguard_component_status gauge
	wanguard_component_status{component_category="bgp_connector",component_name="BGP Connector 1"} 1
	wanguard_component_status{component_category="filter",component_name="Packet Filter 1"} 1
	wanguard_component_status{component_category="sensor",component_name="Flow Sensor 1"} 1
`
}
