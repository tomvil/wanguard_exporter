package collectors

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	wgc "github.com/tomvil/wanguard_exporter/client"
)

type FirewallRulesCollector struct {
	wgClient           *wgc.Client
	FirewallRuleActive *prometheus.Desc
	FirewallRulesTotal *prometheus.Desc
}

type FirewallRulesCount struct {
	Count string
}

type FirewallRule struct {
	Rule_id            string `json:"firewall_rule_id"`
	Attack_id          string
	Source_prefix      string
	Destination_prefix string
	Ip_protocol        string
	From               Time
	Until              Time
	Pkts_s             string `json:"pkts/s"`
	Bits_s             string `json:"bits/s"`
	Max_pkts_s         string `json:"max_pkts/s"`
	Max_bits_s         string `json:"max_bits/s"`
	Pkts               string
	Bits               string
}

func NewFirewallRulesCollector(wgclient *wgc.Client) *FirewallRulesCollector {
	prefix := "wanguard_firewall_rules_"
	return &FirewallRulesCollector{
		wgClient: wgclient,
		FirewallRuleActive: prometheus.NewDesc(prefix+"active", "Active firewall rules at the moment",
			[]string{"rule_id",
				"attack_id",
				"source_prefix",
				"destination_prefix",
				"ip_protocol",
				"from",
				"until",
				"pkts_s",
				"bits_s",
				"max_pkts_s",
				"max_bits_s",
				"pkts",
				"bits"}, nil),
		FirewallRulesTotal: prometheus.NewDesc(prefix+"total", "Total amount of active firewall rules", nil, nil),
	}
}

func (c *FirewallRulesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.FirewallRuleActive
	ch <- c.FirewallRulesTotal
}

func (c *FirewallRulesCollector) Collect(ch chan<- prometheus.Metric) {
	collectFirewallRules(c.FirewallRuleActive, c.wgClient, ch)
	collectFirewallRulesTotal(c.FirewallRulesTotal, c.wgClient, ch)
}

func collectFirewallRules(desc *prometheus.Desc, wgclient *wgc.Client, ch chan<- prometheus.Metric) {
	var rules []map[string]string

	err := wgclient.GetParsed("firewall_rules", &rules)
	if err != nil {
		return
	}

	for _, r := range rules {
		var rule FirewallRule
		err = wgclient.GetParsed(r["href"], &rule)
		if err != nil {
			return
		}

		ch <- prometheus.MustNewConstMetric(
			desc,
			prometheus.GaugeValue,
			1,
			rule.Rule_id,
			rule.Attack_id,
			rule.Source_prefix,
			rule.Destination_prefix,
			rule.Ip_protocol,
			rule.From.Time,
			rule.Until.Time,
			rule.Pkts_s,
			rule.Bits_s,
			rule.Max_pkts_s,
			rule.Max_bits_s,
			rule.Pkts,
			rule.Bits)
	}
}

func collectFirewallRulesTotal(desc *prometheus.Desc, wgclient *wgc.Client, ch chan<- prometheus.Metric) {
	var firewallRulesCount FirewallRulesCount

	err := wgclient.GetParsed("firewall_rules?count=true", &firewallRulesCount)
	if err != nil {
		log.Errorln(err.Error())
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, 0)
		return
	}

	c, err := strconv.ParseFloat(firewallRulesCount.Count, 64)
	if err != nil {
		log.Errorln(err.Error())
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, 0)
		return
	}

	ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, c)
}
