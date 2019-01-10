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
	NilTime        = -1
	timeUnit int64 = int64(time.Millisecond)
	fileName       = "/tmp/profiler"
)

var rate float64 // rate to be used for monitoring

type rootProfiler struct {
	requestId string
	profilers []*profiler
}

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
//Newprofiler returns a new instance of  root profiler
//
func NewRootProfiler(requestId string) *rootProfiler {
	// if !config.GlobalAppConfig.Profiler.Enable {
	// 	return nil
	// }
	rp := new(rootProfiler)
	rp.requestId = requestId
	rp.profilers = make([]*profiler, 0)
	return rp
}

//
//StartProfile starts a profiling using profiler instance p at root level return new profiler
//
func (this *rootProfiler) StartProfile(key string) *profiler {
	// if !config.GlobalAppConfig.Profiler.Enable {
	// 	return
	// }
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
//start profile start profile on child level on key and return new profiler
//also attach the same in root profiler
//

func (this *profiler) StartProfile(key string) *profiler {
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
//EndProfile ends the profiling using profiler instance p for all attached profiles.
//
func (p *rootProfiler) EndProfile() {
	// if !config.GlobalAppConfig.Profiler.Enable {
	// 	return NilTime
	// }

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	for _, eachProfiler := range p.profilers {
		if !eachProfiler.closed {
			eachProfiler.endMemory = bToKb(m.Alloc)
			eachProfiler.endTime = time.Now()
			eachProfiler.closed = true
		}
		if len(eachProfiler.profilers) > 0 {
			eachProfiler.EndProfile()
		}
	}
	formattedString := p.formatRootProfile()
	write2File(formattedString)
}

//
//clear all profiler attach wit parent recursive to clear all profilers
//
func (p *profiler) EndProfile() {

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
			eachProfiler.EndProfile()
		}
	}
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
