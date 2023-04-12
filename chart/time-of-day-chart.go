package chart

import (
	"go-trade-pnl/trading"
	"io"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type TimeOfDayChart struct {
	bar       *charts.Bar
	Month     time.Time
	Portfolio *trading.Portfolio

	tradesByTimeOfDay map[time.Time][]trading.Trade
}

func (c *TimeOfDayChart) Draw(w io.Writer) error {

	c.bar = charts.NewBar()

	intervals := c.generateIntervals()

	c.tradesByTimeOfDay = make(map[time.Time][]trading.Trade)
	easternTimezone, _ := time.LoadLocation("America/New_York")

	for _, interval := range intervals {
		c.tradesByTimeOfDay[interval] = make([]trading.Trade, 0)
	}

	trades := c.Portfolio.GetTradesByMonth(c.Month)
	for _, trade := range trades {
		tradeTime := trade.OpenTime.In(easternTimezone)
		for i, interval := range intervals {
			if tradeTime.Hour() >= interval.Hour() && tradeTime.Minute() >= interval.Minute() {
				if i+1 < len(intervals) && (intervals[i+1].Hour() > tradeTime.Hour() || (intervals[i+1].Hour() == tradeTime.Hour() && intervals[i+1].Minute() > tradeTime.Minute())) {
					c.tradesByTimeOfDay[interval] = append(c.tradesByTimeOfDay[interval], trade)
					break
				}
			}
		}
	}

	timeScale := make([]opts.BarData, 0)

	for _, interval := range intervals {
		timeScale = append(timeScale, opts.BarData{Value: interval.Format("15:04")})
	}

	c.bar.SetXAxis(timeScale)

	c.bar.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{
			AxisLabel: &opts.AxisLabel{Show: true, Interval: "0"}}))

	profits := make([]float64, 0)
	for i, interval := range intervals {
		profits = append(profits, 0)
		for _, trade := range c.tradesByTimeOfDay[interval] {
			profits[i] += trade.GetProfit()
		}
	}

	ySeries := make([]opts.BarData, 0)
	for _, profit := range profits {
		ySeries = append(ySeries, opts.BarData{Value: profit})
	}

	c.bar.AddSeries("Profit", ySeries)

	c.bar.Render(w)

	return nil
}

func (c *TimeOfDayChart) generateIntervals() []time.Time {
	easternTimezone, _ := time.LoadLocation("America/New_York")
	currTime := time.Date(c.Month.Year(), c.Month.Month(), c.Month.Day(), 8, 30, 0, 0, easternTimezone)
	lastTime := time.Date(c.Month.Year(), c.Month.Month(), c.Month.Day(), 18, 0, 0, 0, easternTimezone)

	intervals := make([]time.Time, 0, 18)

	for currTime.Before(lastTime) {
		intervals = append(intervals, currTime)
		currTime = currTime.Add(time.Duration(30) * time.Minute)
	}

	return intervals
}
