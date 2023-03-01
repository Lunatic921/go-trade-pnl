package chart

import (
	"io"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

type Chart interface {
	Draw(w io.Writer) error
}

type LineChart struct {
	Chart
	Line *charts.Line
}

type ScatterChart struct {
	Chart
	Scatter *charts.Scatter
}

func (lc *LineChart) SetChartOptions() {
	lc.Line.SetGlobalOptions(charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeEssos}))
	lc.Line.SetSeriesOptions(
		charts.WithLineChartOpts(opts.LineChart{Smooth: true}),
		charts.WithAreaStyleOpts(opts.AreaStyle{
			Opacity: 0.2,
		}))
	lc.Line.SetGlobalOptions(
		charts.WithVisualMapOpts(opts.VisualMap{
			Left:       "right",
			Min:        0,
			Max:        .0001,
			InRange:    &opts.VisualMapInRange{Color: []string{"red", "green"}},
			Text:       []string{">=0", "<=0"},
			Calculable: true}))
}
