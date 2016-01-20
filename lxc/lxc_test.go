package lxc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckCGroupsWithError(t *testing.T) {
	cgroupPath = "not-exists"
	assert.Error(t, checkCGroups())
}

func TestCheckCGroupsWithSuccess(t *testing.T) {
	cgroupPath = "../test/sys/fs/cgroup"
	assert.NoError(t, checkCGroups())
}

func TestNewWithError(t *testing.T) {
	cgroupPath = "not-exists"
	_, err := New(1)
	assert.Error(t, err)
}

func TestNewErrorKernelVersion(t *testing.T) {
	cgroupPath = "../test/sys/fs/cgroup"
	_, err := New(2)
	assert.Error(t, err)
}

func TestNew(t *testing.T) {
	cgroupPath = "../test/sys/fs/cgroup"
	_, err := New(3)
	assert.NoError(t, err)
}

func TestGetContainersPathWrongKernel(t *testing.T) {
	_, err := getContainersPath(2)
	assert.Error(t, err)
}

func TestGetContainersPath(t *testing.T) {
	res, err := getContainersPath(3)
	assert.NoError(t, err)
	assert.Equal(t, "/sys/fs/cgroup/lxc", res)

	res, err = getContainersPath(4)
	assert.NoError(t, err)
	assert.Equal(t, "/sys/fs/cgroup/cpu,cpuacct/lxc", res)
}

func TestGetContainersEmptyList(t *testing.T) {
	l := &LXC{
		kernelVersion:  3,
		containersPath: "../test/sys/fs/cgroup/lxc_empty",
	}

	assert.Empty(t, l.GetContainers())
}

func TestGetContainersList(t *testing.T) {
	l := &LXC{
		kernelVersion:  3,
		containersPath: "../test/sys/fs/cgroup/lxc",
	}
	list := l.GetContainers()
	assert.Len(t, list, 2)
	assert.Equal(t, list[0], "test-container")
}

func TestContainerExists(t *testing.T) {
	l := &LXC{
		kernelVersion:  3,
		containersPath: "../test/sys/fs/cgroup/lxc",
	}

	assert.False(t, l.containerExists("not-existed"))
	assert.True(t, l.containerExists("test-container"))
}
