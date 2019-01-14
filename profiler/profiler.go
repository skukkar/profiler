package profiler

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

//Niltime used when time is nil or no time captured
const (
	NilTime               = -1
	timeUnit        int64 = int64(time.Millisecond)
	fileName              = "/tmp/profiler"
	PROFILER_SWITCH       = false
)

var rate float64 // rate to be used for monitoring

type rootProfiler struct {
	requestId string
	profilers []*profiler
}

//
//Newprofiler returns a new instance of  root profiler
//
func NewRootProfiler(requestId string) ProfilerInterface {
	if !PROFILER_SWITCH {
		return nil
	}
	rp := new(rootProfiler)
	rp.requestId = requestId
	rp.profilers = make([]*profiler, 0)
	return rp
}

//
//StartProfile starts a profiling using profiler instance p at root level return new profiler
//
func (this *rootProfiler) Start(key string) *profiler {
	if !PROFILER_SWITCH {
		return nil
	}
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	newProfiler := new(profiler)
	newProfiler.startTime = time.Now()
	newProfiler.name = this.requestId + "_" + key
	newProfiler.startMemory = bToKb(m.Alloc)
	this.profilers = append(this.profilers, newProfiler)
	return newProfiler
}

//
//EndProfile ends the profiling using profiler instance p for all attached profiles.
//
func (p *rootProfiler) End() {

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	for _, eachProfiler := range p.profilers {
		if !eachProfiler.closed {
			eachProfiler.endMemory = bToKb(m.Alloc)
			eachProfiler.endTime = time.Now()
			eachProfiler.closed = true
		}
		if len(eachProfiler.profilers) > 0 {
			eachProfiler.End()
		}
	}
	formattedString := p.formatRootProfile()
	write2File(formattedString)
}

//
//format root profiler data
//

func (this *rootProfiler) formatRootProfile() []string {
	formattedString := make([]string, 0)
	for _, eachProfiler := range this.profilers {
		formattedString = append(formattedString, fmt.Sprintf("txnName %s timeTakenMS  %d memoryUsageKB %d  \n", eachProfiler.name, eachProfiler.endTime.Sub(eachProfiler.startTime).Nanoseconds()/timeUnit, eachProfiler.endMemory-eachProfiler.startMemory))
		if len(eachProfiler.profilers) > 0 {
			formattedString = append(formattedString, eachProfiler.formatChildProfile(formattedString)...)
		}
	}
	return formattedString
}

//
// write data into file
//

func write2File(formattedString []string) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		//handle error
	}
	defer f.Close()

	if _, err = f.WriteString(strings.Join(formattedString, "\n")); err != nil {
		//handle error
	}
}

//
//convert bytes to kb
//
func bToKb(b uint64) uint64 {
	return b / 1024
}
