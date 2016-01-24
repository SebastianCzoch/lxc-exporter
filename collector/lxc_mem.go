package collector

import (
	"github.com/SebastianCzoch/lxc-exporter/kernel"
	"github.com/SebastianCzoch/lxc-exporter/lxc"
	"github.com/prometheus/client_golang/prometheus"
)

type lxcMemCollector struct {
	memory  *prometheus.Desc
	lxcStat *lxc.LXC
}

// NewMemStatCollector is a method which return new lxc memory collector
func NewMemStatCollector() Collector {
	return &lxcMemCollector{
		memory: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "memory_usage"),
			"Memory usage in each container in bytes.",
			[]string{"container"}, nil,
		),
	}
}

func (c *lxcMemCollector) Init() error {
	kernelVersion, err := kernel.GetMajorVersion()
	if err != nil {
		return err
	}

	c.lxcStat, err = lxc.New(kernelVersion)
	return err
}

func (c *lxcMemCollector) Update(ch chan<- prometheus.Metric) error {
	for _, containerName := range c.lxcStat.GetContainers() {
		usage, err := c.lxcStat.GetMemStat(containerName)
		if err != nil {
			continue
		}

		ch <- prometheus.MustNewConstMetric(c.memory, prometheus.CounterValue, usage.Usage, containerName)
	}

	return nil
}
