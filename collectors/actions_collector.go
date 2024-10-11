package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	wgc "github.com/tomvil/wanguard_exporter/client"
)

type ActionsCollector struct {
	wgClient     *wgc.Client
	ActionStatus *prometheus.Desc
}

type Response struct {
	ResponseName string `json:"response_name"`
	ResponseHref string `json:"href"`
}

type Action struct {
	ActionName     string `json:"action_name"`
	ActionType     string `json:"action_type"`
	ResponseBranch string `json:"response_branch"`
	ActionHref     string `json:"href"`
}

func NewActionsCollector(wgclient *wgc.Client) *ActionsCollector {
	prefix := "wanguard_action_"
	return &ActionsCollector{
		wgClient:     wgclient,
		ActionStatus: prometheus.NewDesc(prefix+"status", "Status of the response actions", []string{"response_name", "action_name", "action_type", "response_branch"}, nil),
	}
}

func (c *ActionsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.ActionStatus
}

func (c *ActionsCollector) Collect(ch chan<- prometheus.Metric) {
	var responses []Response

	err := c.wgClient.GetParsed("responses", &responses)
	if err != nil {
		return
	}
	for _, response := range responses {
		var actions []Action

		err := c.wgClient.GetParsed(response.ResponseHref+"/actions", &actions)
		if err != nil {
			continue
		}

		for _, action := range actions {
			var params map[string]string

			err := c.wgClient.GetParsed(action.ActionHref+"/status", &params)
			if err != nil {
				log.Errorln(err.Error())
				continue
			}

			if params["status"] == "Active" {
				ch <- prometheus.MustNewConstMetric(c.ActionStatus, prometheus.GaugeValue, 1, response.ResponseName, action.ActionName, action.ActionType, action.ResponseBranch)
			} else {
				ch <- prometheus.MustNewConstMetric(c.ActionStatus, prometheus.GaugeValue, 0, response.ResponseName, action.ActionName, action.ActionType, action.ResponseBranch)
			}

		}
	}
}
