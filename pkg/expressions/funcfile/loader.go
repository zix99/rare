package funcfile

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"rare/pkg/expressions"
	"rare/pkg/logger"
	"strings"
)

func LoadDefinitionsFile(compiler *expressions.KeyBuilder, filename string) (map[string]expressions.KeyBuilderFunction, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return LoadDefinitions(compiler, f)
}

func LoadDefinitions(compiler *expressions.KeyBuilder, r io.Reader) (map[string]expressions.KeyBuilderFunction, error) {
	scanner := bufio.NewScanner(r)
	ret := make(map[string]expressions.KeyBuilderFunction)

	errors := 0
	linenum := 0
	for {
		// read possible multiline
		var sb strings.Builder
		for scanner.Scan() {
			linenum++

			// Get line and split out comments
			line := strings.TrimSpace(trimAfter(scanner.Text(), '#'))
			if line == "" {
				continue
			}

			if strings.HasSuffix(line, "\\") { // multiline
				sb.WriteString(strings.TrimSuffix(line, "\\"))
			} else {
				sb.WriteString(line)
				break
			}
		}
		if sb.Len() == 0 {
			break
		}
		phrase := sb.String()

		// Split arguments
		args := strings.SplitN(phrase, " ", 2)
		if len(args) != 2 {
			continue
		}

		// Compile and save
		fnc, err := createAndAddFunc(compiler, args[0], args[1])
		if err != nil {
			logger.Printf("Error creating function '%s', line %d: %s", args[0], linenum, err)
			errors++
		} else {
			ret[args[0]] = fnc
		}
	}

	if errors > 0 {
		return ret, fmt.Errorf("%d compile errors while loading func spec", errors)
	}
	return ret, nil
}

func trimAfter(s string, r rune) string {
	idx := strings.IndexRune(s, r)
	if idx < 0 {
		return s
	}
	return s[:idx]
}

func createAndAddFunc(compiler *expressions.KeyBuilder, name, expr string) (expressions.KeyBuilderFunction, error) {
	kb, err := compiler.Compile(expr)
	if err != nil {
		return nil, err
	}

	fnc := keyBuilderToFunction(kb)
	compiler.Func(name, fnc)
	return fnc, nil
}
