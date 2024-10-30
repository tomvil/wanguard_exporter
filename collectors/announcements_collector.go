package collectors

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	wgc "github.com/tomvil/wanguard_exporter/client"
)

type AnnouncementsCollector struct {
	wgClient                   *wgc.Client
	AnnouncementActive         *prometheus.Desc
	AnnouncementsFinishedCount *prometheus.Desc
}

type AnnouncementsCount struct {
	Count string
}

type Announcement struct {
	Id     string `json:"bgp_announcement_id"`
	Prefix string
	From   Time
	Until  Time
}

func NewAnnouncementsCollector(wgclient *wgc.Client) *AnnouncementsCollector {
	prefix := "wanguard_announcements_"
	return &AnnouncementsCollector{
		wgClient:                   wgclient,
		AnnouncementActive:         prometheus.NewDesc(prefix+"active", "Active announcements at the moment", []string{"prefix", "from", "until", "announcement_id"}, nil),
		AnnouncementsFinishedCount: prometheus.NewDesc(prefix+"count", "Total count of announcements", nil, nil),
	}
}

func (c *AnnouncementsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.AnnouncementActive
	ch <- c.AnnouncementsFinishedCount
}

func (c *AnnouncementsCollector) Collect(ch chan<- prometheus.Metric) {
	collectActiveAnnouncements(c.AnnouncementActive, c.wgClient, ch)
	collectFinishedAnnouncementsTotal(c.AnnouncementsFinishedCount, c.wgClient, ch)
}

func collectActiveAnnouncements(desc *prometheus.Desc, wgclient *wgc.Client, ch chan<- prometheus.Metric) {
	var announcements []Announcement

	err := wgclient.GetParsed("bgp_announcements?status=Active&fields=prefix,from,until,bgp_announcement_id", &announcements)
	if err != nil {
		return
	}

	for _, announcement := range announcements {
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, 1, announcement.Prefix, announcement.From.Time, announcement.Until.Time, announcement.Id)
	}
}

func collectFinishedAnnouncementsTotal(desc *prometheus.Desc, wgclient *wgc.Client, ch chan<- prometheus.Metric) {
	var finishedAnnouncementsCount AnnouncementsCount

	err := wgclient.GetParsed("bgp_announcements?status=Finished&count=true", &finishedAnnouncementsCount)
	if err != nil {
		log.Errorln(err.Error())
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, 0)
		return
	}

	r, err := strconv.ParseFloat(finishedAnnouncementsCount.Count, 64)
	if err != nil {
		log.Errorln(err.Error())
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, 0)
		return
	}

	ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, r)
}
