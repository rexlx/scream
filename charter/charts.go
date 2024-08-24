package charter

import (
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func createLineChart(seriesName string, data []float64) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemePurplePassion}),
	)
	items := make([]opts.LineData, 0)
	xAxis := []string{}
	smoothLine := opts.LineChart{Smooth: opts.Bool(true)}
	for i := 0; i < len(data); i++ {
		xAxis = append(xAxis, strconv.Itoa(i))
		items = append(items, opts.LineData{Value: data[i]})
	}

	line.SetXAxis(xAxis).
		AddSeries(seriesName, items).
		SetSeriesOptions(charts.WithLineChartOpts(smoothLine))
	return line
}
