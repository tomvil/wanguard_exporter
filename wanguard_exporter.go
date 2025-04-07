package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	wgc "github.com/tomvil/wanguard_exporter/client"
	"github.com/tomvil/wanguard_exporter/collectors"
)

const (
	version string = "1.6"
)

type collectorsList struct {
	enabled   *bool
	collector prometheus.Collector
}

var (
	showVersion = flag.Bool("version", false, "Print version and other information about wanguard_exporter")
	listenAddr  = flag.String("web.listen-address", ":9868", "The address to listen on for HTTP requests")
	metricsPath = flag.String("web.metrics-path", "/metrics", "Path under which metrics will be exposed")
	apiAddress  = flag.String("api.address", "http://127.0.0.1:81", "WANGuard API address")
	apiUsername = flag.String("api.username", "admin", "WANGuard API username")
	apiPassword = flag.String("api.password", "", "WANGuard API password")

	licenseCollectorEnabled       = flag.Bool("collector.license", true, "Expose license metrics")
	announcementsCollectorEnabled = flag.Bool("collector.announcements", true, "Expose announcements metrics")
	anomaliesCollectorEnabled     = flag.Bool("collector.anomalies", true, "Expose anomalies metrics")
	componentsCollectorEnabled    = flag.Bool("collector.components", true, "Expose components metrics")
	actionsCollectorEnabled       = flag.Bool("collector.actions", true, "Expose actions metrics")
	sensorsCollectorEnabled       = flag.Bool("collector.sensors", true, "Expose sensors metrics")
	trafficCollectorEnabled       = flag.Bool("collector.traffic", true, "Expose traffic metrics")
	firewallRulesCollectorEnabled = flag.Bool("collector.firewall_rules", true, "Expose firewall rules metrics")

	cl []collectorsList
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Println("wanguard_exporter")
		fmt.Println("Version:", version)
		fmt.Println("Author: Tomas Vilemaitis")
		fmt.Println("Metric exporter for WANGuard")
		os.Exit(0)
	}

	if *apiPassword == "" {
		*apiPassword = os.Getenv("WANGUARD_PASSWORD")
		if *apiPassword == "" {
			log.Errorln(`Please set the WANGuard API Password!
		API Password can be set with api.password flag or
		by setting WANGUARD_PASSWORD environment variable.`)
			os.Exit(1)
		}
	}

	wgClient := wgc.NewClient(*apiAddress, *apiUsername, *apiPassword)
	cl = []collectorsList{
		{licenseCollectorEnabled, collectors.NewLicenseCollector(wgClient)},
		{announcementsCollectorEnabled, collectors.NewAnnouncementsCollector(wgClient)},
		{anomaliesCollectorEnabled, collectors.NewAnomaliesCollector(wgClient)},
		{componentsCollectorEnabled, collectors.NewComponentsCollector(wgClient)},
		{actionsCollectorEnabled, collectors.NewActionsCollector(wgClient)},
		{sensorsCollectorEnabled, collectors.NewSensorsCollector(wgClient)},
		{trafficCollectorEnabled, collectors.NewTrafficCollector(wgClient)},
		{firewallRulesCollectorEnabled, collectors.NewFirewallRulesCollector(wgClient)},
	}

	startServer()
}

func startServer() {
	var landingPage = []byte(`<html>
	<head><title>WANguard Exporter (Version ` + version + `)</title></head>
	<body>
	<h1>WANGuard Exporter</h1>
	<p><a href="` + *metricsPath + `">Metrics</a></p>
	<h2>More information:</h2>
	<p><a href="https://github.com/tomvil/wanguard_exporter">github.com/tomvil/wanguard_exporter</a></p>
	</body>
	</html>`)

	log.Infof("Starting WANGuard exporter (Version: %s)", version)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write(landingPage); err != nil {
			log.Fatal(err.Error())
		}
	})
	http.HandleFunc(*metricsPath, handleMetricsRequest)

	log.Infof("Listening for %s on %s", *metricsPath, *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func handleMetricsRequest(w http.ResponseWriter, r *http.Request) {
	registry := prometheus.NewRegistry()

	for _, c := range cl {
		if *c.enabled {
			registry.MustRegister(c.collector)
		}
	}

	promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorLog:      log.NewErrorLogger(),
		ErrorHandling: promhttp.ContinueOnError}).ServeHTTP(w, r)
}
