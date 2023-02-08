package chart

import (
	"go-trade-pnl/trading"
	"io"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type IntradayChart struct {
	LineChart
	Day       time.Time
	Portfolio *trading.Portfolio
}

func (c *IntradayChart) Draw(w io.Writer) error {

	trades := c.Portfolio.GetTradesByDay(c.Day)

	tradeTimeAxisData := make([]string, len(trades), len(trades))
	for i, trade := range trades {
		tradeTimeAxisData[i] = trade.CloseTime.Format("15:04:05")
	}

	tradeProfitChartData := make([]opts.LineData, len(trades), len(trades))
	var previousProfit float64 = 0
	for i, trade := range trades {
		profit := trade.GetProfit()
		tradeProfitChartData[i] = opts.LineData{Value: profit + previousProfit}
		previousProfit += profit
	}

	c.Line = charts.NewLine()

	c.Line.SetXAxis(tradeTimeAxisData)
	c.Line.AddSeries("P&L", tradeProfitChartData)

	c.SetChartOptions()

	c.Line.Render(w)

	return nil
}
