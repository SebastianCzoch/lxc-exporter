package kernel

import (
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var (
	kernelVersionFile = "/proc/version"
)

func GetVersion() (string, error) {
	content, err := getVersionFile()
	if err != nil {
		return "", err
	}

	return getKernelVersionFromContent(content), nil
}

func GetMajorVersion() (int, error) {
	content, err := getVersionFile()
	if err != nil {
		return 0, err
	}

	return getMajorVersionFromString(getKernelVersionFromContent(content)), nil
}

func getVersionFile() (string, error) {
	content, err := ioutil.ReadFile(kernelVersionFile)
	return string(content), err
}

func getKernelVersionFromContent(content string) string {
	pattern := regexp.MustCompile("Linux version ([0-9]\\.[0-9]\\.[0-9])+")
	found := pattern.FindString(content)
	return strings.Replace(found, "Linux version ", "", -1)
}

func getMajorVersionFromString(version string) int {
	v, _ := strconv.Atoi(string(version[0]))
	return v
}
