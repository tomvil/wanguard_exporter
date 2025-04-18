package collectors

import (
	"os"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	wgc "github.com/tomvil/wanguard_exporter/client"
)

func TestAnnouncementsCollector(t *testing.T) {
	wgcClient := wgc.NewClient(os.Getenv("TEST_SERVER_URL"), "u", "p")
	AnnouncementsCollector := NewAnnouncementsCollector(wgcClient)

	metricsCount := testutil.CollectAndCount(AnnouncementsCollector)
	if metricsCount != 2 {
		t.Errorf("Expected 2 metrics, got %d", metricsCount)
	}

	lintErrors, err := testutil.CollectAndLint(AnnouncementsCollector)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	for _, lintError := range lintErrors {
		t.Errorf("metric %v has lint error: %v", lintError.Metric, lintError.Text)
	}

	expectedMetrics := announcementsExpectedMetrics()
	err = testutil.CollectAndCompare(AnnouncementsCollector, strings.NewReader(expectedMetrics),
		"wanguard_announcements_active",
		"wanguard_announcements_finished")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func announcementsExpectedMetrics() string {
	return `
	# HELP wanguard_announcements_active Active announcements at the moment
	# TYPE wanguard_announcements_active gauge
	wanguard_announcements_active{announcement_id="1",bgp_connector_name="Connector 1",from="2024-10-23 09:31:01",prefix="10.10.10.10/32",until=""} 1
	# HELP wanguard_announcements_finished Total amount of finished announcements
	# TYPE wanguard_announcements_finished gauge
	wanguard_announcements_finished 1
`
}
