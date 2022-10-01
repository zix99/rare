package cmd

import (
	"errors"
	"fmt"
	"rare/pkg/color"
	"rare/pkg/expressions"
	"rare/pkg/expressions/exprofiler"
	"rare/pkg/expressions/stdlib"
	"rare/pkg/humanize"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

func expressionFunction(c *cli.Context) error {
	var (
		expString  = c.Args().First()
		noOptimize = c.Bool("no-optimize")
		data       = c.StringSlice("data")
		keys       = c.StringSlice("key")
		benchmark  = c.Bool("benchmark")
		stats      = c.Bool("stats")
	)

	if c.NArg() != 1 {
		return errors.New("expected exactly 1 expression argument")
	}

	if expString == "" {
		return errors.New("empty expression")
	}

	builder := stdlib.NewStdKeyBuilderEx(!noOptimize)
	compiled, err := builder.Compile(expString)
	expCtx := expressions.KeyBuilderContextArray{
		Elements: data,
		Keys:     parseKeyValuesIntoMap(keys...),
	}

	if err != nil {
		return err
	}

	fmt.Printf("Expression: %s\n", color.Wrap(color.BrightWhite, expString))
	result := compiled.BuildKey(&expCtx)
	fmt.Printf("Result:     %s\n", color.Wrap(color.BrightYellow, result))

	if stats {
		stats := exprofiler.GetMetrics(compiled, &expCtx)

		fmt.Println()
		fmt.Println("Stats")
		fmt.Printf("  Stages:        %s\n", color.Wrapi(color.BrightWhite, compiled.StageCount()))
		fmt.Printf("  Match Lookups: %s\n", color.Wrapi(color.BrightWhite, stats.MatchLookups))
		fmt.Printf("  Key   Lookups: %s\n", color.Wrapi(color.BrightWhite, stats.KeyLookups))
	}

	if benchmark {
		fmt.Println()
		duration, iterations := exprofiler.Benchmark(compiled, &expCtx)
		perf := (duration / time.Duration(iterations)).String()
		fmt.Printf("Benchmark: %s ", color.Wrap(color.BrightWhite, perf))
		fmt.Print(color.Wrapf(color.BrightBlack, "(%s iterations in %s)", humanize.Hi(iterations), duration.String()))
		fmt.Print("\n")
	}

	return nil
}

// Parse multiple kv's into a map
func parseKeyValuesIntoMap(kvs ...string) map[string]string {
	ret := make(map[string]string)
	for _, item := range kvs {
		k, v := parseKeyValue(item)
		ret[k] = v
	}
	return ret
}

// parse keys and values separated by '='
func parseKeyValue(s string) (string, string) {
	idx := strings.IndexByte(s, '=')
	if idx < 0 {
		return s, s
	}
	return s[:idx], s[idx+1:]
}

func expressionCommand() *cli.Command {
	return &cli.Command{
		Name:        "expression",
		Usage:       "Test and benchmark expressions",
		Description: "Given an expression, and optionally some data, test the output and performance of an expression",
		ArgsUsage:   "<expression>",
		Aliases:     []string{"exp"},
		Action:      expressionFunction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "benchmark",
				Aliases: []string{"b"},
				Usage:   "Benchmark the expression (slow)",
			},
			&cli.BoolFlag{
				Name:    "stats",
				Aliases: []string{"s"},
				Usage:   "Display stats about the expression",
			},
			&cli.StringSliceFlag{
				Name:    "data",
				Aliases: []string{"d"},
				Usage:   "Specify positional data in the expression",
			},
			&cli.StringSliceFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Usage:   "Specify a named argument, a=b",
			},
			&cli.BoolFlag{
				Name:  "no-optimize",
				Usage: "Disable expression static analysis optimization",
			},
		},
	}
}
