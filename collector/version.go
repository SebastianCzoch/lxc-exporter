package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	Version  string
	Revision string
	Branch   string
)

type versionCollector struct {
	metric *prometheus.GaugeVec
}

func init() {
	Factories["version"] = NewVersionCollector
}

func NewVersionCollector() (Collector, error) {
	metric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "lxc_exporter_build_info",
			Help: "A metric with a constant '1' value labeled by version, revision, and branch from which the lxc_exporter was built.",
		},
		[]string{"version", "revision", "branch"},
	)
	metric.WithLabelValues(Version, Revision, Branch).Set(1)
	return &versionCollector{
		metric: metric,
	}, nil
}

func (c *versionCollector) Update(ch chan<- prometheus.Metric) (err error) {
	c.metric.Collect(ch)
	return err
}
