package collectors

import (
	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/tomvil/countries"
	wgc "github.com/tomvil/wanguard_exporter/client"
)

type TrafficCollector struct {
	wgClient            *wgc.Client
	CountryTopPPSIn     *prometheus.Desc
	CountryTopPPSOut    *prometheus.Desc
	CountryTopBPSIn     *prometheus.Desc
	CountryTopBPSOut    *prometheus.Desc
	IPVersionTopPPSIn   *prometheus.Desc
	IPVersionTopPPSOut  *prometheus.Desc
	IPVersionTopBPSIn   *prometheus.Desc
	IPVersionTopBPSOut  *prometheus.Desc
	IPProtocolTopPPSIn  *prometheus.Desc
	IPProtocolTopPPSOut *prometheus.Desc
	IPProtocolTopBPSIn  *prometheus.Desc
	IPProtocolTopBPSOut *prometheus.Desc
	TalkersTopPPSIn     *prometheus.Desc
	TalkersTopPPSOut    *prometheus.Desc
	TalkersTopBPSIn     *prometheus.Desc
	TalkersTopBPSOut    *prometheus.Desc
}

type CountryTop struct {
	Top map[string]Country
}

type Country struct {
	Country string
	Value   int
}

type IPVersionTop struct {
	Top map[string]IPVersion
}

type IPVersion struct {
	IPVersion string `json:"Description"`
	Value     int
}

type IPProtocolTop struct {
	Top map[string]IPProtocol
}

type IPProtocol struct {
	IPProtocol int `json:"ip_protocol"`
	Value      int
}

type TalkerTop struct {
	Top map[string]Talker
}

type Talker struct {
	IPAddress string `json:"ip_address"`
	Value     int
}

var protocol = map[int]string{
	1:  "ICMP",
	2:  "IGMP",
	6:  "TCP",
	17: "UDP",
	47: "GRE",
	50: "ESP",
	51: "AH",
	58: "ICMPv6",
	88: "EIGRP",
	89: "OSPF",
}

func NewTrafficCollector(wgclient *wgc.Client) *TrafficCollector {
	prefix := "wanguard_traffic_"
	return &TrafficCollector{
		wgClient:            wgclient,
		CountryTopPPSIn:     prometheus.NewDesc(prefix+"country_packets_per_second_in", "Packets per second in by country", []string{"country", "country_code"}, nil),
		CountryTopPPSOut:    prometheus.NewDesc(prefix+"country_packets_per_second_out", "Packets per second out by country", []string{"country", "country_code"}, nil),
		CountryTopBPSIn:     prometheus.NewDesc(prefix+"country_bits_per_second_in", "Bits per second in by country", []string{"country", "country_code"}, nil),
		CountryTopBPSOut:    prometheus.NewDesc(prefix+"country_bits_per_second_out", "Bits per second out by country", []string{"country", "country_code"}, nil),
		IPVersionTopPPSIn:   prometheus.NewDesc(prefix+"ip_version_packets_per_second_in", "Packets per second in by IP version", []string{"ip_version"}, nil),
		IPVersionTopPPSOut:  prometheus.NewDesc(prefix+"ip_version_packets_per_second_out", "Packets per second out by IP version", []string{"ip_version"}, nil),
		IPVersionTopBPSIn:   prometheus.NewDesc(prefix+"ip_version_bits_per_second_in", "Bits per second in by IP version", []string{"ip_version"}, nil),
		IPVersionTopBPSOut:  prometheus.NewDesc(prefix+"ip_version_bits_per_second_out", "Bits per second out by IP version", []string{"ip_version"}, nil),
		IPProtocolTopPPSIn:  prometheus.NewDesc(prefix+"ip_protocol_packets_per_second_in", "Packets per second in by IP protocol", []string{"ip_protocol"}, nil),
		IPProtocolTopPPSOut: prometheus.NewDesc(prefix+"ip_protocol_packets_per_second_out", "Packets per second out by IP protocol", []string{"ip_protocol"}, nil),
		IPProtocolTopBPSIn:  prometheus.NewDesc(prefix+"ip_protocol_bits_per_second_in", "Bits per second in by IP protocol", []string{"ip_protocol"}, nil),
		IPProtocolTopBPSOut: prometheus.NewDesc(prefix+"ip_protocol_bits_per_second_out", "Bits per second out by IP protocol", []string{"ip_protocol"}, nil),
		TalkersTopPPSIn:     prometheus.NewDesc(prefix+"talkers_packets_per_second_in", "Packets per second in by IP address", []string{"ip_address"}, nil),
		TalkersTopPPSOut:    prometheus.NewDesc(prefix+"talkers_packets_per_second_out", "Packets per second out by IP address", []string{"ip_address"}, nil),
		TalkersTopBPSIn:     prometheus.NewDesc(prefix+"talkers_bits_per_second_in", "Bits per second in by IP address", []string{"ip_address"}, nil),
		TalkersTopBPSOut:    prometheus.NewDesc(prefix+"talkers_bits_per_second_out", "Bits per second out by IP address", []string{"ip_address"}, nil),
	}
}

func (c *TrafficCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.CountryTopPPSIn
	ch <- c.CountryTopPPSOut
	ch <- c.CountryTopBPSIn
	ch <- c.CountryTopBPSOut
	ch <- c.IPVersionTopPPSIn
	ch <- c.IPVersionTopPPSOut
	ch <- c.IPVersionTopBPSIn
	ch <- c.IPVersionTopBPSOut
	ch <- c.IPProtocolTopPPSIn
	ch <- c.IPProtocolTopPPSOut
	ch <- c.IPProtocolTopBPSIn
	ch <- c.IPProtocolTopBPSOut
	ch <- c.TalkersTopPPSIn
	ch <- c.TalkersTopPPSOut
	ch <- c.TalkersTopBPSIn
	ch <- c.TalkersTopBPSOut
}

func (c *TrafficCollector) Collect(ch chan<- prometheus.Metric) {
	var wsync sync.WaitGroup
	wsync.Add(16)

	go collectTopTrafficByCountry(c.CountryTopPPSIn, ch, c.wgClient, &wsync, "Packets", "Inbound")
	go collectTopTrafficByCountry(c.CountryTopBPSIn, ch, c.wgClient, &wsync, "Bits", "Inbound")
	go collectTopTrafficByCountry(c.CountryTopPPSOut, ch, c.wgClient, &wsync, "Packets", "Outbound")
	go collectTopTrafficByCountry(c.CountryTopBPSOut, ch, c.wgClient, &wsync, "Bits", "Outbound")

	go collectTopTrafficByIPVersion(c.IPVersionTopPPSIn, ch, c.wgClient, &wsync, "Packets", "Inbound")
	go collectTopTrafficByIPVersion(c.IPVersionTopBPSIn, ch, c.wgClient, &wsync, "Bits", "Inbound")
	go collectTopTrafficByIPVersion(c.IPVersionTopPPSOut, ch, c.wgClient, &wsync, "Packets", "Outbound")
	go collectTopTrafficByIPVersion(c.IPVersionTopBPSOut, ch, c.wgClient, &wsync, "Bits", "Outbound")

	go collectTopTrafficByIPProtocol(c.IPProtocolTopPPSIn, ch, c.wgClient, &wsync, "Packets", "Inbound")
	go collectTopTrafficByIPProtocol(c.IPProtocolTopBPSIn, ch, c.wgClient, &wsync, "Bits", "Inbound")
	go collectTopTrafficByIPProtocol(c.IPProtocolTopPPSOut, ch, c.wgClient, &wsync, "Packets", "Outbound")
	go collectTopTrafficByIPProtocol(c.IPProtocolTopBPSOut, ch, c.wgClient, &wsync, "Bits", "Outbound")

	go collectTopTrafficByTalkers(c.TalkersTopPPSIn, ch, c.wgClient, &wsync, "Packets", "Inbound")
	go collectTopTrafficByTalkers(c.TalkersTopBPSIn, ch, c.wgClient, &wsync, "Bits", "Inbound")
	go collectTopTrafficByTalkers(c.TalkersTopPPSOut, ch, c.wgClient, &wsync, "Packets", "Outbound")
	go collectTopTrafficByTalkers(c.TalkersTopBPSOut, ch, c.wgClient, &wsync, "Bits", "Outbound")

	wsync.Wait()

}

func collectTopTrafficByCountry(desc *prometheus.Desc, ch chan<- prometheus.Metric, wgclient *wgc.Client, wsync *sync.WaitGroup, unit string, direction string) {
	var countryTop CountryTop

	href := "sensor_live_tops?top_type=Countries" + "&unit=" + unit + "&direction=" + direction

	err := wgclient.GetParsed(href, &countryTop)
	if err != nil {
		log.Errorln(err.Error())
	}

	for i := 1; i <= len(countryTop.Top); i++ {
		k := strconv.Itoa(i)
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(countryTop.Top[k].Value), countryTop.Top[k].Country, countries.ByName(countryTop.Top[k].Country).Alpha2())
	}

	defer wsync.Done()
}

func collectTopTrafficByIPVersion(desc *prometheus.Desc, ch chan<- prometheus.Metric, wgclient *wgc.Client, wsync *sync.WaitGroup, unit string, direction string) {
	var ipVersionTop IPVersionTop

	href := "sensor_live_tops?top_type=IP%20Versions" + "&unit=" + unit + "&direction=" + direction

	err := wgclient.GetParsed(href, &ipVersionTop)
	if err != nil {
		log.Errorln(err.Error())
	}

	for i := 1; i <= len(ipVersionTop.Top); i++ {
		k := strconv.Itoa(i)
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(ipVersionTop.Top[k].Value), ipVersionTop.Top[k].IPVersion)
	}

	defer wsync.Done()
}

func collectTopTrafficByIPProtocol(desc *prometheus.Desc, ch chan<- prometheus.Metric, wgclient *wgc.Client, wsync *sync.WaitGroup, unit string, direction string) {
	var ipProtocolTop IPProtocolTop

	href := "sensor_live_tops?top_type=IP%20Protocols" + "&unit=" + unit + "&direction=" + direction

	err := wgclient.GetParsed(href, &ipProtocolTop)
	if err != nil {
		log.Errorln(err.Error())
	}

	for i := 1; i <= len(ipProtocolTop.Top); i++ {
		k := strconv.Itoa(i)
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(ipProtocolTop.Top[k].Value), protocol[ipProtocolTop.Top[k].IPProtocol])
	}

	defer wsync.Done()
}

func collectTopTrafficByTalkers(desc *prometheus.Desc, ch chan<- prometheus.Metric, wgclient *wgc.Client, wsync *sync.WaitGroup, unit string, direction string) {
	var talkerTop TalkerTop

	href := "sensor_live_tops?top_type=Talkers" + "&unit=" + unit + "&direction=" + direction

	err := wgclient.GetParsed(href, &talkerTop)
	if err != nil {
		log.Errorln(err.Error())
	}

	for i := 1; i <= len(talkerTop.Top); i++ {
		k := strconv.Itoa(i)
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(talkerTop.Top[k].Value), talkerTop.Top[k].IPAddress)
	}

	defer wsync.Done()
}
