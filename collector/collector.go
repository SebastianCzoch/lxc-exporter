package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Namespace is a prefix for all metrics
const Namespace = "lxc"

// Collector is interface for collectors
type Collector interface {
	Update(ch chan<- prometheus.Metric) error
	Init() error
}
