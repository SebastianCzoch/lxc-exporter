package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

const Namespace = "lxc"

type Collector interface {
	Update(ch chan<- prometheus.Metric) error
	Init() error
}
