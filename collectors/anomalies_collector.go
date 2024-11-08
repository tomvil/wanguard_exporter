package collectors

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	wgc "github.com/tomvil/wanguard_exporter/client"
)

type AnomaliesCollector struct {
	wgClient          *wgc.Client
	AnomalyActive     *prometheus.Desc
	AnomaliesFinished *prometheus.Desc
}

type AnomaliesCount struct {
	Count string
}

type Anomaly struct {
	AnomalyId string `json:"anomaly_id"`
	Prefix    string
	Anomaly   string
	Duration  string
	Pkts_s    string `json:"pkts/s"`
	Bits_s    string `json:"bits/s"`
	Packets   string
	Bits      string
}

func NewAnomaliesCollector(wgclient *wgc.Client) *AnomaliesCollector {
	prefix := "wanguard_anomalies_"
	return &AnomaliesCollector{
		wgClient:          wgclient,
		AnomalyActive:     prometheus.NewDesc(prefix+"active", "Active anomalies at the moment", []string{"prefix", "anomaly", "anomaly_id", "duration", "pkts_s", "packets", "bits_s", "bits"}, nil),
		AnomaliesFinished: prometheus.NewDesc(prefix+"finished", "Number of finished anomalies", nil, nil),
	}
}

func (c *AnomaliesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.AnomalyActive
	ch <- c.AnomaliesFinished
}

func (c *AnomaliesCollector) Collect(ch chan<- prometheus.Metric) {
	collectActiveAnomalies(c.AnomalyActive, c.wgClient, ch)
	collectFinishedAnomaliesTotal(c.AnomaliesFinished, c.wgClient, ch)
}

func collectActiveAnomalies(desc *prometheus.Desc, wgclient *wgc.Client, ch chan<- prometheus.Metric) {
	var anomalies []Anomaly

	err := wgclient.GetParsed("anomalies?status=Active&fields=anomaly_id,anomaly,prefix,duration,pkts/s,packets,bits/s,bits", &anomalies)
	if err != nil {
		return
	}

	for _, anomaly := range anomalies {
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, 1,
			anomaly.Prefix,
			anomaly.Anomaly,
			anomaly.AnomalyId,
			anomaly.Duration,
			anomaly.Pkts_s,
			anomaly.Packets,
			anomaly.Bits_s,
			anomaly.Bits)
	}
}

func collectFinishedAnomaliesTotal(desc *prometheus.Desc, wgclient *wgc.Client, ch chan<- prometheus.Metric) {
	var finishedAnomaliesCount AnomaliesCount

	err := wgclient.GetParsed("anomalies?status=Finished&count=true", &finishedAnomaliesCount)
	if err != nil {
		log.Errorln(err.Error())
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, 0)
		return
	}

	r, err := strconv.ParseFloat(finishedAnomaliesCount.Count, 64)
	if err != nil {
		log.Errorln(err.Error())
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, 0)
		return
	}

	ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, r)
}
