package service

import (
	"net/http"
	"sync"
	"time"

	"github.com/SebastianCzoch/lxc-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const subsystem = "exporter"

var (
	collectorLabelNames = []string{"collector", "result"}

	scrapeDurations = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: collector.Namespace,
			Subsystem: subsystem,
			Name:      "scrape_duration_seconds",
			Help:      "node_exporter: Duration of a scrape job.",
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
	collectors := map[string]collector.Collector{}
	collectors["lxc_cpu"], _ = collector.NewStatCollector()

	nodeCollector := NodeCollector{collectors: collectors}
	prometheus.MustRegister(nodeCollector)
}

// NodeCollector implements the prometheus.Collector interface.
type NodeCollector struct {
	collectors map[string]collector.Collector
}

// Describe implements the prometheus.Collector interface.
func (n NodeCollector) Describe(ch chan<- *prometheus.Desc) {
	scrapeDurations.Describe(ch)
}

// Collect implements the prometheus.Collector interface.
func (n NodeCollector) Collect(ch chan<- prometheus.Metric) {
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
	var result string

	if err != nil {
		log.Errorf("ERROR: %s collector failed after %fs: %s", name, duration.Seconds(), err)
		result = "error"
	} else {
		log.Debugf("OK: %s collector succeeded after %fs.", name, duration.Seconds())
		result = "success"
	}
	scrapeDurations.WithLabelValues(name, result).Observe(duration.Seconds())
}
