package aggregation

import (
	"errors"
	"rare/pkg/aggregation/sorting"
	"rare/pkg/expressions"
	"rare/pkg/stringSplitter"
	"strings"
)

type GroupKey string

type accumulatorDataDefn struct {
	name    string
	expr    *expressions.CompiledKeyBuilder
	initial string
}

type accumulatorGroupDefn struct {
	name string
	expr *expressions.CompiledKeyBuilder
}

type exprAccumulatorContext struct {
	current   string
	match     string
	keyLookup func(string) string
}

func (s *exprAccumulatorContext) GetMatch(idx int) (ret string) {
	if idx == 0 {
		return s.match
	}

	// Index 1+, parse the string as if it's a range (Without heap alloc)
	splitter := stringSplitter.Splitter{S: s.match, Delim: expressions.ArraySeparatorString}
	for i := 0; i < idx; i++ {
		ret = splitter.Next()
	}
	return
}

func (s *exprAccumulatorContext) GetKey(key string) string {
	if key == "." {
		return s.current
	}
	if s.keyLookup != nil {
		return s.keyLookup(key)
	}
	return ""
}

type AccumulatingGroup struct {
	data map[GroupKey][]string // group -> columns

	compiler     *expressions.KeyBuilder
	groupDef     []*accumulatorGroupDefn
	colDef       []*accumulatorDataDefn // colname -> expr
	colIdxLookup map[string]int         // name to col-index
	sortExpr     *expressions.CompiledKeyBuilder

	context *exprAccumulatorContext
}

func NewAccumulatingGroup(compiler *expressions.KeyBuilder) *AccumulatingGroup {
	return &AccumulatingGroup{
		data:         make(map[GroupKey][]string),
		colIdxLookup: make(map[string]int),
		compiler:     compiler,
		context:      &exprAccumulatorContext{},
	}
}

func (s *AccumulatingGroup) AddGroupExpr(name, expr string) error {
	if len(s.data) > 0 {
		return errors.New("unable to add new group to existing data")
	}
	for _, item := range s.groupDef {
		if item.name == name {
			return errors.New("duplicate group")
		}
	}

	kb, err := s.compiler.Compile(expr)
	if err != nil {
		return err
	}
	s.groupDef = append(s.groupDef, &accumulatorGroupDefn{
		expr: kb,
		name: name,
	})
	return nil
}

func (s *AccumulatingGroup) AddDataExpr(name, expr, initial string) error {
	if len(s.data) > 0 {
		return errors.New("unable to add new expression to existing data")
	}
	if _, ok := s.colIdxLookup[name]; ok {
		return errors.New("duplicate data expression")
	}

	kb, err := s.compiler.Compile(expr)
	if err != nil {
		return err
	}

	s.colDef = append(s.colDef, &accumulatorDataDefn{
		name:    name,
		expr:    kb,
		initial: initial,
	})
	s.colIdxLookup[name] = len(s.colDef) - 1

	return nil
}

func (s *AccumulatingGroup) SetSort(expr string) error {
	compiled, err := s.compiler.Compile(expr)
	if err != nil {
		return err
	}
	s.sortExpr = compiled
	return nil
}

func (s *AccumulatingGroup) Sample(element string) {
	// init shared context
	// this is not multithreaded code! So this saves an alloc
	ctx := s.context
	ctx.match = element
	ctx.current = ""
	ctx.keyLookup = nil

	// Get which group this will belong to
	groupKey := s.buildGroupKey(ctx)

	rowData, hasRow := s.data[groupKey]
	if !hasRow {
		// Create new row & initialize
		rowData = make([]string, len(s.colDef))
		for i, colDef := range s.colDef {
			rowData[i] = colDef.initial
		}
		s.data[groupKey] = rowData
	}

	// Now that row are defined, allow retrieving them
	ctx.keyLookup = func(key string) string {
		if idx, ok := s.colIdxLookup[key]; ok {
			return rowData[idx]
		}
		return ""
	}

	// Sample each data point in row group
	for idx, dataExpr := range s.colDef {
		ctx.current = rowData[idx]
		rowData[idx] = dataExpr.expr.BuildKey(ctx)
	}
}

func (s *AccumulatingGroup) ParseErrors() uint64 {
	return 0
}

func (s *AccumulatingGroup) buildGroupKey(ctx expressions.KeyBuilderContext) GroupKey {
	if len(s.groupDef) == 0 {
		return ""
	}

	if len(s.groupDef) == 1 {
		return GroupKey(s.groupDef[0].expr.BuildKey(ctx))
	}

	var sb strings.Builder
	for i, gexpr := range s.groupDef {
		if i > 0 {
			sb.WriteRune(expressions.ArraySeparator)
		}
		sb.WriteString(gexpr.expr.BuildKey(ctx))
	}
	return GroupKey(sb.String())
}

// Split a group key into its individual parts
func (s GroupKey) Parts() []string {
	if s == "" {
		return make([]string, 0)
	}
	return strings.Split(string(s), expressions.ArraySeparatorString)
}

func (s *AccumulatingGroup) GroupCols() []string {
	ret := make([]string, len(s.groupDef))
	for i, gdef := range s.groupDef {
		ret[i] = gdef.name
	}
	return ret
}

func (s *AccumulatingGroup) DataCols() []string {
	ret := make([]string, len(s.colDef))
	for id, def := range s.colDef {
		ret[id] = def.name
	}
	return ret
}

// context used to evaluate a row's sort-order
type accumulatorGroupSortContext struct {
	groupKey  string
	rowLookup func(string) string
}

func (s *accumulatorGroupSortContext) GetMatch(idx int) (ret string) {
	splitter := stringSplitter.Splitter{
		S:     s.groupKey,
		Delim: expressions.ArraySeparatorString,
	}
	for i := 0; i <= idx; i++ {
		ret = splitter.Next()
	}
	return ret
}

func (s *accumulatorGroupSortContext) GetKey(key string) string {
	if key == "." {
		return s.groupKey
	}
	return s.rowLookup(key)
}

// All possible values that were found for groups (as GroupKey) as a new slice
func (s *AccumulatingGroup) Groups(sort sorting.NameSorter) []GroupKey {
	ret := make([]GroupKey, 0, len(s.data))
	for g := range s.data {
		ret = append(ret, g)
	}
	if s.sortExpr != nil {
		ctx := accumulatorGroupSortContext{}
		sorting.SortBy(ret, sort, func(x GroupKey) string {
			ctx.groupKey = string(x)
			ctx.rowLookup = func(row string) string {
				return s.data[x][s.colIdxLookup[row]]
			}
			return s.sortExpr.BuildKey(&ctx)
		})
	} else {
		sorting.SortBy(ret, sort, func(x GroupKey) string {
			return string(x)
		})
	}
	return ret
}

// Number of defined group columns
func (s *AccumulatingGroup) GroupColCount() int {
	return len(s.groupDef)
}

func (s *AccumulatingGroup) ColCount() int {
	return len(s.groupDef) + len(s.colDef)
}

// Returns a new slice of data-groups for key
func (s *AccumulatingGroup) Data(groupKey GroupKey) []string {
	ret := make([]string, len(s.colDef))
	copy(ret, s.data[groupKey])
	return ret
}

func (s *AccumulatingGroup) DataNoCopy(groupKey GroupKey) []string {
	return s.data[groupKey]
}

func (s *AccumulatingGroup) DataCount() int {
	return len(s.data)
}
