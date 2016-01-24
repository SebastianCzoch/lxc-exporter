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

func TestForceToFloat64(t *testing.T) {
	assert.Equal(t, float64(1), forceToFloat64("1"))
	assert.Equal(t, float64(0), forceToFloat64("a1"))
	assert.Equal(t, float64(0), forceToFloat64("1a"))
	assert.Equal(t, float64(1.21), forceToFloat64("1.21"))
	assert.Equal(t, float64(12), forceToFloat64("1.2e1"))
	assert.Equal(t, float64(0), forceToFloat64("1.2s1"))
}

func TestParseProcStat(t *testing.T) {
	procStatContent := "cpu  371569 70421 342711 39660594 23304 0 19646 0 0 0\ncpu0 94768 15317 87632 9890650 14025 0 19607 0 0 0\ncpu1 92180 6489 86278 9933852 3233 0 11 0 0 0\ncpu2 93749 19132 85917 9920595 2771 0 12 0 0 0\ncpu3 90872 29483 82884 9915497 3275 0 16 0 0 0\nintr 50722203 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 45332079 0 0 534 0 0 0 0 0 2347113 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0\nctxt 39220113\nbtime 1452954785\nprocesses 904649\nprocs_running 2\nprocs_blocked 0\nsoftirq 108313225 0 40488244 4 2450241 0 0 89676 33072501 606 32211953\n"
	object := parseProcStat([]byte(procStatContent))
	assert.Equal(t, ProcStat{
		User:       float64(371569),
		System:     float64(70421),
		Nice:       float64(342711),
		Idle:       float64(39660594),
		Wait:       float64(23304),
		Irq:        float64(0),
		Srq:        float64(19646),
		Zero:       float64(0),
		prevUser:   float64(371569),
		prevSystem: float64(70421),
		prevNice:   float64(342711),
		prevIdle:   float64(39660594),
		prevWait:   float64(23304),
		prevIrq:    float64(0),
		prevSrq:    float64(19646),
		prevZero:   float64(0),
	}, object)
}

func TestRefreshError(t *testing.T) {
	procStatPath = "../test/proc/stat"
	object, err := GetProcStat()
	assert.NoError(t, err)
	procStatPath = "/not-existed-directory"
	_, err = object.Refresh()
	assert.Error(t, err)
}

func TestRefresh(t *testing.T) {
	procStatPath = "../test/proc/stat"
	object, err := GetProcStat()
	assert.NoError(t, err)
	object, err = object.Refresh()
	assert.NoError(t, err)
	assert.Equal(t, ProcStat{
		User:       float64(0),
		System:     float64(0),
		Nice:       float64(0),
		Idle:       float64(0),
		Wait:       float64(0),
		Irq:        float64(0),
		Srq:        float64(0),
		Zero:       float64(0),
		prevUser:   float64(371569),
		prevSystem: float64(70421),
		prevNice:   float64(342711),
		prevIdle:   float64(39660594),
		prevWait:   float64(23304),
		prevIrq:    float64(0),
		prevSrq:    float64(19646),
		prevZero:   float64(0),
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
		User:       float64(371569),
		System:     float64(70421),
		Nice:       float64(342711),
		Idle:       float64(39660594),
		Wait:       float64(23304),
		Irq:        float64(0),
		Srq:        float64(19646),
		Zero:       float64(0),
		prevUser:   float64(371569),
		prevSystem: float64(70421),
		prevNice:   float64(342711),
		prevIdle:   float64(39660594),
		prevWait:   float64(23304),
		prevIrq:    float64(0),
		prevSrq:    float64(19646),
		prevZero:   float64(0),
	}, object)
}
