// package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"regexp"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/SebastianCzoch/lxc-exporter/cpu"
// 	"github.com/SebastianCzoch/lxc-exporter/lxc"
// )

// var (
// 	usage           string
// 	idle            float64
// 	total           float64
// 	containersTotal = make(map[string]float64)
// 	containersIdle  = make(map[string]float64)
// )

// func main() {
// 	manager, _ := lxc.GetLXCManager()
// 	containers := manager.GetContainers()

// 	for {

		content, _ := ioutil.ReadFile("/proc/stat")
		reg := regexp.MustCompile("\\s\\s+")
		content = reg.ReplaceAll(content, []byte(" "))
		lines := strings.Split(string(content), "\n")
		cpuSum := strings.Split(lines[0], " ")

		cpuStruct := cpu.ProcStat{
			User:   forceToInt(cpuSum[1]),
			System: forceToInt(cpuSum[2]),
			Nice:   forceToInt(cpuSum[3]),
			Idle:   forceToInt(cpuSum[4]),
			Wait:   forceToInt(cpuSum[5]),
			Irq:    forceToInt(cpuSum[6]),
			Srq:    forceToInt(cpuSum[7]),
			Zero:   forceToInt(cpuSum[8]),
		}

// 		usage, idle, total = cpuStruct.CalculateUsage(idle, total)
// 		fmt.Println("Physical usage: ", usage, "%")

// 		for _, c := range containers {
// 			containerUsage, _ := manager.GetCPUUsage(c)
// 			usage, containersIdle[c], containersTotal[c] = cpuStruct.CalculateCGroups(containerUsage.User, containerUsage.System, containersIdle[c], containersTotal[c])
// 			fmt.Println("Container ", c, " usage: ", usage, "%")
// 		}
// 		time.Sleep(3 * time.Second)

// 	}
// }

// func forceToInt(variable string) int {
// 	value, _ := strconv.Atoi(variable)
// 	return value
// }
