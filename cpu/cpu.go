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
	User       float64
	System     float64
	Nice       float64
	Idle       float64
	Wait       float64
	Irq        float64
	Srq        float64
	Zero       float64
	prevUser   float64
	prevSystem float64
	prevNice   float64
	prevIdle   float64
	prevWait   float64
	prevIrq    float64
	prevSrq    float64
	prevZero   float64
}

func GetProcStat() (ProcStat, error) {
	content, err := fetchProcStat()
	if err != nil {
		return ProcStat{}, err
	}

	return parseProcStat(content), nil
}

func (p *ProcStat) Refresh() (ProcStat, error) {
	content, err := fetchProcStat()
	if err != nil {
		return ProcStat{}, err
	}

	newProc := parseProcStat(content)
	return ProcStat{
		User:       newProc.User - p.prevUser,
		System:     newProc.System - p.prevSystem,
		Nice:       newProc.Nice - p.prevNice,
		Idle:       newProc.Idle - p.prevIdle,
		Wait:       newProc.Wait - p.prevWait,
		Irq:        newProc.Irq - p.prevIrq,
		Srq:        newProc.Srq - p.prevSrq,
		Zero:       newProc.Zero - p.prevZero,
		prevUser:   newProc.User,
		prevSystem: newProc.System,
		prevNice:   newProc.Nice,
		prevIdle:   newProc.Idle,
		prevWait:   newProc.Wait,
		prevIrq:    newProc.Irq,
		prevSrq:    newProc.Srq,
		prevZero:   newProc.Zero,
	}, nil
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
		User:       forceToFloat64(cpuSum[1]),
		System:     forceToFloat64(cpuSum[2]),
		Nice:       forceToFloat64(cpuSum[3]),
		Idle:       forceToFloat64(cpuSum[4]),
		Wait:       forceToFloat64(cpuSum[5]),
		Irq:        forceToFloat64(cpuSum[6]),
		Srq:        forceToFloat64(cpuSum[7]),
		Zero:       forceToFloat64(cpuSum[8]),
		prevUser:   forceToFloat64(cpuSum[1]),
		prevSystem: forceToFloat64(cpuSum[2]),
		prevNice:   forceToFloat64(cpuSum[3]),
		prevIdle:   forceToFloat64(cpuSum[4]),
		prevWait:   forceToFloat64(cpuSum[5]),
		prevIrq:    forceToFloat64(cpuSum[6]),
		prevSrq:    forceToFloat64(cpuSum[7]),
		prevZero:   forceToFloat64(cpuSum[8]),
	}
}

func forceToFloat64(variable string) float64 {
	value, _ := strconv.ParseFloat(variable, 64)
	return value
}
