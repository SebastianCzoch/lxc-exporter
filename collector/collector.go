package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

const Namespace = "lxc"

var Factories = make(map[string]func() (Collector, error))

type Collector interface {
	Update(ch chan<- prometheus.Metric) error
	Init() error
}
