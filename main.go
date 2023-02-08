package main

import (
	"flag"
	"fmt"
	"go-learn-pl/chart"
	"go-learn-pl/server"
	"go-learn-pl/trading"
	"time"
)

func main() {
	historyFilePathPtr := flag.String("trades", "", "TOS trade history")
	flag.Parse()

	p := trading.NewPortfolio(*historyFilePathPtr)

	for _, trade := range p.Trades {
		fmt.Println("Ticker: ", trade.Ticker)
		fmt.Println("Profit: ", fmt.Sprintf("%.2f", trade.GetProfit()))
		fmt.Println("----------------------")
	}

	tradingDays := p.GetTradingDays()

	for _, date := range tradingDays {
		intraDayChart := &chart.IntradayChart{Day: date, Portfolio: p}
		srv := server.TradePage{Chart: intraDayChart}
		srv.CreatePath(fmt.Sprintf("/%s", date.Format("2006-01-02")), intraDayChart)
	}

	loc, _ := time.LoadLocation("America/Phoenix")

	monthAfterEndMonth := time.Date(tradingDays[len(tradingDays)-1].Year(), tradingDays[len(tradingDays)-1].Month()+1, 0, 0, 0, 0, 0, loc)

	for currentMonth := time.Date(tradingDays[0].Year(), tradingDays[0].Month(), 1, 0, 0, 0, 0, loc); currentMonth.Before(monthAfterEndMonth); currentMonth = currentMonth.AddDate(0, 1, 0) {
		monthlyChart := &chart.MonthlyChart{Month: currentMonth, Portfolio: p}
		srv := server.TradePage{Chart: monthlyChart}
		srv.CreatePath(fmt.Sprintf("/%s", currentMonth.Format("2006-01")), monthlyChart)
	}

	firstYear := time.Date(tradingDays[0].Year(), 1, 1, 0, 0, 0, 0, loc)
	yearAfterLastYear := time.Date(tradingDays[len(tradingDays)-1].Year()+1, 1, 1, 0, 0, 0, 0, loc)

	for currYear := firstYear; currYear.Before(yearAfterLastYear); currYear = currYear.AddDate(1, 0, 0) {
		fmt.Println("Year:", currYear.Year())
		yearlyChart := &chart.YearlyChart{Year: currYear, Portfolio: p}
		srv := server.TradePage{Chart: yearlyChart}
		srv.CreatePath(fmt.Sprintf("/%s", currYear.Format("2006")), yearlyChart)
	}

	server.Serve(8081)

}
