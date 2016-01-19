package cpu

import "strconv"

const (
	fixedOne = 1.0
)

type ProcStat struct {
	User   int
	System int
	Nice   int
	Idle   int
	Wait   int
	Irq    int
	Srq    int
	Zero   int
}

type ProcessorsStat struct {
	Sum  ProcStat
	Cpus []ProcStat
}

func (p *ProcStat) CalculateUsage(prevIdle, prevTotal float64) (string, float64, float64) {
	total := float64(p.User + p.System + p.Nice + p.Idle + p.Irq + p.Srq + p.Zero + p.Wait)
	idle := float64(p.Idle + p.Wait)
	diffIdle := idle - prevIdle
	diffTotal := total - prevTotal
	usage := (diffTotal - diffIdle) / diffTotal * 100

	return strconv.FormatFloat(usage, 'f', 2, 32), idle, total
}

func (p *ProcStat) CalculateCGroups(user, system int, prevIdle, prevTotal float64) (string, float64, float64) {
	total := float64(user + system + p.Idle + p.Wait)
	idle := float64(p.Idle + p.Wait)
	diffIdle := idle - prevIdle
	diffTotal := total - prevTotal
	usage := (diffTotal - diffIdle) / diffTotal * 100

	return strconv.FormatFloat(usage, 'f', 2, 32), idle, total
}
