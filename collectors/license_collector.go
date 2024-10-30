package collectors

import (
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	wgc "github.com/tomvil/wanguard_exporter/client"
)

type LicenseCollector struct {
	wgClient                       *wgc.Client
	SoftwareVersion                *prometheus.Desc
	LicensedSensors                *prometheus.Desc
	LicensedSensorsUsed            *prometheus.Desc
	LicensedSensorsRemaining       *prometheus.Desc
	LicensedDpdkEngines            *prometheus.Desc
	LicensedDpdkEnginesUsed        *prometheus.Desc
	LicensedDpdkEnginesRemaining   *prometheus.Desc
	LicensedFilters                *prometheus.Desc
	LicensedFiltersUsed            *prometheus.Desc
	LicensedFiltersRemaining       *prometheus.Desc
	LicenseSecondsRemaining        *prometheus.Desc
	LicenseSupportSecondsRemaining *prometheus.Desc
}

type License struct {
	SoftwareVersion              string      `json:"software_version"`
	LicensedSensors              interface{} `json:"licensed_sensors"`
	LicensedSensorsUsed          interface{} `json:"licensed_sensors_used"`
	LicensedSensorsRemaining     interface{} `json:"licensed_sensors_remaining"`
	LicensedDpdkEngines          interface{} `json:"licensed_dpdk_engines"`
	LicensedDpdkEnginesUsed      interface{} `json:"licensed_dpdk_engines_used"`
	LicensedDpdkEnginesRemaining interface{} `json:"licensed_dpdk_engines_remaining"`
	LicensedFilters              interface{} `json:"licensed_filters"`
	LicensedFiltersUsed          interface{} `json:"licensed_filters_used"`
	LicensedFiltersRemaining     interface{} `json:"licensed_filters_remaining"`
	LicenseDaysRemaining         interface{} `json:"license_expiry_date_remaining"`
	LicenseSupportDaysRemaining  interface{} `json:"support_expiry_date_remaining"`
}

func NewLicenseCollector(wgclient *wgc.Client) *LicenseCollector {
	prefix := "wanguard_license_"
	return &LicenseCollector{
		wgClient:                       wgclient,
		SoftwareVersion:                prometheus.NewDesc(prefix+"software_version", "Software version", []string{"software_version"}, nil),
		LicensedSensors:                prometheus.NewDesc(prefix+"sensors_count", "Licensed sensors count", nil, nil),
		LicensedSensorsUsed:            prometheus.NewDesc(prefix+"sensors_used", "Licensed sensors used", nil, nil),
		LicensedSensorsRemaining:       prometheus.NewDesc(prefix+"sensors_remaining", "Licensed sensors remaining", nil, nil),
		LicensedDpdkEngines:            prometheus.NewDesc(prefix+"dpdk_engines_count", "Licensed DPDK engines count", nil, nil),
		LicensedDpdkEnginesUsed:        prometheus.NewDesc(prefix+"dpdk_engines_used", "Licensed DPDK engines used", nil, nil),
		LicensedDpdkEnginesRemaining:   prometheus.NewDesc(prefix+"dpdk_engines_remaining", "Licensed DPDK engines remaining", nil, nil),
		LicensedFilters:                prometheus.NewDesc(prefix+"filters_count", "Licensed filters count", nil, nil),
		LicensedFiltersUsed:            prometheus.NewDesc(prefix+"filters_used", "Licensed filters used", nil, nil),
		LicensedFiltersRemaining:       prometheus.NewDesc(prefix+"filters_remaining", "Licensed filters available", nil, nil),
		LicenseSecondsRemaining:        prometheus.NewDesc(prefix+"seconds_remaining", "License seconds remaining", nil, nil),
		LicenseSupportSecondsRemaining: prometheus.NewDesc(prefix+"support_seconds_remaining", "Support license seconds remaining", nil, nil),
	}
}

func (c *LicenseCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.SoftwareVersion
	ch <- c.LicensedSensors
	ch <- c.LicensedSensorsUsed
	ch <- c.LicensedSensorsRemaining
	ch <- c.LicensedDpdkEngines
	ch <- c.LicensedDpdkEnginesUsed
	ch <- c.LicensedDpdkEnginesRemaining
	ch <- c.LicensedFilters
	ch <- c.LicensedFiltersUsed
	ch <- c.LicensedFiltersRemaining
	ch <- c.LicenseSecondsRemaining
	ch <- c.LicenseSupportSecondsRemaining
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

	ch <- prometheus.MustNewConstMetric(c.SoftwareVersion, prometheus.GaugeValue, 1, license.SoftwareVersion)
	ch <- prometheus.MustNewConstMetric(c.LicensedSensors, prometheus.GaugeValue, getFloat64(license.LicensedSensors))
	ch <- prometheus.MustNewConstMetric(c.LicensedSensorsUsed, prometheus.GaugeValue, getFloat64(license.LicensedSensorsUsed))
	ch <- prometheus.MustNewConstMetric(c.LicensedSensorsRemaining, prometheus.GaugeValue, getFloat64(license.LicensedSensorsRemaining))
	ch <- prometheus.MustNewConstMetric(c.LicensedDpdkEngines, prometheus.GaugeValue, getFloat64(license.LicensedDpdkEngines))
	ch <- prometheus.MustNewConstMetric(c.LicensedDpdkEnginesUsed, prometheus.GaugeValue, getFloat64(license.LicensedDpdkEnginesUsed))
	ch <- prometheus.MustNewConstMetric(c.LicensedDpdkEnginesRemaining, prometheus.GaugeValue, getFloat64(license.LicensedDpdkEnginesRemaining))
	ch <- prometheus.MustNewConstMetric(c.LicensedFilters, prometheus.GaugeValue, getFloat64(license.LicensedFilters))
	ch <- prometheus.MustNewConstMetric(c.LicensedFiltersUsed, prometheus.GaugeValue, getFloat64(license.LicensedFiltersUsed))
	ch <- prometheus.MustNewConstMetric(c.LicensedFiltersRemaining, prometheus.GaugeValue, getFloat64(license.LicensedFiltersRemaining))
	ch <- prometheus.MustNewConstMetric(c.LicenseSecondsRemaining, prometheus.GaugeValue, toSeconds(getFloat64(license.LicenseDaysRemaining)))
	ch <- prometheus.MustNewConstMetric(c.LicenseSupportSecondsRemaining, prometheus.GaugeValue, toSeconds(getFloat64(license.LicenseSupportDaysRemaining)))
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

func toSeconds(days float64) float64 {
	return days * 86400
}
