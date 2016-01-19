package lxc0

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

const (
	containersPath = "/sys/fs/cgroup/cpu/lxc"
)

type LXC struct {
}

type Usage struct {
	System int
	User   int
}

func (l LXC) GetContainers() []string {
	var containers = []string{}
	files, _ := ioutil.ReadDir(containersPath)
	for _, f := range files {
		if f.IsDir() {
			containers = append(containers, f.Name())
		}
	}

	return containers
}

func (l LXC) GetCPUUsage(container string) (Usage, error) {
	content, err := ioutil.ReadFile(fmt.Sprint(containersPath, "/", container, "/cpuacct.stat"))
	if err != nil {
		return Usage{}, err
	}

	reg := regexp.MustCompile("\\s\\s+")
	content = reg.ReplaceAll(content, []byte(" "))
	lines := strings.Split(string(content), "\n")
	user := strings.Split(lines[0], " ")
	system := strings.Split(lines[1], " ")

	return Usage{User: forceToInt(user[1]), System: forceToInt(system[1])}, nil
}

func forceToInt(variable string) int {
	value, _ := strconv.Atoi(variable)
	return value
}
