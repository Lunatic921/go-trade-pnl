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

type MonthlyPriceShareChart struct {
	ScatterChart
	Month     time.Time
	Portfolio *trading.Portfolio
}

func (c *MonthlyPriceShareChart) Draw(w io.Writer) error {
	var stockPricesAxisData []float64
	var stockShareSeriesData []opts.ScatterData

	c.Scatter = charts.NewScatter()

	monthsTrades := c.Portfolio.FilterTrades(c.Month.Year(), int(c.Month.Month()), -1)
	sort.Slice(monthsTrades, func(i, j int) bool {
		return monthsTrades[i].GetOpeningPriceAvg() < monthsTrades[j].GetOpeningPriceAvg()
	})

	avgProfitLoss := 0.0
	for _, trade := range monthsTrades {
		avgProfitLoss += math.Abs(trade.GetProfit() / float64(trade.TotalShareCount))
	}

	avgProfitLoss /= float64(len(monthsTrades))

	for _, trade := range monthsTrades {
		stockPricesAxisData = append(stockPricesAxisData, math.Round(trade.GetOpeningPriceAvg()))
		stockShareSeriesData = append(stockShareSeriesData, opts.ScatterData{Value: trade.TotalShareCount, Name: trade.Ticker})
	}

	c.Scatter.SetXAxis(stockPricesAxisData)
	c.Scatter.AddSeries("Shares by Price", stockShareSeriesData)

	c.Scatter.Render(w)

	return nil
}

func (c *MonthlyPriceShareChart) SetTradeMode(includeSwings bool) {
	c.Portfolio.IncludeSwing = includeSwings
}
