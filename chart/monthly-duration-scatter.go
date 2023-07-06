package chart

import (
	"go-trade-pnl/trading"
	"io"
	"math"
	"sort"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type MonthlyDurationScatterChart struct {
	ScatterChart
	Month     time.Time
	Portfolio *trading.Portfolio
}

func (c *MonthlyDurationScatterChart) Draw(w io.Writer) error {
	var stockDurationAxisData []float64
	var stockDurationSeriesData []opts.ScatterData

	c.Scatter = charts.NewScatter()

	monthsTrades := c.Portfolio.GetTradesByMonth(c.Month)
	sort.Slice(monthsTrades, func(i, j int) bool {
		return monthsTrades[i].GetDuration().Seconds() < monthsTrades[j].GetDuration().Seconds()
	})

	avgProfitLoss := 0.0
	for _, trade := range monthsTrades {
		avgProfitLoss += math.Abs(trade.GetProfit() / float64(trade.TotalShareCount))
	}

	avgProfitLoss /= float64(len(monthsTrades))

	for _, trade := range monthsTrades {
		profit := trade.GetProfit()

		if math.Abs(profit)/float64(trade.TotalShareCount) < avgProfitLoss*10.0 {
			stockDurationAxisData = append(stockDurationAxisData, trade.GetDuration().Seconds())
			stockDurationSeriesData = append(stockDurationSeriesData, opts.ScatterData{Value: profit})
		}

	}

	c.Scatter.SetXAxis(stockDurationAxisData)
	c.Scatter.AddSeries("Profit by Duration", stockDurationSeriesData)

	c.Scatter.Render(w)

	return nil
}

func (c *MonthlyDurationScatterChart) SetTradeMode(includeSwings bool) {
	c.Portfolio.IncludeSwing = includeSwings
}
