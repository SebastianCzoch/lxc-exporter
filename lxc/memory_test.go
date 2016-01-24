package lxc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMemStatNotExistedContainer(t *testing.T) {
	l := &LXC{
		kernelVersion:  4,
		containersPath: "../test/sys/fs/cgroup/lxc",
	}

	_, err := l.GetMemStat("not-existed")
	assert.Error(t, err)
}

func TestGetMemStatPathFail(t *testing.T) {
	l := &LXC{
		kernelVersion:  2,
		containersPath: "../test/sys/fs/cgroup/lxc",
	}

	_, err := l.getMemStatPath("test")
	assert.Error(t, err)
}

func TestGetMemStatPath(t *testing.T) {
	cgroupPath = "/sys/fs/cgroup"
	l := &LXC{
		kernelVersion: 4,
	}

	res, err := l.getMemStatPath("test")
	assert.NoError(t, err)
	assert.Equal(t, "/sys/fs/cgroup/memory/lxc/test/memory.usage_in_bytes", res)
}

func TestFetchMemStatFail(t *testing.T) {
	l := &LXC{
		kernelVersion:  4,
		containersPath: "../test/sys/fs/cgroup/lxc",
	}

	_, err := l.fetchMemStat("not-existed")
	assert.Error(t, err)
}

func TestFetchMemStat(t *testing.T) {
	l := &LXC{
		kernelVersion: 4,
	}

	cgroupPath = "../test/sys/fs/cgroup"
	res, err := l.fetchMemStat("test-container")
	assert.NoError(t, err)
	assert.Equal(t, "3242834328\n", string(res))
}

func TestParseMemStat(t *testing.T) {
	res := parseMemStat([]byte("3092210213\n"))
	assert.Equal(t, float64(3092210213), res.Usage)
}

func TestGetMemStat(t *testing.T) {
	l := &LXC{
		kernelVersion:  4,
		containersPath: "../test/sys/fs/cgroup/cpu,cpuacct/lxc",
	}

	cgroupPath = "../test/sys/fs/cgroup"
	res, err := l.GetMemStat("test-container")
	assert.NoError(t, err)
	assert.Equal(t, float64(3242834328), res.Usage)
}
