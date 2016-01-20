package collector

import (
	"github.com/SebastianCzoch/lxc-exporter/kernel"
	"github.com/SebastianCzoch/lxc-exporter/lxc"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	userHz = 100
)

type lxcCPUCollector struct {
	cpu          *prometheus.Desc
	intr         *prometheus.Desc
	ctxt         *prometheus.Desc
	forks        *prometheus.Desc
	btime        *prometheus.Desc
	procsRunning *prometheus.Desc
	procsBlocked *prometheus.Desc
}

func init() {
	Factories["stat"] = NewStatCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// kernel/system statistics.
func NewStatCollector() (Collector, error) {
	return &lxcCPUCollector{
		cpu: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "cpu"),
			"Seconds the cpus spent in each mode.",
			[]string{"cpu", "mode", "container"}, nil,
		),
	}, nil
}

// Expose kernel and system statistics.
func (c *lxcCPUCollector) Update(ch chan<- prometheus.Metric) error {
	kernelVersion, _ := kernel.GetMajorVersion()
	lxcstat, _ := lxc.New(kernelVersion)

	for _, containerName := range lxcstat.GetContainers() {
		containerStat, _ := lxcstat.GetProcStat(containerName)
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, float64(containerStat.User), "cpu", "user", containerName)
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, float64(containerStat.System), "cpu", "system", containerName)
	}

	return nil
}
