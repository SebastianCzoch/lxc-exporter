package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchProcStatFail(t *testing.T) {
	procStatPath = "/not-existed-directory"
	_, err := fetchProcStat()
	assert.Error(t, err)
}

func TestFetchProcStat(t *testing.T) {
	procStatPath = "../test/proc/stat"
	res, err := fetchProcStat()
	assert.NoError(t, err)
	assert.Equal(t, "cpu  371569 70421 342711 39660594 23304 0 19646 0 0 0\ncpu0 94768 15317 87632 9890650 14025 0 19607 0 0 0\ncpu1 92180 6489 86278 9933852 3233 0 11 0 0 0\ncpu2 93749 19132 85917 9920595 2771 0 12 0 0 0\ncpu3 90872 29483 82884 9915497 3275 0 16 0 0 0\nintr 50722203 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 45332079 0 0 534 0 0 0 0 0 2347113 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0\nctxt 39220113\nbtime 1452954785\nprocesses 904649\nprocs_running 2\nprocs_blocked 0\nsoftirq 108313225 0 40488244 4 2450241 0 0 89676 33072501 606 32211953\n", string(res))
}

func TestForceToInt(t *testing.T) {
	assert.Equal(t, 1, forceToInt("1"))
	assert.Equal(t, 0, forceToInt("a1"))
	assert.Equal(t, 0, forceToInt("1a"))
}

func TestParseProcStat(t *testing.T) {
	procStatContent := "cpu  371569 70421 342711 39660594 23304 0 19646 0 0 0\ncpu0 94768 15317 87632 9890650 14025 0 19607 0 0 0\ncpu1 92180 6489 86278 9933852 3233 0 11 0 0 0\ncpu2 93749 19132 85917 9920595 2771 0 12 0 0 0\ncpu3 90872 29483 82884 9915497 3275 0 16 0 0 0\nintr 50722203 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 45332079 0 0 534 0 0 0 0 0 2347113 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0\nctxt 39220113\nbtime 1452954785\nprocesses 904649\nprocs_running 2\nprocs_blocked 0\nsoftirq 108313225 0 40488244 4 2450241 0 0 89676 33072501 606 32211953\n"
	object := parseProcStat([]byte(procStatContent))
	assert.Equal(t, ProcStat{
		User:   371569,
		System: 70421,
		Nice:   342711,
		Idle:   39660594,
		Wait:   23304,
		Irq:    0,
		Srq:    19646,
		Zero:   0,
	}, object)
}

func TestGetProcStatFail(t *testing.T) {
	procStatPath = "/not-existed-directory"
	_, err := GetProcStat()
	assert.Error(t, err)
}

func TestGetProcStat(t *testing.T) {
	procStatPath = "../test/proc/stat"
	object, err := GetProcStat()
	assert.NoError(t, err)
	assert.Equal(t, ProcStat{
		User:   371569,
		System: 70421,
		Nice:   342711,
		Idle:   39660594,
		Wait:   23304,
		Irq:    0,
		Srq:    19646,
		Zero:   0,
	}, object)
}

func TestCalculateUsageInPrecentage(t *testing.T) {
	object := ProcStat{
		User:   371569,
		System: 78721,
		Nice:   342711,
		Idle:   39660594,
		Wait:   23304,
		Irq:    0,
		Srq:    19646,
		Zero:   0,
	}

	assert.Equal(t, float32(2), object.CalculateUsageInPrecentage())
}

func TestCalculateUsageInPrecentageSecond(t *testing.T) {
	object := ProcStat{
		User:      371569,
		System:    67121,
		Nice:      342711,
		Idle:      39660594,
		Wait:      23304,
		Irq:       0,
		Srq:       19646,
		Zero:      0,
		prevIdle:  39650594,
		prevTotal: 39690594,
	}

	assert.Equal(t, float32(95.80), object.CalculateUsageInPrecentage())
}
