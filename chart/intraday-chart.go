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

type TickerSummary struct {
	Ticker            string
	TradeCount        int
	WinningTradeCount int
	Profit            float64
}

const intraDayTmplPath = "chart/templates/intraday-details.html"

func (c *IntradayChart) Draw(w io.Writer) error {

	trades := c.Portfolio.FilterTrades(c.Day.Year(), int(c.Day.Month()), c.Day.Day())

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
		Trades        []ExtendedTrade
		TickerSummary []*TickerSummary
	}{
		Trades:        extendedTrades,
		TickerSummary: c.getProfitsPerTicker(trades),
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

func (c *IntradayChart) getProfitsPerTicker(trades []*trading.Trade) []*TickerSummary {
	tickerProfits := make(map[string]*TickerSummary)

	keys := make([]string, 0, len(trades))

	for _, trade := range trades {
		tickerProfit, ok := tickerProfits[trade.Ticker]
		tradeProfit := trade.GetProfit()
		if !ok {
			tickerProfit = &TickerSummary{
				Ticker:            trade.Ticker,
				TradeCount:        0,
				WinningTradeCount: 0,
				Profit:            0.0,
			}

			tickerProfits[trade.Ticker] = tickerProfit

			keys = append(keys, trade.Ticker)
		}

		tickerProfit.Profit += tradeProfit
		tickerProfit.TradeCount += 1
		if tradeProfit > 0 {
			tickerProfit.WinningTradeCount += 1
		}
	}

	tickerProfitList := make([]*TickerSummary, len(keys))

	for i, key := range keys {
		tickerProfitList[i] = tickerProfits[key]
	}

	sort.SliceStable(tickerProfitList, func(i, j int) bool {
		return tickerProfitList[i].Profit > tickerProfitList[j].Profit
	})

	return tickerProfitList
}

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
		openTotal += (trade.NetPrice * float64(trade.Qty))
		openCount += trade.Qty
	}

	return fmt.Sprintf("$%.2f", openTotal/float64(openCount))
}

func (et *ExtendedTrade) GetExitPrice() string {
	exitTotal := 0.0
	exitCount := 0

	for _, trade := range et.Trade.CloseExecutions {
		exitTotal += (trade.NetPrice * float64(trade.Qty))
		exitCount += trade.Qty
	}

	return fmt.Sprintf("$%.2f", math.Abs(exitTotal/float64(exitCount)))
}

func (ts *TickerSummary) GetWinPercentage() string {
	return fmt.Sprintf("%.1f %%", (100 * float64(ts.WinningTradeCount) / float64(ts.TradeCount)))
}

func (ts *TickerSummary) GetProfit() string {
	return fmt.Sprintf("$%.2f", ts.Profit)
}
