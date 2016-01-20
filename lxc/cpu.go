package lxc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type ProcStat struct {
	User   int
	System int
}

var (
	errorContainerNotFound = errors.New("container not found")
	lxcCPUStatPathPattern  = map[int]string{
		3: "%s/lxc/%s/cpuacct.stat",
		4: "%s/cpu,cpuacct/lxc/%s/cpuacct.stat",
	}
)

func (l *LXC) GetProcStat(containerName string) (ProcStat, error) {
	if !l.containerExists(containerName) {
		return ProcStat{}, errorContainerNotFound
	}

	cpuStat, err := l.fetchProcStat(containerName)
	if err != nil {
		return ProcStat{}, err
	}

	return parseProcStat(cpuStat), nil
}

func (l *LXC) fetchProcStat(containerName string) ([]byte, error) {
	path, err := l.getCPUStatPath(containerName)
	if err != nil {
		return []byte{}, err
	}

	return ioutil.ReadFile(path)
}

func (l *LXC) getCPUStatPath(containerName string) (string, error) {
	if _, ok := lxcCPUStatPathPattern[l.kernelVersion]; !ok {
		return "", errorKernelNotSupported
	}

	return fmt.Sprintf(lxcCPUStatPathPattern[l.kernelVersion], cgroupPath, containerName), nil
}

func parseProcStat(content []byte) ProcStat {
	reg := regexp.MustCompile("\\s\\s+")
	content = reg.ReplaceAll(content, []byte(" "))
	lines := strings.Split(string(content), "\n")
	user := strings.Split(lines[0], " ")
	system := strings.Split(lines[1], " ")

	return ProcStat{User: forceToInt(user[1]), System: forceToInt(system[1])}
}

func forceToInt(variable string) int {
	value, _ := strconv.Atoi(variable)
	return value
}
