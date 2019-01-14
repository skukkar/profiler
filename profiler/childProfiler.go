package profiler

import (
	"fmt"
	"runtime"
	"time"
)

type profiler struct {
	name        string
	startTime   time.Time
	startMemory uint64
	endMemory   uint64
	endTime     time.Time
	profilers   []*profiler
	closed      bool
}

//
//start profile start profile on child level on key and return new profiler
//also attach the same in root profiler
//

func (this *profiler) Start(key string) *profiler {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	newProfiler := new(profiler)
	newProfiler.startTime = time.Now()
	newProfiler.name = this.name + "_" + key
	newProfiler.startMemory = bToKb(m.Alloc)
	this.profilers = append(this.profilers, newProfiler)
	return newProfiler
}

//
//clear all profiler attach wit parent recursive to clear all profilers
//
func (p *profiler) End() {

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	if !p.closed {
		p.endTime = time.Now()
		p.endMemory = bToKb(m.Alloc)
		p.closed = true
	}

	for _, eachProfiler := range p.profilers {
		if !eachProfiler.closed {
			eachProfiler.endTime = time.Now()
			eachProfiler.endMemory = bToKb(m.Alloc)
			eachProfiler.closed = true
		}
		if len(eachProfiler.profilers) > 0 {
			eachProfiler.End()
		}
	}
}

//
// format child profiler data
//
func (this *profiler) formatChildProfile(formattedString []string) []string {
	for _, eachProfiler := range this.profilers {
		formattedString = append(formattedString, fmt.Sprintf("txnName %s timeTakenMS  %d memoryUsageKB %d  \n", eachProfiler.name, eachProfiler.endTime.Sub(eachProfiler.startTime).Nanoseconds()/timeUnit, eachProfiler.endMemory-eachProfiler.startMemory))
		if len(eachProfiler.profilers) > 0 {
			formattedString = append(formattedString, eachProfiler.formatChildProfile(formattedString)...)
		}
	}
	return formattedString
}
