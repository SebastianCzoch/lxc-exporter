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
	User   int
	System int
	Nice   int
	Idle   int
	Wait   int
	Irq    int
	Srq    int
	Zero   int
}

func GetProcStat() (ProcStat, error) {
	content, err := fetchProcStat()
	if err != nil {
		return ProcStat{}, err
	}

	return parseProcStat(content), nil
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
