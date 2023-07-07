package chart

import (
	"fmt"
	"go-trade-pnl/trading"
	"io"
	"math"
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

type ExtendedTrade struct {
	*trading.Trade
}

const intraDayTmplPath = "chart/templates/intraday-details.html"

func (c *IntradayChart) Draw(w io.Writer) error {

	trades := c.Portfolio.GetTradesByDay(c.Day)

	sort.SliceStable(trades, func(i, j int) bool {
		return trades[i].CloseTime.Compare(trades[j].CloseTime) == -1
	})

	easternTimezone, _ := time.LoadLocation("America/New_York")

	tradeTimeAxisData := make([]string, len(trades))
	for i, trade := range trades {
		tradeTimeAxisData[i] = trade.CloseTime.In(easternTimezone).Format("15:04:05")
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

	extendedTrades := make([]ExtendedTrade, len(trades))
	for i, trade := range trades {
		extTrade := ExtendedTrade{trade}
		extendedTrades[i] = extTrade
	}

	data := struct {
		Trades []ExtendedTrade
	}{
		Trades: extendedTrades,
	}

	tmpl := template.Must(template.ParseFiles(intraDayTmplPath))
	err := tmpl.Execute(w, data)
	if err != nil {
		fmt.Printf("Err: %s\n", err.Error())
	}

	return nil
}

func (c *IntradayChart) SetTradeMode(includeSwings bool) {
	c.Portfolio.IncludeSwing = includeSwings
}

// func (c *IntradayChart) getProfitsPerTicker(trades []*trading.Trade) []string {
// 	tickerProfits := make(map[string]float64)

// 	keys := make([]string, 0, len(trades))

// 	for _, trade := range trades {
// 		_, ok := tickerProfits[trade.Ticker]
// 		if ok {
// 			tickerProfits[trade.Ticker] += trade.GetProfit()
// 		} else {
// 			tickerProfits[trade.Ticker] = trade.GetProfit()
// 			keys = append(keys, trade.Ticker)
// 		}
// 	}

// 	sort.SliceStable(keys, func(i, j int) bool {
// 		return tickerProfits[keys[i]] > tickerProfits[keys[j]]
// 	})

// 	tickerProfitStrs := make([]string, len(keys))

// 	for i, key := range keys {
// 		tickerProfitStrs[i] = fmt.Sprintf("%s: $%0.2f", key, tickerProfits[key])
// 	}

// 	return tickerProfitStrs
// }

// func (c *IntradayChart) getTradeDetails(trades []*trading.Trade) []string {
// 	tradeDetails := make([]string, len(trades))

// 	for i, trade := range trades {
// 		avgOpenPrice := trade.GetOpeningPriceAvg()
// 		avgClosePrice := trade.GetClosingPriceAvg()

// 		detail := fmt.Sprintf("%d) %s: %d shares ($%0.2f->$%0.2f)   $%0.2f", i+1, trade.Ticker,
// 			trade.TotalShareCount, avgOpenPrice, avgClosePrice, trade.GetProfit())

// 		tradeDetails[i] = detail
// 	}

// 	return tradeDetails
// }

func (et *ExtendedTrade) GetOpenTime() string {
	return et.OpenTime.Format("15:04:05")
}

func (et *ExtendedTrade) GetCloseTime() string {
	return et.CloseTime.Format("15:04:05")
}

func (et *ExtendedTrade) GetProfit() string {
	tradeProfit := et.Trade.GetProfit()
	minusSign := ""
	if tradeProfit < 0 {
		minusSign = "-"
	}

	return fmt.Sprintf("%s$%.2f", minusSign, math.Abs(et.Trade.GetProfit()))
}

func (et *ExtendedTrade) GetEntryPrice() string {
	openTotal := 0.0
	openCount := 0

	for _, trade := range et.Trade.OpenExecutions {
		openTotal += trade.NetPrice
		openCount += 1
	}

	return fmt.Sprintf("$%.2f", openTotal/float64(openCount))
}

func (et *ExtendedTrade) GetExitPrice() string {
	exitTotal := 0.0
	exitCount := 0

	for _, trade := range et.Trade.CloseExecutions {
		exitTotal += trade.NetPrice
		exitCount += 1
	}

	return fmt.Sprintf("$%.2f", math.Abs(exitTotal/float64(exitCount)))
}
