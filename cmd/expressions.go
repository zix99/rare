package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/zix99/rare/pkg/color"
	"github.com/zix99/rare/pkg/expressions"
	"github.com/zix99/rare/pkg/expressions/exprofiler"
	"github.com/zix99/rare/pkg/expressions/funclib"
	"github.com/zix99/rare/pkg/humanize"
	"github.com/zix99/rare/pkg/minijson"

	"github.com/urfave/cli/v2"
)

func expressionFunction(c *cli.Context) error {
	var (
		expString   = c.Args().First()
		noOptimize  = c.Bool("no-optimize")
		data        = c.StringSlice("data")
		keyPairs    = c.StringSlice("key")
		benchmark   = c.Bool("benchmark")
		stats       = c.Bool("stats")
		skipNewline = c.Bool("skip-newline")
		detailed    = stats || benchmark
		listFuncs   = c.Bool("listfuncs")
		raw         = c.Bool("raw")
	)

	if listFuncs {
		fmt.Println(color.Wrap(color.Bold, "Builtin:  "), strings.Join(extractFuncNames(funclib.Builtins), ", "))
		if len(funclib.Additional) > 0 {
			fmt.Println(color.Wrap(color.Bold, "FuncsFile:"), strings.Join(extractFuncNames(funclib.Additional), ", "))
		}
		return nil
	}

	if c.NArg() != 1 {
		return errors.New("expected exactly 1 expression argument. Use - for stdin")
	}

	if expString == "" {
		return errors.New("empty expression")
	}
	if expString == "-" {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			return errors.New("error reading input")
		}
		expString = string(b)
	}

	builder := funclib.NewKeyBuilderEx(!noOptimize)
	compiled, err := builder.Compile(expString)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return errors.New("compile error")
	}

	expCtx := expressions.KeyBuilderContextArray{
		Elements: data,
		Keys:     parseKeyValuesIntoMap(keyPairs...),
	}

	// Emulate special keys
	{
		keys := parseKeyValuesIntoMap(keyPairs...)
		expCtx.Keys["src"] = "<args>"
		expCtx.Keys["line"] = "0"
		expCtx.Keys["."] = buildSpecialKeyJson(nil, keys)
		expCtx.Keys["#"] = buildSpecialKeyJson(data, nil)
		expCtx.Keys[".#"] = buildSpecialKeyJson(data, keys)
		expCtx.Keys["#."] = expCtx.Keys[".#"]
		expCtx.Keys["@"] = expressions.MakeArray(data...)
	}

	// Output results
	result := compiled.BuildKey(&expCtx)

	if !raw {
		result = smartFormatResult(result)
	}

	if detailed {
		fmt.Printf("Expression: %s\n", color.Wrap(color.BrightWhite, expString))
		fmt.Printf("Result:     %s\n", color.Wrap(color.BrightYellow, result))
	} else {
		fmt.Print(result)
		if !skipNewline {
			fmt.Println()
		}
	}

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
		fmt.Print(color.Wrapf(color.BrightBlack, "(%s iterations in %s)", humanize.Hi32(iterations), duration.String()))
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

func buildSpecialKeyJson(matches []string, values map[string]string) string {
	var json minijson.JsonObjectBuilder
	json.Open()
	for i, val := range matches {
		json.WriteString(strconv.Itoa(i), val)
	}
	for k, v := range values {
		json.WriteString(k, v)
	}
	json.Close()
	return json.String()
}

func smartFormatResult(s string) string {
	if strings.ContainsRune(s, expressions.ArraySeparator) {
		// Output array
		var sb strings.Builder
		sb.WriteRune('[')
		for idx, val := range strings.Split(s, expressions.ArraySeparatorString) {
			if idx > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(val)
		}
		sb.WriteRune(']')
		return sb.String()
	}
	return s
}

func extractFuncNames(lib funclib.FunctionSet) []string {
	ret := make([]string, 0, len(lib))
	for name := range lib {
		ret = append(ret, name)
	}
	sort.Strings(ret)
	return ret
}

func expressionCommand() *cli.Command {
	return &cli.Command{
		Name:        "expression",
		Usage:       "Evaluate and benchmark expressions",
		Description: "Given an expression, and optionally some data, test the output and performance of an expression",
		ArgsUsage:   "<expression|->",
		Aliases:     []string{"exp", "expr"},
		Action:      expressionFunction,
		Category:    cmdCatHelp,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "listfuncs",
				Usage: "Lists all available expression functions",
			},
			&cli.BoolFlag{
				Name:    "skip-newline",
				Aliases: []string{"n"},
				Usage:   "Don't add a newline character when printing plain result",
			},
			&cli.BoolFlag{
				Name:    "raw",
				Aliases: []string{"r"},
				Usage:   "Don't format arrays, output raw with null-separators",
			},
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
