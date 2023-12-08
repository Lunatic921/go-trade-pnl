package chart

import (
	"fmt"
	"go-trade-pnl/trading"
	"io"
	"sort"
	"text/template"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type YearlyChart struct {
	LineChart
	Year      time.Time
	Portfolio *trading.Portfolio
}

type YearlyStat struct {
	StatName  string
	StatValue string
}

const yearlyStatsTmplPath = "chart/templates/year-stats.html"

func (c *YearlyChart) Draw(w io.Writer) error {
	timeNow := time.Now()
	_, weekOfYear := timeNow.ISOWeek()
	if c.Year.Year() != timeNow.Year() {
		year := c.Year
		_, weekOfYear = year.Add(-time.Minute).ISOWeek()
	}
	profitsByWeek := make([]float64, weekOfYear)

	tradesThisYear := c.Portfolio.FilterTrades(c.Year.Year(), -1, -1)

	sort.Slice(tradesThisYear, func(i, j int) bool {
		return tradesThisYear[i].CloseTime.Compare(tradesThisYear[j].CloseTime) == -1
	})

	tradeTimeAxisData := make([]string, weekOfYear)

	currWeek := 0
	for _, trade := range tradesThisYear {
		_, tradeWeek := trade.CloseTime.ISOWeek()
		if tradeWeek != currWeek {
			for tradeWeek != currWeek {
				currWeek += 1
				tradeTimeAxisData[currWeek-1] = trade.CloseTime.Format("01/02")

				profitsByWeek[currWeek-1] = 0
				if currWeek-1 > 0 {
					profitsByWeek[currWeek-1] = profitsByWeek[currWeek-2]
				}
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

	data := struct {
		Stats []YearlyStat
	}{
		Stats: c.getYearlyStats(),
	}

	tmpl := template.Must(template.ParseFiles(yearlyStatsTmplPath))
	err := tmpl.Execute(w, data)
	if err != nil {
		fmt.Printf("Err: %s\n", err.Error())
	}

	return nil
}

func (c *YearlyChart) SetTradeMode(includeSwings bool) {
	c.Portfolio.IncludeSwing = includeSwings
}

func (c *YearlyChart) getYearlyStats() []YearlyStat {

	yearProfit := c.Portfolio.GetProfit(c.Year.Year(), -1, -1)
	tradingDays := c.Portfolio.GetTradingDays(c.Year.Year(), -1, -1)
	greenDays, redDays := c.Portfolio.GetGreenVsRedDays(c.Year.Year(), -1, -1)
	avgDailyPl := fmt.Sprintf("$%0.2f", yearProfit/float64(len(tradingDays)))
	dailyProfitPerShare := fmt.Sprintf("$%0.2f", c.Portfolio.GetProfitPerShare(c.Year.Year(), -1, -1)/float64(len(tradingDays)))
	avgTradePl := fmt.Sprintf("$%0.4f", c.Portfolio.GetTradePl(c.Year.Year(), -1, -1))

	stats := []YearlyStat{
		{StatName: "Trading Days", StatValue: fmt.Sprintf("%d", len(tradingDays))},
		{StatName: "Green Days", StatValue: fmt.Sprintf("%d", greenDays)},
		{StatName: "Red Days", StatValue: fmt.Sprintf("%d", redDays)},
		{StatName: "Trades", StatValue: fmt.Sprintf("%d", len(c.Portfolio.FilterTrades(c.Year.Year(), -1, -1)))},
		{StatName: "Total P/L", StatValue: fmt.Sprintf("$%0.2f", yearProfit)},
		{StatName: "Trade P/L", StatValue: avgTradePl},
		{StatName: "Win Pct", StatValue: fmt.Sprintf("%.2f%%", 100.0*c.Portfolio.GetWinPercentage(c.Year.Year(), -1, -1))},
		{StatName: "Daily P/L", StatValue: avgDailyPl},
		{StatName: "Daily $/Share", StatValue: dailyProfitPerShare},
	}

	return stats
}
