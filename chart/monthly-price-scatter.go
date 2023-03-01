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

type MonthlyPriceProfitChart struct {
	ScatterChart
	Month     time.Time
	Portfolio *trading.Portfolio
}

func (c *MonthlyPriceProfitChart) Draw(w io.Writer) error {
	var stockPricesAxisData []float64
	var stockProfitSeriesData []opts.ScatterData

	c.Scatter = charts.NewScatter()

	monthsTrades := c.Portfolio.GetTradesByMonth(c.Month)
	sort.Slice(monthsTrades, func(i, j int) bool {
		return monthsTrades[i].GetOpeningPriceAvg() < monthsTrades[j].GetOpeningPriceAvg()
	})

	avgProfitLoss := 0.0
	for _, trade := range monthsTrades {
		avgProfitLoss += math.Abs(trade.GetProfit() / float64(trade.TotalShareCount))
	}

	avgProfitLoss /= float64(len(monthsTrades))

	for _, trade := range monthsTrades {
		profit := trade.GetProfit()

		if math.Abs(profit)/float64(trade.TotalShareCount) < avgProfitLoss*10.0 {
			stockPricesAxisData = append(stockPricesAxisData, math.Round(trade.GetOpeningPriceAvg()))
			stockProfitSeriesData = append(stockProfitSeriesData, opts.ScatterData{Value: profit / float64(trade.TotalShareCount)})
		}

	}

	c.Scatter.SetXAxis(stockPricesAxisData)
	c.Scatter.AddSeries("Profit by Price", stockProfitSeriesData)

	c.Scatter.Render(w)

	return nil
}
