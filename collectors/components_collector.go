package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	wgc "github.com/tomvil/wanguard_exporter/client"
)

type ComponentsCollector struct {
	wgClient             *wgc.Client
	ComponentsCategories []string
	ComponentStatus      *prometheus.Desc
}

func NewComponentsCollector(wgclient *wgc.Client) *ComponentsCollector {
	prefix := "wanguard_component_"
	return &ComponentsCollector{
		wgClient:             wgclient,
		ComponentsCategories: []string{"bgp_connector", "filter", "sensor"},
		ComponentStatus:      prometheus.NewDesc(prefix+"status", "Status of the component", []string{"component_name", "component_category"}, nil),
	}
}

func (c *ComponentsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.ComponentStatus
}

func (c *ComponentsCollector) Collect(ch chan<- prometheus.Metric) {
	for _, category := range c.ComponentsCategories {
		var components []map[string]string

		err := c.wgClient.GetParsed(category+"s", &components)
		if err != nil {
			continue
		}

		for _, component := range components {
			var params map[string]string

			err := c.wgClient.GetParsed(component["href"]+"/status", &params)
			if err != nil {
				log.Errorln(err.Error())
				continue
			}

			if params["status"] == "Active" {
				ch <- prometheus.MustNewConstMetric(c.ComponentStatus, prometheus.GaugeValue, 1, component[category+"_name"], category)
			} else {
				ch <- prometheus.MustNewConstMetric(c.ComponentStatus, prometheus.GaugeValue, 0, component[category+"_name"], category)
			}

		}
	}
}
