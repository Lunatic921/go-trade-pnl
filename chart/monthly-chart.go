package chart

import (
	"go-trade-pnl/trading"
	"io"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type MonthlyChart struct {
	LineChart
	Month     time.Time
	Portfolio *trading.Portfolio
}

func (c *MonthlyChart) Draw(w io.Writer) error {
	loc, _ := time.LoadLocation("America/Phoenix")
	startDay := time.Date(c.Month.Year(), c.Month.Month(), 1, 0, 0, 0, 0, loc)
	nextMonth := time.Date(c.Month.Year(), c.Month.Month()+1, 1, 0, 0, 0, 0, loc)

	monthsTrades := c.Portfolio.FilterTrades(c.Month.Year(), int(c.Month.Month()), -1)

	var tradeTimeAxisData []string
	var tradeProfitChartData []opts.LineData

	runningTotal := 0.0
	for currDay := startDay; currDay.Before(nextMonth); currDay = currDay.AddDate(0, 0, 1) {
		dayProfit := 0.0
		tradeCount := 0
		for _, trade := range monthsTrades {
			if trade.CloseTime.Day() != currDay.Day() {
				if tradeCount > 0 {
					tradeTimeAxisData = append(tradeTimeAxisData, currDay.Format("1/2"))
					runningTotal += dayProfit
					tradeProfitChartData = append(tradeProfitChartData, opts.LineData{Value: runningTotal})
					monthsTrades = monthsTrades[tradeCount:]

					tradeCount = 0
					dayProfit = 0
				}
				break
			}

			tradeCount++
			dayProfit += trade.GetProfit()
		}

		//Add the last day
		if tradeCount > 0 {
			tradeTimeAxisData = append(tradeTimeAxisData, currDay.Format("1/2"))
			runningTotal += dayProfit
			tradeProfitChartData = append(tradeProfitChartData, opts.LineData{Value: runningTotal})
		}
	}

	c.Line = charts.NewLine()

	c.Line.SetXAxis(tradeTimeAxisData)
	c.Line.AddSeries("P&L", tradeProfitChartData)

	c.SetChartOptions()

	c.Line.Render(w)

	return nil
}

func (c *MonthlyChart) SetTradeMode(includeSwings bool) {
	c.Portfolio.IncludeSwing = includeSwings
}
