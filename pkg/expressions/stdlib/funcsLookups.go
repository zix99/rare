package stdlib

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/zix99/rare/pkg/expressions"
)

var DisableLoad = false

// {load "filename"}
// loads static file as string
func kfLoadFile(args []expressions.KeyBuilderStage) (expressions.KeyBuilderStage, error) {
	if DisableLoad {
		return stageErrorf(ErrFile, "loading disabled")
	}

	if len(args) != 1 {
		return stageErrArgCount(args, 1)
	}

	filename, ok := expressions.EvalStaticStage(args[0])
	if !ok {
		return stageError(ErrConst)
	}

	f, err := os.Open(filename)
	if err != nil {
		return stageErrorf(ErrFile, "Unable to open file: "+filename)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return stageErrorf(ErrFile, "Error reading file: "+filename)
	}

	sContent := string(content)

	return func(context expressions.KeyBuilderContext) string {
		return sContent
	}, nil
}

func buildLookupTable(content string, commentPrefix string) map[string]string {
	lookup := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := scanner.Text()

		if commentPrefix != "" && strings.HasPrefix(line, commentPrefix) {
			continue
		}

		parts := strings.Fields(line)
		switch len(parts) {
		case 0: //noop
		case 1:
			lookup[parts[0]] = ""
		case 2:
			lookup[parts[0]] = parts[1]
		}
	}
	return lookup
}

// {lookup key "table" [commentPrefix]}
func kfLookupKey(args []expressions.KeyBuilderStage) (expressions.KeyBuilderStage, error) {
	if !isArgCountBetween(args, 2, 3) {
		return stageErrArgRange(args, "2-3")
	}

	content, ok := expressions.EvalStaticStage(args[1])
	if !ok {
		return stageArgError(ErrConst, 1)
	}

	commentPrefix := EvalStageIndexOrDefault(args, 2, "")

	lookup := buildLookupTable(content, commentPrefix)

	return func(context expressions.KeyBuilderContext) string {
		key := args[0](context)
		return lookup[key]
	}, nil
}

// {haskey key "table" [commentprefix]}
func kfHasKey(args []expressions.KeyBuilderStage) (expressions.KeyBuilderStage, error) {
	if !isArgCountBetween(args, 2, 3) {
		return stageErrArgRange(args, "2-3")
	}

	content, ok := expressions.EvalStaticStage(args[1])
	if !ok {
		return stageArgError(ErrConst, 1)
	}

	commentPrefix := EvalStageIndexOrDefault(args, 2, "")

	lookup := buildLookupTable(content, commentPrefix)

	return func(context expressions.KeyBuilderContext) string {
		key := args[0](context)
		_, has := lookup[key]
		return expressions.TruthyStr(has)
	}, nil
}
