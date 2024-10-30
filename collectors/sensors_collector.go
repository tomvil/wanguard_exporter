package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	wgc "github.com/tomvil/wanguard_exporter/client"
)

type SensorsCollector struct {
	wgClient          *wgc.Client
	SensorInternalIPS *prometheus.Desc
	SensorExternalIPS *prometheus.Desc
	SensorPPSIn       *prometheus.Desc
	SensorPPSOut      *prometheus.Desc
	SensorBPSIn       *prometheus.Desc
	SensorBPSOut      *prometheus.Desc
	SensorDroppedIn   *prometheus.Desc
	SensorDroppedOut  *prometheus.Desc
	SensorUsageIn     *prometheus.Desc
	SensorUsageOut    *prometheus.Desc
	SensorLoad        *prometheus.Desc
	SensorCpu         *prometheus.Desc
	SensorRam         *prometheus.Desc
}

type Sensor struct {
	Sensor struct {
		InterfaceName string `json:"sensor_interface_name"`
		InterfaceID   string `json:"sensor_interface_id"`
	}
	InternalIPS         string `json:"internal_ips"`
	ExternalIPS         string `json:"external_ips"`
	PacketsPerSecondIN  string `json:"packets/s_in"`
	PacketsPerSecondOUT string `json:"packets/s_out"`
	BitsPerSecondIN     string `json:"bits/s_in"`
	BitsPerSecondOUT    string `json:"bits/s_out"`
	DroppedIN           string `json:"dropped_in"`
	DroppedOUT          string `json:"dropped_out"`
	UsageIN             string `json:"usage_in"`
	UsageOUT            string `json:"usage_out"`
	Load                string
	Cpu                 string `json:"cpu%"`
	Ram                 int
}

func NewSensorsCollector(wgclient *wgc.Client) *SensorsCollector {
	prefix := "wanguard_sensor_"
	return &SensorsCollector{
		wgClient:          wgclient,
		SensorInternalIPS: prometheus.NewDesc(prefix+"internal_ips", "Total number of internal ip addresses", []string{"sensor_name", "sensor_id"}, nil),
		SensorExternalIPS: prometheus.NewDesc(prefix+"external_ips", "Total number of external ip addresses", []string{"sensor_name", "sensor_id"}, nil),
		SensorPPSIn:       prometheus.NewDesc(prefix+"packets_per_second_in", "Incoming packets per second", []string{"sensor_name", "sensor_id"}, nil),
		SensorPPSOut:      prometheus.NewDesc(prefix+"packets_per_second_out", "Incoming packets per second", []string{"sensor_name", "sensor_id"}, nil),
		SensorBPSIn:       prometheus.NewDesc(prefix+"bytes_per_second_in", "Incoming bytes per second", []string{"sensor_name", "sensor_id"}, nil),
		SensorBPSOut:      prometheus.NewDesc(prefix+"bytes_per_second_out", "Outgoing bytes per second", []string{"sensor_name", "sensor_id"}, nil),
		SensorDroppedIn:   prometheus.NewDesc(prefix+"dropped_in", "Total number of dropped packets in", []string{"sensor_name", "sensor_id"}, nil),
		SensorDroppedOut:  prometheus.NewDesc(prefix+"dropped_out", "Total number of dropped packets out", []string{"sensor_name", "sensor_id"}, nil),
		SensorUsageIn:     prometheus.NewDesc(prefix+"usage_in", "Interface incoming traffic usage", []string{"sensor_name", "sensor_id"}, nil),
		SensorUsageOut:    prometheus.NewDesc(prefix+"usage_out", "Interface outgoing traffic usage", []string{"sensor_name", "sensor_id"}, nil),
		SensorLoad:        prometheus.NewDesc(prefix+"load", "Sensors load", []string{"sensor_name", "sensor_id"}, nil),
		SensorCpu:         prometheus.NewDesc(prefix+"cpu", "Sensors CPU usage", []string{"sensor_name", "sensor_id"}, nil),
		SensorRam:         prometheus.NewDesc(prefix+"ram", "Sensors ram usage", []string{"sensor_name", "sensor_id"}, nil),
	}
}

func (c *SensorsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.SensorInternalIPS
	ch <- c.SensorExternalIPS
	ch <- c.SensorPPSIn
	ch <- c.SensorPPSOut
	ch <- c.SensorBPSIn
	ch <- c.SensorBPSOut
	ch <- c.SensorDroppedIn
	ch <- c.SensorDroppedOut
	ch <- c.SensorUsageIn
	ch <- c.SensorUsageOut
	ch <- c.SensorLoad
	ch <- c.SensorCpu
	ch <- c.SensorRam
}

func (c *SensorsCollector) Collect(ch chan<- prometheus.Metric) {
	var sensors []Sensor

	err := c.wgClient.GetParsed("sensor_live_stats", &sensors)
	if err != nil {
		log.Errorln(err.Error())
	}

	for _, s := range sensors {
		ch <- prometheus.MustNewConstMetric(c.SensorInternalIPS, prometheus.GaugeValue, stringToFloat64(s.InternalIPS), s.Sensor.InterfaceName, s.Sensor.InterfaceID)
		ch <- prometheus.MustNewConstMetric(c.SensorExternalIPS, prometheus.GaugeValue, stringToFloat64(s.ExternalIPS), s.Sensor.InterfaceName, s.Sensor.InterfaceID)
		ch <- prometheus.MustNewConstMetric(c.SensorPPSIn, prometheus.GaugeValue, stringToFloat64(s.PacketsPerSecondIN), s.Sensor.InterfaceName, s.Sensor.InterfaceID)
		ch <- prometheus.MustNewConstMetric(c.SensorPPSOut, prometheus.GaugeValue, stringToFloat64(s.PacketsPerSecondOUT), s.Sensor.InterfaceName, s.Sensor.InterfaceID)
		ch <- prometheus.MustNewConstMetric(c.SensorBPSIn, prometheus.GaugeValue, bitsToBytes(stringToFloat64(s.BitsPerSecondIN)), s.Sensor.InterfaceName, s.Sensor.InterfaceID)
		ch <- prometheus.MustNewConstMetric(c.SensorBPSOut, prometheus.GaugeValue, bitsToBytes(stringToFloat64(s.BitsPerSecondOUT)), s.Sensor.InterfaceName, s.Sensor.InterfaceID)
		ch <- prometheus.MustNewConstMetric(c.SensorDroppedIn, prometheus.GaugeValue, stringToFloat64(s.DroppedIN), s.Sensor.InterfaceName, s.Sensor.InterfaceID)
		ch <- prometheus.MustNewConstMetric(c.SensorDroppedOut, prometheus.GaugeValue, stringToFloat64(s.DroppedOUT), s.Sensor.InterfaceName, s.Sensor.InterfaceID)
		ch <- prometheus.MustNewConstMetric(c.SensorUsageIn, prometheus.GaugeValue, stringToFloat64(s.UsageIN), s.Sensor.InterfaceName, s.Sensor.InterfaceID)
		ch <- prometheus.MustNewConstMetric(c.SensorUsageOut, prometheus.GaugeValue, stringToFloat64(s.UsageOUT), s.Sensor.InterfaceName, s.Sensor.InterfaceID)
		ch <- prometheus.MustNewConstMetric(c.SensorLoad, prometheus.GaugeValue, stringToFloat64(s.Load), s.Sensor.InterfaceName, s.Sensor.InterfaceID)
		ch <- prometheus.MustNewConstMetric(c.SensorCpu, prometheus.GaugeValue, stringToFloat64(s.Cpu), s.Sensor.InterfaceName, s.Sensor.InterfaceID)
		ch <- prometheus.MustNewConstMetric(c.SensorRam, prometheus.GaugeValue, float64(s.Ram), s.Sensor.InterfaceName, s.Sensor.InterfaceID)
	}
}
