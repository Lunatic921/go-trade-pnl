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

type IntradayChart struct {
	LineChart
	Day       time.Time
	Portfolio *trading.Portfolio
}

const intraDayTmplPath = "chart/templates/intraday-details.html"

func (c *IntradayChart) Draw(w io.Writer) error {

	trades := c.Portfolio.GetTradesByDay(c.Day)

	tradeTimeAxisData := make([]string, len(trades))
	for i, trade := range trades {
		tradeTimeAxisData[i] = trade.CloseTime.Format("15:04:05")
	}

	tradeProfitChartData := make([]opts.LineData, len(trades))
	var previousProfit float64 = 0.0
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

	data := struct {
		ProfitsByTicker []string
		TradeDetails    []string
	}{
		ProfitsByTicker: c.getProfitsPerTicker(trades),
		TradeDetails:    c.getTradeDetails(trades),
	}

	tmpl := template.Must(template.ParseFiles(intraDayTmplPath))
	tmpl.Execute(w, data)

	return nil
}

func (c *IntradayChart) getProfitsPerTicker(trades []*trading.Trade) []string {
	tickerProfits := make(map[string]float64)

	keys := make([]string, 0, len(trades))

	for _, trade := range trades {
		_, ok := tickerProfits[trade.Ticker]
		if ok {
			tickerProfits[trade.Ticker] += trade.GetProfit()
		} else {
			tickerProfits[trade.Ticker] = trade.GetProfit()
			keys = append(keys, trade.Ticker)
		}
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return tickerProfits[keys[i]] > tickerProfits[keys[j]]
	})

	tickerProfitStrs := make([]string, len(keys))

	for i, key := range keys {
		tickerProfitStrs[i] = fmt.Sprintf("%s: $%0.2f", key, tickerProfits[key])
	}

	return tickerProfitStrs
}

func (c *IntradayChart) getTradeDetails(trades []*trading.Trade) []string {
	tradeDetails := make([]string, len(trades))

	// sort.Slice(trades, func(i, j int) bool {
	// 	return trades[i].GetProfit() > trades[j].GetProfit()
	// })

	for i, trade := range trades {
		avgOpenPrice := trade.GetOpeningPriceAvg()
		avgClosePrice := trade.GetClosingPriceAvg()

		detail := fmt.Sprintf("%d) %s: %d shares ($%0.2f->$%0.2f)   $%0.2f", i+1, trade.Ticker,
			trade.TotalShareCount, avgOpenPrice, avgClosePrice, trade.GetProfit())

		tradeDetails[i] = detail
	}

	return tradeDetails
}
