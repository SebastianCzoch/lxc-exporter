package lxc

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

// ProcStat is a method which represent cpu stats from cgroups
type ProcStat struct {
	User   float64
	System float64
}

var (
	lxcCPUStatPathPattern = map[int]string{
		3: "%s/lxc/%s/cpuacct.stat",
		4: "%s/cpu,cpuacct/lxc/%s/cpuacct.stat",
	}
)

// GetProcStat is a method which returns ProcStat struct for selected container
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

	return ProcStat{User: forceToFloat64(user[1]), System: forceToFloat64(system[1])}
}
