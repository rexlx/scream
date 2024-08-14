package main

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func createLineChart(data []float64) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemePurplePassion}),
	)
	items := make([]opts.LineData, 0)
	smoothLine := opts.LineChart{Smooth: opts.Bool(true)}
	for i := 0; i < len(data); i++ {
		items = append(items, opts.LineData{Value: data[i]})
	}
	line.SetXAxis(nil).
		AddSeries("Category A", items).
		SetSeriesOptions(charts.WithLineChartOpts(smoothLine))
	return line
}
