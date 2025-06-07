package main

import (
	"fmt"
	"os"

	"github.com/zix99/rare/cmd"
	"github.com/zix99/rare/cmd/helpers"
	"github.com/zix99/rare/pkg/color"
	"github.com/zix99/rare/pkg/expressions/funcfile"
	"github.com/zix99/rare/pkg/expressions/funclib"
	"github.com/zix99/rare/pkg/expressions/stdlib"
	"github.com/zix99/rare/pkg/humanize"
	"github.com/zix99/rare/pkg/logger"
	"github.com/zix99/rare/pkg/matchers/fastregex"
	"github.com/zix99/rare/pkg/multiterm"
	"github.com/zix99/rare/pkg/multiterm/termunicode"

	"github.com/urfave/cli/v2"
)

type appModifier func(app *cli.App)

var appModifiers []appModifier

func buildApp() *cli.App {
	app := cli.NewApp()

	app.Usage = "A fast regex parser, extractor and realtime aggregator"

	app.Version = fmt.Sprintf("%s, %s; regex: %s", version, buildSha, fastregex.Version)

	app.Description = `Aggregate and display information parsed from text files using
	regex and a simple handlebars-like expressions.

	Run "rare docs overview" or go to https://rare.zdyn.net for more information
	
	https://github.com/zix99/rare`

	app.Copyright = `rare  Copyright (C) 2019 Chris LaPointe
	This program comes with ABSOLUTELY NO WARRANTY.
	This is free software, and you are welcome to redistribute it
	under certain conditions`

	app.UseShortOptionHandling = true
	app.Suggest = true

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "nocolor",
			Aliases: []string{"nc"},
			Usage:   "Disables color output",
		},
		&cli.BoolFlag{
			Name:    "noformat",
			Aliases: []string{"nf"},
			Usage:   "Disable number formatting",
		},
		&cli.BoolFlag{
			Name:    "nounicode",
			Aliases: []string{"nu"},
			Usage:   "Disable usage of unicode characters",
		},
		&cli.BoolFlag{
			Name:    "noload",
			Aliases: []string{"nl"},
			Usage:   "Disable external file loading in expressions",
		},
		&cli.StringSliceFlag{
			Name:    "funcs",
			EnvVars: []string{"RARE_FUNC_FILES"},
			Usage:   "Specify filenames to load expressions from",
		},
		&cli.BoolFlag{
			Name:  "color",
			Usage: "Force-enable color output",
		},
		&cli.BoolFlag{
			Name:  "notrim",
			Usage: "By default, rare will trim output text for in-place updates. Setting this flag will disable that",
		},
	}

	// When showing default help, exit with an error code
	app.Action = func(c *cli.Context) error {
		var err error

		args := c.Args()
		if args.Present() {
			err = cli.ShowCommandHelp(c, args.First())
		} else {
			err = cli.ShowAppHelp(c)
		}

		if err != nil {
			return err
		}
		return cli.Exit("", helpers.ExitCodeInvalidUsage)
	}

	app.Commands = cmd.GetSupportedCommands()

	app.Before = cli.BeforeFunc(func(c *cli.Context) error {
		if c.Bool("nocolor") {
			color.Enabled = false
		} else if c.Bool("color") {
			color.Enabled = true
		}
		if c.Bool("noformat") {
			humanize.Enabled = false
		}
		if c.Bool("notrim") {
			multiterm.AutoTrim = false
		}
		if c.Bool("nounicode") {
			termunicode.UnicodeEnabled = false
		}
		if c.Bool("noload") {
			stdlib.DisableLoad = true
		}
		if funcs := c.StringSlice("funcs"); len(funcs) > 0 {
			cmplr := funclib.NewKeyBuilder()
			for _, ff := range funcs {
				funclib.TryAddFunctions(funcfile.LoadDefinitionsFile(cmplr, ff))
			}
		}
		return nil
	})

	app.ExitErrHandler = func(c *cli.Context, err error) {
		// Suppress built-in handler (Which will exit before running any After())
		// Handle exit-codes in main()
		// This also allows for better unit testing...
	}

	// Apply any plugin/modifiers
	for _, modifier := range appModifiers {
		modifier(app)
	}

	return app
}

func cliMain(args ...string) error {
	return buildApp().Run(args)
}

func main() {
	err := cliMain(os.Args...)
	if err != nil {
		if msg := err.Error(); msg != "" {
			logger.Print(msg)
		}
		if v, ok := err.(cli.ExitCoder); ok {
			os.Exit(v.ExitCode())
		}
		os.Exit(helpers.ExitCodeInvalidUsage)
	}
}
