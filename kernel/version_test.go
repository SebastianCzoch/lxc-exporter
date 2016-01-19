package kernel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersionFile(t *testing.T) {
	kernelVersionFile = "worng-path"
	_, err := getVersionFile()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "open worng-path: no such file or directory")
}

func TestGetKernelVersionFromContent(t *testing.T) {
	content := "Linux version 4.1.6-249 (build-bot@example.com) (gcc version 4.9.1 (Ubuntu/Linaro 4.9.1-10ubuntu2) ) #1 SMP Mon Aug 31 12:34:47 UTC 2015"
	result := getKernelVersionFromContent(content)
	assert.Equal(t, "4.1.6", result)
}

func TestGetMajorVersionFromString(t *testing.T) {
	content := "4.1.6-249"
	result := getMajorVersionFromString(content)
	assert.Equal(t, 4, result)
}

func TestGetMajorVersionError(t *testing.T) {
	kernelVersionFile = "not-existing-path"
	result, err := GetMajorVersion()
	assert.Empty(t, result)
	assert.Error(t, err)
}

func TestGetMajorVersionSuccess(t *testing.T) {
	kernelVersionFile = "../test/proc/version"
	result, err := GetMajorVersion()
	assert.Equal(t, result, 4)
	assert.NoError(t, err)
}

func TestGetVersionError(t *testing.T) {
	kernelVersionFile = "../test/proc/version"
	result, err := GetVersion()
	assert.Equal(t, result, "4.1.6")
	assert.NoError(t, err)
}
