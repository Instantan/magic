package magic_test

import (
	"runtime"
	"testing"
	"time"

	"github.com/Instantan/magic"
)

func BenchmarkSignals(b *testing.B) {

	// for a simple memory benchmark we assume that there are 200_000 clients with each 1000 signals
	// that means we benchmark the memory allocation of 20 million signals

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	getters := []magic.Get[int]{}
	setters := []magic.Set[int]{}

	for i := 0; i < 20_000_000; i++ {
		g, s := magic.Signal(i)
		getters = append(getters, g)
		setters = append(setters, s)
	}

	b.Logf("Alloc = %v MiB", bToMb(m.Alloc))
	b.Logf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	b.Logf("\tSys = %v MiB", bToMb(m.Sys))
	b.Logf("\tNumGC = %v\n", m.NumGC)

	b.Logf("Allocated %v signal getters", len(getters))
	b.Logf("Allocated %v signals setters", len(setters))
	time.Sleep(time.Second * 5)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
