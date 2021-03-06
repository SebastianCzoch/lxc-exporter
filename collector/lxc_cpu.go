package collector

import (
	"fmt"

	"github.com/SebastianCzoch/lxc-exporter/cpu"
	"github.com/SebastianCzoch/lxc-exporter/kernel"
	"github.com/SebastianCzoch/lxc-exporter/lxc"
	"github.com/prometheus/client_golang/prometheus"
)

type lxcCPUCollector struct {
	cpu                       *prometheus.Desc
	cpuPrecentage             *prometheus.Desc
	cpuRealPhysical           *prometheus.Desc
	cpuRealPhysicalPrecentage *prometheus.Desc
	lxcStat                   *lxc.LXC
}

var (
	containersStat = make(map[string]lxc.ProcStat)
	physicalStat   cpu.ProcStat
)

// NewCPUStatCollector is a function which return new lxc cpu collector
func NewCPUStatCollector() Collector {
	return &lxcCPUCollector{
		cpu: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "cpu"),
			"Seconds the cpus spent in each mode.",
			[]string{"mode", "container"}, nil,
		),
		cpuPrecentage: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "cpu_precentage"),
			"Precentage of usage processor",
			[]string{"container"}, nil,
		),
		cpuRealPhysical: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "cpu_physical_real"),
			"Seconds the real physical cpu spent in each mode.",
			[]string{"mode"}, nil,
		),
		cpuRealPhysicalPrecentage: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "cpu_physical_real_precentage"),
			"Precentage of usage processor",
			[]string{}, nil,
		),
	}
}

func (c *lxcCPUCollector) Init() error {
	kernelVersion, err := kernel.GetMajorVersion()
	if err != nil {
		return err
	}

	c.lxcStat, err = lxc.New(kernelVersion)
	physicalStat, _ = cpu.GetProcStat()
	c.scrapMetrics()
	return err
}

func (c *lxcCPUCollector) Update(ch chan<- prometheus.Metric) error {
	physicalStat, _ = physicalStat.Refresh()
	var totalSystem float64
	var totalUser float64
	for _, containerName := range c.lxcStat.GetContainers() {
		containerStat, err := c.getContainerStat(containerName)
		if err != nil {
			continue
		}
		totalSystem += containerStat.System
		totalUser += containerStat.User

		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, containerStat.User, "user", containerName)
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, containerStat.System, "system", containerName)
		ch <- prometheus.MustNewConstMetric(c.cpuPrecentage, prometheus.CounterValue, c.getPrecentage(&containerStat), containerName)
	}
	ch <- prometheus.MustNewConstMetric(c.cpuRealPhysical, prometheus.CounterValue, physicalStat.System-totalSystem, "system")
	ch <- prometheus.MustNewConstMetric(c.cpuRealPhysical, prometheus.CounterValue, physicalStat.User-totalUser, "user")
	ch <- prometheus.MustNewConstMetric(c.cpuRealPhysical, prometheus.CounterValue, physicalStat.Idle, "idle")
	ch <- prometheus.MustNewConstMetric(c.cpuRealPhysical, prometheus.CounterValue, physicalStat.Wait, "wait")
	ch <- prometheus.MustNewConstMetric(c.cpuRealPhysical, prometheus.CounterValue, physicalStat.Nice, "nice")
	ch <- prometheus.MustNewConstMetric(c.cpuRealPhysical, prometheus.CounterValue, physicalStat.Srq, "sqr")
	ch <- prometheus.MustNewConstMetric(c.cpuRealPhysical, prometheus.CounterValue, physicalStat.Irq, "irq")
	ch <- prometheus.MustNewConstMetric(c.cpuRealPhysicalPrecentage, prometheus.CounterValue, c.getRealPhysicalPrecentage(totalUser, totalSystem, &physicalStat))

	return nil
}

func (c *lxcCPUCollector) getPrecentage(procStat *lxc.ProcStat) float64 {
	idle := float64(physicalStat.Idle + physicalStat.Wait)
	total := float64(procStat.User+procStat.System) + idle
	precentage := (total - idle) / total * 100

	return float64(int(precentage*100)) / 100
}

func (c *lxcCPUCollector) getRealPhysicalPrecentage(totalUser, totalSystem float64, procStat *cpu.ProcStat) float64 {
	idle := procStat.Idle + procStat.Wait
	total := procStat.User + procStat.System - totalSystem - totalUser + idle
	precentage := (total - idle) / total * 100

	precentage = float64(int(precentage*100)) / 100
	if precentage < 0 {
		return 0
	}

	return precentage
}

func (c *lxcCPUCollector) getContainerStat(containerName string) (lxc.ProcStat, error) {
	prevStat, err := c.getPrevStat(containerName)
	if err != nil {
		c.scrapMetric(containerName)
		return prevStat, err
	}

	acctStat, _ := c.lxcStat.GetProcStat(containerName)
	containersStat[containerName] = acctStat
	return lxc.ProcStat{
		User:   acctStat.User - prevStat.User,
		System: acctStat.System - prevStat.System,
	}, nil
}

func (c *lxcCPUCollector) scrapMetrics() {
	for _, containerName := range c.lxcStat.GetContainers() {
		c.scrapMetric(containerName)
	}
}

func (c *lxcCPUCollector) scrapMetric(containerName string) {
	containersStat[containerName], _ = c.lxcStat.GetProcStat(containerName)
}

func (c *lxcCPUCollector) getPrevStat(containerName string) (lxc.ProcStat, error) {
	if _, ok := containersStat[containerName]; !ok {
		return lxc.ProcStat{}, fmt.Errorf("container %s doesn't exists", containerName)
	}

	return containersStat[containerName], nil
}
