//go:build !rare_no_pprof

package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/urfave/cli/v2"
)

var profilerDone chan bool

func startProfiler(basename string) {
	fcpu, _ := os.Create(basename + ".cpu.prof")
	pprof.StartCPUProfile(fcpu)
	profilerDone = make(chan bool)

	go func() {
		idx := 0
	OUTER_LOOP:
		for {
			select {
			case <-time.After(500 * time.Millisecond):
				filename := fmt.Sprintf("%s_%03d.prof", basename, idx)
				if f, err := os.Create(filename); err == nil {
					pprof.WriteHeapProfile(f)
				}
				idx++
			case <-profilerDone:
				break OUTER_LOOP
			}
		}
	}()
}

func stopProfile() {
	pprof.StopCPUProfile()
	profilerDone <- true
}

func startMemoryRecorder() (stop func()) {
	var samples int64
	var peakHeap, peakAlloc, peakStack uint64
	var avgHeap, avgAlloc, avgStack int64

	doneSignal := make(chan struct{})

	go func() {
		var stats runtime.MemStats
		for {
			select {
			case <-time.After(100 * time.Millisecond):
				runtime.ReadMemStats(&stats)

				samples++
				peakHeap = max(peakHeap, stats.HeapInuse)
				avgHeap += (int64(stats.HeapInuse) - avgHeap) / samples

				peakAlloc = max(peakAlloc, stats.HeapAlloc)
				avgAlloc = max(int64(stats.HeapAlloc)-avgAlloc) / samples

				peakStack = max(peakStack, stats.StackInuse)
				avgStack += (int64(stats.StackInuse) - avgStack) / samples
			case <-doneSignal:
				return
			}
		}
	}()

	return func() {
		doneSignal <- struct{}{}
		fmt.Fprintf(os.Stderr, "MemStat: samples=%d, peakHeap=%d, avgHeap=%d, peakAlloc=%d, avgAlloc=%d, peakStack=%d, avgStack=%d\n", samples, peakHeap, avgHeap, peakAlloc, avgAlloc, peakStack, avgStack)
	}
}

func init() {
	appModifiers = append(appModifiers, func(app *cli.App) {
		app.Flags = append(app.Flags, &cli.StringFlag{
			Name:  "profile",
			Usage: "Write application profiling information as part of execution. Specify base-name",
		}, &cli.BoolFlag{
			Name:  "metrics",
			Usage: "Outputs runtime memory metrics after a program runs",
		}, &cli.BoolFlag{
			Name:  "metrics-memory",
			Usage: "Records memory metrics every 100ms to get peaks/averages",
		})

		var startMem runtime.MemStats
		var start time.Time
		var memRecordStop func()

		oldBefore := app.Before
		app.Before = func(c *cli.Context) error {
			if c.IsSet("profile") {
				basename := c.String("profile")
				startProfiler(basename)
			}
			if c.Bool("metrics") {
				runtime.ReadMemStats(&startMem)
			}
			if c.Bool("metrics-memory") {
				memRecordStop = startMemoryRecorder()
			}

			start = time.Now()

			if oldBefore != nil {
				return oldBefore(c)
			}
			return nil
		}

		oldAfter := app.After
		app.After = func(c *cli.Context) error {
			stop := time.Now()

			if c.Bool("metrics") {
				var after runtime.MemStats
				runtime.ReadMemStats(&after)
				fmt.Fprintf(os.Stderr, "Runtime: %s\n", stop.Sub(start).String())
				fmt.Fprintf(os.Stderr, "Memory : total=%d; malloc=%d; free=%d; numgc=%d; pausegc=%s\n",
					after.TotalAlloc-startMem.TotalAlloc,
					after.Mallocs-startMem.Mallocs,
					after.Frees-startMem.Frees,
					after.NumGC-startMem.NumGC,
					time.Duration(after.PauseTotalNs-startMem.PauseTotalNs).String())
			}
			if memRecordStop != nil {
				memRecordStop()
			}
			if c.IsSet("profile") {
				stopProfile()
			}

			if oldAfter != nil {
				return oldAfter(c)
			}
			return nil
		}
	})
}
