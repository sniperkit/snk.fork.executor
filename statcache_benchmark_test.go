/*
Sniperkit-Bot
- Status: analyzed
*/

package executor

import (
	"math/rand"
	"runtime"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

var sink *uint64
var statNames [100000000]string

func init() {
	sink = new(uint64)

	for i := 0; i < len(statNames); i++ {
		statNames[i] = strconv.Itoa(i)
	}
}

type benchmarkSource struct{}

func (benchmarkSource) Timer(name string) Timer {
	return func(dur time.Duration) { atomic.AddUint64(sink, 1) }
}

func (benchmarkSource) Counter(name string) Counter {
	return func(delta int) { atomic.AddUint64(sink, 1) }
}

func benchmarkCacheN(b *testing.B, n int, cache statCache) {
	runtime.GC()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			set := cache.get(statNames[rand.Intn(n)])
			set.Success(1)
		}
	})
}

var nameRange = []int{10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000}

func benchmarkCache(b *testing.B, cache func(StatSource) statCache) {
	for _, n := range nameRange {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			benchmarkCacheN(b, n, cache(benchmarkSource{}))
		})
	}
}

func BenchmarkMutexCache(b *testing.B) {
	benchmarkCache(b, newMutexCache)
}

func BenchmarkSyncMapCache(b *testing.B) {
	benchmarkCache(b, newSyncMapCache)
}
