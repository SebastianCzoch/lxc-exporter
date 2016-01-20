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
	User       int
	System     int
	Nice       int
	Idle       int
	Wait       int
	Irq        int
	Srq        int
	Zero       int
	prevUser   int
	prevSystem int
	prevNice   int
	prevIdle   int
	prevWait   int
	prevIrq    int
	prevSrq    int
	prevZero   int
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
		User:       forceToInt(cpuSum[1]),
		System:     forceToInt(cpuSum[2]),
		Nice:       forceToInt(cpuSum[3]),
		Idle:       forceToInt(cpuSum[4]),
		Wait:       forceToInt(cpuSum[5]),
		Irq:        forceToInt(cpuSum[6]),
		Srq:        forceToInt(cpuSum[7]),
		Zero:       forceToInt(cpuSum[8]),
		prevUser:   forceToInt(cpuSum[1]),
		prevSystem: forceToInt(cpuSum[2]),
		prevNice:   forceToInt(cpuSum[3]),
		prevIdle:   forceToInt(cpuSum[4]),
		prevWait:   forceToInt(cpuSum[5]),
		prevIrq:    forceToInt(cpuSum[6]),
		prevSrq:    forceToInt(cpuSum[7]),
		prevZero:   forceToInt(cpuSum[8]),
	}
}

func forceToInt(variable string) int {
	value, _ := strconv.Atoi(variable)
	return value
}
