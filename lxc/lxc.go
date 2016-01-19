package lxc

import (
	"errors"
	"fmt"
	"os"

	"github.com/SebastianCzoch/lxc-exporter/lxc/lxc0"
)

type Usage struct {
	System int
	User   int
}

type LXC interface {
	GetContainers() []string
	GetCPUUsage(container string) (lxc0.Usage, error)
}

type LXCVersion string

const (
	LXC0       = "lxc0"
	LXC1       = "lxc1"
	cgroupPath = "test/sys/fs/cgroup" //TODO FIXME!

	errorNoCGroupsFound = "no cgroups found"
	errorNoLXCFound     = "no LXC found"
)

func GetLXCManager() (LXC, error) {
	err := checkCGroups()
	if err != nil {
		return nil, err
	}

	if isLXC0() {
		return lxc0.LXC{}, nil
	}

	if isLXC1() {
		return nil, errors.New("not supported yet")
	}

	return nil, errors.New(errorNoLXCFound)
}

func isLXC0() bool {
	_, err := os.Stat(fmt.Sprintf("%s/%s", cgroupPath, "lxc"))
	return !os.IsNotExist(err)
}

func isLXC1() bool {
	_, err := os.Stat(fmt.Sprintf("%s/%s/%s", cgroupPath, "cpu", "lxc"))
	return !os.IsNotExist(err)
}

func checkCGroups() error {
	_, err := os.Stat(cgroupPath)
	if os.IsNotExist(err) {
		return errors.New(errorNoCGroupsFound)
	}

	return nil
}
