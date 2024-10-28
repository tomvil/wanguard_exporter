package collectors

import (
	"os"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	wgc "github.com/tomvil/wanguard_exporter/client"
)

func TestAnomaliesCollector(t *testing.T) {
	wgcClient := wgc.NewClient(os.Getenv("TEST_SERVER_URL"), "u", "p")
	AnomaliesCollector := NewAnomaliesCollector(wgcClient)

	metricsCount := testutil.CollectAndCount(AnomaliesCollector)
	if metricsCount != 2 {
		t.Errorf("Expected 2 metrics, got %d", metricsCount)
	}

	expectedMetrics := anomaliesExpectedMetrics()
	err := testutil.CollectAndCompare(AnomaliesCollector, strings.NewReader(expectedMetrics),
		"wanguard_anomalies_active",
		"wanguard_anomalies_total")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func anomaliesExpectedMetrics() string {
	return `
	# HELP wanguard_anomalies_active Active anomalies at the moment
	# TYPE wanguard_anomalies_active gauge
	wanguard_anomalies_active{anomaly="ICMP pkts/s > 1",bits="169576384000",bits_s="9014400",duration="60",packets="320020500",pkts_s="17500",prefix="10.10.10.10/32"} 1
	
	# HELP wanguard_anomalies_total Total amount of anomalies
	# TYPE wanguard_anomalies_total gauge
	wanguard_anomalies_total 1
`
}