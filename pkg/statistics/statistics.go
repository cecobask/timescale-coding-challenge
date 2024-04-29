package statistics

import (
	"slices"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type Statistics struct {
	Count       int
	TotalTime   time.Duration
	MinimumTime time.Duration
	MaximumTime time.Duration
	AverageTime time.Duration
	MedianTime  time.Duration
	Elements    []time.Duration
}

func NewStatistics(benchmarks []time.Duration) Statistics {
	if len(benchmarks) == 0 {
		return Statistics{}
	}
	slices.Sort(benchmarks)
	count := len(benchmarks)
	minimum := benchmarks[0]
	maximum := benchmarks[count-1]
	total := time.Duration(0)
	for _, b := range benchmarks {
		total += b
	}
	average := total / time.Duration(count)
	median := benchmarks[count/2]
	if count%2 == 0 {
		median = (benchmarks[count/2-1] + benchmarks[count/2]) / 2
	}
	return Statistics{
		Count:       count,
		TotalTime:   total,
		MinimumTime: minimum,
		MaximumTime: maximum,
		AverageTime: average,
		MedianTime:  median,
		Elements:    benchmarks,
	}
}

func (s Statistics) BenchmarksTable(sort bool) string {
	if sort {
		slices.Sort(s.Elements)
	}
	var rows []table.Row
	for _, b := range s.Elements {
		rows = append(rows, table.Row{b})
	}
	benchTable := table.NewWriter()
	benchTable.SetAutoIndex(true)
	benchTable.SetStyle(table.StyleBold)
	benchTable.Style().Options.SeparateRows = true
	benchTable.Style().Title.Align = text.AlignCenter
	benchTable.AppendHeader(table.Row{"Duration"})
	benchTable.AppendRows(rows)
	return benchTable.Render()
}

func (s Statistics) StatisticsTable() string {
	var rows []table.Row
	rows = append(rows,
		table.Row{
			"Count",
			s.Count,
		}, table.Row{
			"Total Time",
			s.TotalTime,
		}, table.Row{
			"Minimum Time",
			s.MinimumTime,
		}, table.Row{
			"Maximum Time",
			s.MaximumTime,
		}, table.Row{
			"Average Time",
			s.AverageTime,
		}, table.Row{
			"Median Time",
			s.MedianTime,
		},
	)
	statsTable := table.NewWriter()
	statsTable.SetStyle(table.StyleBold)
	statsTable.Style().Options.SeparateRows = true
	statsTable.Style().Title.Align = text.AlignCenter
	statsTable.SetTitle("QUERY STATISTICS")
	statsTable.AppendRows(rows)
	return statsTable.Render()
}
