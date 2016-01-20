package cpu

import (
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var (
	procStatPath = "/proc/stat"
)

type ProcStat struct {
	User      int
	System    int
	Nice      int
	Idle      int
	Wait      int
	Irq       int
	Srq       int
	Zero      int
	prevTotal float32
	prevIdle  float32
}

func GetProcStat() (ProcStat, error) {
	content, err := fetchProcStat()
	if err != nil {
		return ProcStat{}, err
	}

	return parseProcStat(content), nil
}

func (p *ProcStat) CalculateUsageInPrecentage() float32 {
	total := float32(p.User + p.System + p.Nice + p.Idle + p.Irq + p.Srq + p.Zero + p.Wait)
	idle := float32(p.Idle + p.Wait)
	diffIdle := idle - p.prevIdle
	diffTotal := total - p.prevTotal
	usage := (diffTotal - diffIdle) / diffTotal * 100

	p.prevIdle = idle
	p.prevTotal = total
	return float32(int(usage*100)) / 100
}

func fetchProcStat() ([]byte, error) {
	return ioutil.ReadFile(procStatPath)
}

func parseProcStat(content []byte) ProcStat {
	reg := regexp.MustCompile("\\s\\s+")
	content = reg.ReplaceAll(content, []byte(" "))
	lines := strings.Split(string(content), "\n")
	cpuSum := strings.Split(lines[0], " ")

	return ProcStat{
		User:   forceToInt(cpuSum[1]),
		System: forceToInt(cpuSum[2]),
		Nice:   forceToInt(cpuSum[3]),
		Idle:   forceToInt(cpuSum[4]),
		Wait:   forceToInt(cpuSum[5]),
		Irq:    forceToInt(cpuSum[6]),
		Srq:    forceToInt(cpuSum[7]),
		Zero:   forceToInt(cpuSum[8]),
	}
}

func forceToInt(variable string) int {
	value, _ := strconv.Atoi(variable)
	return value
}
