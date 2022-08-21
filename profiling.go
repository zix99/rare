//go:build !rare_no_pprof

package main

import (
	"fmt"
	"os"
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
		})

		oldBefore := app.Before
		app.Before = func(c *cli.Context) error {
			if c.IsSet("profile") {
				basename := c.String("profile")
				startProfiler(basename)
			}

			if oldBefore != nil {
				return oldBefore(c)
			}
			return nil
		}

		oldAfter := app.After
		app.After = func(c *cli.Context) error {
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
