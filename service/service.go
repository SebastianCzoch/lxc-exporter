package service

import (
	"net/http"
	"sync"
	"time"

	"github.com/SebastianCzoch/lxc-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	collectors          = make(map[string]collector.Collector)
	collectorLabelNames = []string{"collector", "result"}

	scrapeDurations = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: collector.Namespace,
			Subsystem: "exporter",
			Name:      "scrape_duration_seconds",
			Help:      "lxc_exporter: Duration of a scrape job.",
		},
		collectorLabelNames,
	)
)

func StartServer(addr *string) {
	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
            <head><title>LXC Exporter</title></head>
            <body>
            <h1>LXC Exporter</h1>
            <p><a href="/metrics">Metrics</a></p>
            </body>
            </html>`))
	})

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}

func StartColectors() {
	loadCollectors()
	for name, c := range collectors {
		err := c.Init()
		if err != nil {
			delete(collectors, name)
			log.Errorf("collector %s disabled, because: %s", name, err.Error())
		}
	}

	prometheus.MustRegister(LXCCollector{collectors: collectors})
}

// LXCCollector implements the prometheus.Collector interface.
type LXCCollector struct {
	collectors map[string]collector.Collector
}

// Describe implements the prometheus.Collector interface.
func (n LXCCollector) Describe(ch chan<- *prometheus.Desc) {
	scrapeDurations.Describe(ch)
}

// Collect implements the prometheus.Collector interface.
func (n LXCCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(n.collectors))
	for name, c := range n.collectors {
		go func(name string, c collector.Collector) {
			execute(name, c, ch)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
	scrapeDurations.Collect(ch)
}

func execute(name string, c collector.Collector, ch chan<- prometheus.Metric) {
	begin := time.Now()
	err := c.Update(ch)
	duration := time.Since(begin)

	if err != nil {
		log.Errorf("ERROR: %s collector failed after %fs: %s", name, duration.Seconds(), err)
		scrapeDurations.WithLabelValues(name, "error").Observe(duration.Seconds())
		return
	}

	scrapeDurations.WithLabelValues(name, "success").Observe(duration.Seconds())
}

func loadCollectors() {
	collectors["lxc_cpu"] = collector.NewCPUStatCollector()
	collectors["lxc_mem"] = collector.NewMemStatCollector()
}
