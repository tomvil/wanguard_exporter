package collectors

import (
	"os"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	wgc "github.com/tomvil/wanguard_exporter/client"
)

func TestActionsCollector(t *testing.T) {
	wgcClient := wgc.NewClient(os.Getenv("TEST_SERVER_URL"), "u", "p")
	ActionsCollector := NewActionsCollector(wgcClient)

	metricsCount := testutil.CollectAndCount(ActionsCollector)
	if metricsCount != 1 {
		t.Errorf("Expected 1 metric, got %d", metricsCount)
	}

	expectedMetrics := actionsExpectedMetrics()
	err := testutil.CollectAndCompare(ActionsCollector, strings.NewReader(expectedMetrics),
		"wanguard_action_status")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func actionsExpectedMetrics() string {
	return `
	# HELP wanguard_action_status Status of the response actions
	# TYPE wanguard_action_status gauge
	wanguard_action_status{action_name="Action 1",action_type="Send a custom Syslog message",response_branch="When an anomaly is detected",response_name="Response 1"} 1
`
}
