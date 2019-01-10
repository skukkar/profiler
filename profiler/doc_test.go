package profiler_test

import (
	"sahil/profiler"
	"time"
)

func Factorial(n int, sleepTime int64) int {
	fact := 1
	for i := 1; i <= n; i++ {
		fact = fact * i
	}
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	return fact
}

func Example() {
	prof := profiler.NewRootProfiler("testProfiler123")
	defer prof.EndProfile()
	p := prof.StartProfile("fact")
	_ = Factorial(101211111, 2)
	k := p.StartProfile("redis")
	_ = Factorial(1111, 1)
	k.EndProfile()
	m := p.StartProfile("mongo")
	_ = Factorial(1111, 4)
	m.EndProfile()
	p.EndProfile()
	prof.StartProfile("fibonici")
}
