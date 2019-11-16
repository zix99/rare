package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

var profilerDone chan bool

func startProfiler(basename string) {
	fcpu, _ := os.Create(basename + ".cpu.prof")
	pprof.StartCPUProfile(fcpu)
	profilerDone = make(chan bool)

	go func() {
		idx := 0
		for {
			select {
			case <-time.After(500 * time.Millisecond):
				filename := fmt.Sprintf("%s_%03d.prof", basename, idx)
				if f, err := os.Create(filename); err == nil {
					pprof.WriteHeapProfile(f)
				}
				idx++
			case <-profilerDone:
				break
			}
		}
	}()
}

func stopProfile() {
	pprof.StopCPUProfile()
	profilerDone <- true
}
