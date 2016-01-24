package lxc

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// MemStat is a struct which contains usage of memory read from cgroups
type MemStat struct {
	Usage float64
}

var (
	lxcMemStatPathPattern = map[int]string{
		4: "%s/memory/lxc/%s/memory.usage_in_bytes",
	}
)

// GetMemStat is a method which returns MemStat struct for selected container
func (l *LXC) GetMemStat(containerName string) (MemStat, error) {
	if !l.containerExists(containerName) {
		return MemStat{}, errorContainerNotFound
	}

	memStat, err := l.fetchMemStat(containerName)
	if err != nil {
		return MemStat{}, err
	}

	return parseMemStat(memStat), nil
}

func (l *LXC) fetchMemStat(containerName string) ([]byte, error) {
	path, err := l.getMemStatPath(containerName)
	if err != nil {
		return []byte{}, err
	}

	return ioutil.ReadFile(path)
}

func (l *LXC) getMemStatPath(containerName string) (string, error) {
	if _, ok := lxcMemStatPathPattern[l.kernelVersion]; !ok {
		return "", errorKernelNotSupported
	}

	return fmt.Sprintf(lxcMemStatPathPattern[l.kernelVersion], cgroupPath, containerName), nil
}

func parseMemStat(content []byte) MemStat {
	s := strings.Replace(string(content), "\n", "", -1)
	return MemStat{Usage: forceToFloat64(s)}
}
