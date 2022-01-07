package collectors

import (
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	wgc "github.com/tomvil/wanguard_exporter/client"
)

type LicenseCollector struct {
	wgClient                    *wgc.Client
	LicensedSensors             *prometheus.Desc
	LicensedSensorsUsed         *prometheus.Desc
	LicensedSensorsRemaining    *prometheus.Desc
	LicensedFilters             *prometheus.Desc
	LicensedFiltersUsed         *prometheus.Desc
	LicensedFiltersRemaining    *prometheus.Desc
	LicenseDaysRemaining        *prometheus.Desc
	LicenseSupportDaysRemaining *prometheus.Desc
}

type License struct {
	LicensedSensors             interface{} `json:"licensed_sensors"`
	LicensedSensorsUsed         interface{} `json:"licensed_sensors_used"`
	LicensedSensorsRemaining    interface{} `json:"licensed_sensors_remaining"`
	LicensedFilters             interface{} `json:"licensed_filters"`
	LicensedFiltersUsed         interface{} `json:"licensed_filters_used"`
	LicensedFiltersRemaining    interface{} `json:"licensed_filters_remaining"`
	LicenseDaysRemaining        interface{} `json:"license_expiry_date_remaining"`
	LicenseSupportDaysRemaining interface{} `json:"support_expiry_date_remaining"`
}

func NewLicenseCollector(wgclient *wgc.Client) *LicenseCollector {
	prefix := "wanguard_license_"
	return &LicenseCollector{
		wgClient:                    wgclient,
		LicensedSensors:             prometheus.NewDesc(prefix+"sensors_total", "Licensed sensors total", nil, nil),
		LicensedSensorsUsed:         prometheus.NewDesc(prefix+"sensors_used", "Licensed sensors used", nil, nil),
		LicensedSensorsRemaining:    prometheus.NewDesc(prefix+"sensors_remaining", "Licensed sensors remaining", nil, nil),
		LicensedFilters:             prometheus.NewDesc(prefix+"filters", "Licensed filters total", nil, nil),
		LicensedFiltersUsed:         prometheus.NewDesc(prefix+"filters_used", "Licensed filters total", nil, nil),
		LicensedFiltersRemaining:    prometheus.NewDesc(prefix+"filters_remaining", "Licensed filters available", nil, nil),
		LicenseDaysRemaining:        prometheus.NewDesc(prefix+"days_remaining", "License days remaining", nil, nil),
		LicenseSupportDaysRemaining: prometheus.NewDesc(prefix+"support_days_remaining", "Support license days remaining", nil, nil),
	}
}

func (c *LicenseCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.LicensedSensors
	ch <- c.LicensedSensorsUsed
	ch <- c.LicensedSensorsRemaining
	ch <- c.LicensedFilters
	ch <- c.LicensedFiltersUsed
	ch <- c.LicensedFiltersRemaining
	ch <- c.LicenseDaysRemaining
	ch <- c.LicenseSupportDaysRemaining
}

func (c *LicenseCollector) Collect(ch chan<- prometheus.Metric) {
	var license License

	err := c.wgClient.GetParsed(
		"license_manager",
		&license,
	)
	if err != nil {
		log.Errorln(
			err.Error(),
		)
	}

	ch <- prometheus.MustNewConstMetric(c.LicensedSensors, prometheus.GaugeValue, getFloat64(license.LicensedSensors))
	ch <- prometheus.MustNewConstMetric(c.LicensedSensorsUsed, prometheus.GaugeValue, getFloat64(license.LicensedSensorsUsed))
	ch <- prometheus.MustNewConstMetric(c.LicensedSensorsRemaining, prometheus.GaugeValue, getFloat64(license.LicensedSensorsRemaining))
	ch <- prometheus.MustNewConstMetric(c.LicensedFilters, prometheus.GaugeValue, getFloat64(license.LicensedFilters))
	ch <- prometheus.MustNewConstMetric(c.LicensedFiltersUsed, prometheus.GaugeValue, getFloat64(license.LicensedFiltersUsed))
	ch <- prometheus.MustNewConstMetric(c.LicensedFiltersRemaining, prometheus.GaugeValue, getFloat64(license.LicensedFiltersRemaining))
	ch <- prometheus.MustNewConstMetric(c.LicenseDaysRemaining, prometheus.GaugeValue, getFloat64(license.LicenseDaysRemaining))
	ch <- prometheus.MustNewConstMetric(c.LicenseSupportDaysRemaining, prometheus.GaugeValue, getFloat64(license.LicenseSupportDaysRemaining))
}

func getFloat64(v interface{}) float64 {
	switch v := v.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case nil:
		return 0
	case string:
		r := strings.NewReplacer(
			" days", "",
			"∞", "9999")

		result, err := strconv.ParseFloat(r.Replace(v), 64)
		if err != nil {
			log.Errorf("was not able to parse %T to float64!", v)
			return 0
		}

		return float64(result)
	default:
		log.Errorf("conversion to float64 from %T is not supported", v)
		return 0
	}
}
