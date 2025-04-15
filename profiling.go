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

func init() {
	appModifiers = append(appModifiers, func(app *cli.App) {
		app.Flags = append(app.Flags, &cli.StringFlag{
			Name:  "profile",
			Usage: "Write application profiling information as part of execution. Specify base-name",
		}, &cli.BoolFlag{
			Name:  "metrics",
			Usage: "Outputs runtime memory metrics after a program runs",
		})

		var beforeMem runtime.MemStats
		var start time.Time

		oldBefore := app.Before
		app.Before = func(c *cli.Context) error {
			if c.IsSet("profile") {
				basename := c.String("profile")
				startProfiler(basename)
			}
			if c.Bool("metrics") {
				runtime.ReadMemStats(&beforeMem)
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
				fmt.Printf("Runtime: %s\n", stop.Sub(start).String())
				fmt.Printf("Memory : total=%d; malloc=%d; free=%d; numgc=%d; pausegc=%s\n",
					after.TotalAlloc-beforeMem.TotalAlloc,
					after.Mallocs-beforeMem.Mallocs,
					after.Frees-beforeMem.Frees,
					after.NumGC-beforeMem.NumGC,
					time.Duration(after.PauseTotalNs-beforeMem.PauseTotalNs).String())
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
