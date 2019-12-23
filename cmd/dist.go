package cmd

import (
	"rare/cmd/helpers"
	"rare/cmd/readProgress"
	"rare/pkg/aggregation"
	"rare/pkg/fuzzy"
	"rare/pkg/multiterm"

	"github.com/urfave/cli"
)

type distCounter struct {
	counter *aggregation.MatchCounter
	lookup  *fuzzy.FuzzyTable
}

func (s *distCounter) Sample(element string) {
	m, _ := s.lookup.GetMatchId(element)
	s.counter.Sample(s.lookup.GetString(m))
}

func (s *distCounter) ParseErrors() uint64 {
	return s.counter.ParseErrors()
}

func distFunction(c *cli.Context) error {
	counter := &distCounter{
		counter: aggregation.NewCounter(),
		lookup:  fuzzy.NewFuzzyTable(0.9),
	}
	writer := multiterm.NewHistogram(multiterm.New(), 10)
	writer.ShowBar = true
	writer.ShowPercentage = true

	ext := helpers.BuildExtractorFromArguments(c)

	helpers.RunAggregationLoop(ext, counter, func() {
		writer.UpdateSamples(counter.counter.Count())
		items := counter.counter.ItemsSorted(10, false)
		for idx, item := range items {
			writer.WriteForLine(idx, item.Name, item.Item.Count())
		}
		writer.InnerWriter().WriteForLine(11, readProgress.GetReadFileString())
	})

	writer.InnerWriter().Close()

	return nil
}

func distCommand() *cli.Command {
	return helpers.AdaptCommandForExtractor(cli.Command{
		Name:   "dist",
		Action: distFunction,
	})
}
