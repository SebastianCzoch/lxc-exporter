package lxc

import (
	"testing"

	"github.com/SebastianCzoch/lxc-exporter/cpu"
	"github.com/stretchr/testify/assert"
)

func TestGetProcStatNotExistedContainer(t *testing.T) {
	l := &LXC{
		kernelVersion:  3,
		containersPath: "../test/sys/fs/cgroup/lxc",
	}

	_, err := l.GetProcStat("not-existed")
	assert.Error(t, err)
}

func TestGetCPUStatPathFail(t *testing.T) {
	l := &LXC{
		kernelVersion:  2,
		containersPath: "../test/sys/fs/cgroup/lxc",
	}

	_, err := l.getCPUStatPath("test")
	assert.Error(t, err)
}

func TestGetCPUStatPath(t *testing.T) {
	l := &LXC{
		kernelVersion:  3,
		containersPath: "../test/sys/fs/cgroup/lxc",
	}

	res, err := l.getCPUStatPath("test")
	assert.NoError(t, err)
	assert.Equal(t, "/sys/fs/cgroup/lxc/test/cpuacct.stat", res)
}

func TestFetchProcStatFail(t *testing.T) {
	l := &LXC{
		kernelVersion:  3,
		containersPath: "../test/sys/fs/cgroup/lxc",
	}

	_, err := l.fetchProcStat("not-existed")
	assert.Error(t, err)
}

func TestFetchProcStat(t *testing.T) {
	l := &LXC{
		kernelVersion:  3,
		containersPath: "../test/sys/fs/cgroup/lxc",
	}

	cgroupPath = "../test/sys/fs/cgroup"
	res, err := l.fetchProcStat("test-container")
	assert.NoError(t, err)
	assert.Equal(t, "user 3092210213\nsystem 311472536\n", string(res))
}

func TestForceToInt(t *testing.T) {
	assert.Equal(t, 1, forceToInt("1"))
	assert.Equal(t, 0, forceToInt("a1"))
	assert.Equal(t, 0, forceToInt("1a"))
}

func TestParseProcStat(t *testing.T) {
	res := parseProcStat([]byte("user 3092210213\nsystem 311472536\n"))
	assert.Equal(t, 3092210213, res.User)
	assert.Equal(t, 311472536, res.System)
}

func TestGetProcStat(t *testing.T) {
	l := &LXC{
		kernelVersion:  3,
		containersPath: "../test/sys/fs/cgroup/lxc",
	}

	cgroupPath = "../test/sys/fs/cgroup"
	res, err := l.GetProcStat("test-container")
	assert.NoError(t, err)
	assert.Equal(t, 3092210213, res.User)
	assert.Equal(t, 311472536, res.System)
}

func TestCalculateUsageInPrecentage(t *testing.T) {
	physical := cpu.ProcStat{
		User:   371569,
		System: 78721,
		Nice:   342711,
		Idle:   39660594,
		Wait:   23304,
		Irq:    0,
		Srq:    19646,
		Zero:   0,
	}

	object := ProcStat{
		User:   2210213,
		System: 79536,
	}

	assert.Equal(t, float32(5.45), object.CalculateUsageInPrecentage(physical))
}
