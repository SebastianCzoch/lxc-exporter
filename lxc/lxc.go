package lxc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	containerDirectoryPattern = "%s/%s"
)

var (
	cgroupPath = "/sys/fs/cgroup"
	lxcPath    = map[int]string{
		3: fmt.Sprintf("%s/lxc", cgroupPath),
		4: fmt.Sprintf("%s/cpu,cpuacct/lxc", cgroupPath),
	}

	errorNoCGroupsFound     = fmt.Errorf("no cgroups found at %s", cgroupPath)
	errorKernelNotSupported = errors.New("yours version of kernel is not supported")
)

type LXC struct {
	kernelVersion  int
	containersPath string
}

func New(kernelVersion int) (*LXC, error) {
	err := checkCGroups()
	if err != nil {
		return nil, err
	}

	containersPath, err := getContainersPath(kernelVersion)
	if err != nil {
		return nil, err
	}

	return &LXC{
		kernelVersion:  kernelVersion,
		containersPath: containersPath,
	}, nil
}

func (l *LXC) GetContainers() []string {
	var containers = []string{}
	files, _ := ioutil.ReadDir(l.containersPath)
	for _, f := range files {
		if f.IsDir() {
			containers = append(containers, f.Name())
		}
	}

	return containers
}

func (l *LXC) containerExists(containerName string) bool {
	_, err := os.Stat(fmt.Sprintf(containerDirectoryPattern, l.containersPath, containerName))
	return err == nil
}

func getContainersPath(kernelVersion int) (string, error) {
	if _, ok := lxcPath[kernelVersion]; !ok {
		return "", errorKernelNotSupported
	}

	return lxcPath[kernelVersion], nil
}

func checkCGroups() error {
	_, err := os.Stat(cgroupPath)
	if os.IsNotExist(err) {
		return errorNoCGroupsFound
	}

	return nil
}
