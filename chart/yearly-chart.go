package chart

import (
	"go-trade-pnl/trading"
	"io"
	"sort"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type YearlyChart struct {
	LineChart
	Year      time.Time
	Portfolio *trading.Portfolio
}

func (c *YearlyChart) Draw(w io.Writer) error {
	_, weekOfYear := time.Now().ISOWeek()
	profitsByWeek := make([]float64, weekOfYear)

	tradesThisYear := c.Portfolio.GetTradesByYear(c.Year)

	sort.Slice(tradesThisYear, func(i, j int) bool {
		return tradesThisYear[i].CloseTime.Compare(tradesThisYear[j].CloseTime) == -1
	})

	tradeTimeAxisData := make([]string, weekOfYear)

	currWeek := 0
	for _, trade := range tradesThisYear {
		_, tradeWeek := trade.CloseTime.ISOWeek()
		if tradeWeek != currWeek {
			currWeek = tradeWeek
			tradeTimeAxisData[currWeek-1] = trade.CloseTime.Format("01/02")

			profitsByWeek[currWeek-1] = 0
			if currWeek-1 > 0 {
				profitsByWeek[currWeek-1] = profitsByWeek[currWeek-2]
			}
		}

		profitsByWeek[currWeek-1] += trade.GetProfit()
	}

	tradeProfitChartData := make([]opts.LineData, len(profitsByWeek))
	for i, profit := range profitsByWeek {
		tradeProfitChartData[i] = opts.LineData{Value: profit}
	}

	c.Line = charts.NewLine()

	c.Line.SetXAxis(tradeTimeAxisData)
	c.Line.AddSeries("P&L", tradeProfitChartData)

	c.SetChartOptions()

	c.Line.Render(w)

	return nil
}

func (c *YearlyChart) SetTradeMode(includeSwings bool) {
	c.Portfolio.IncludeSwing = includeSwings
}
