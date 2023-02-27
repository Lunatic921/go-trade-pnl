package main

import (
	"flag"
	"fmt"
	"go-trade-pnl/chart"
	"go-trade-pnl/server"
	"go-trade-pnl/trading"
	"time"
)

func main() {
	historyFilePathPtr := flag.String("trades", "", "TOS trade history")
	flag.Parse()

	p := trading.NewPortfolio(*historyFilePathPtr)

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

		monthlyCalendar := &chart.MonthlyCalendar{Month: currentMonth, Portfolio: p}
		srv2 := server.TradePage{Chart: monthlyCalendar}
		srv2.CreatePath(fmt.Sprintf("/calendar/%s", currentMonth.Format("2006-01")), monthlyCalendar)
	}

	firstYear := time.Date(tradingDays[0].Year(), 1, 1, 0, 0, 0, 0, loc)
	yearAfterLastYear := time.Date(tradingDays[len(tradingDays)-1].Year()+1, 1, 1, 0, 0, 0, 0, loc)

	for currYear := firstYear; currYear.Before(yearAfterLastYear); currYear = currYear.AddDate(1, 0, 0) {
		yearlyChart := &chart.YearlyChart{Year: currYear, Portfolio: p}
		srv := server.TradePage{Chart: yearlyChart}
		srv.CreatePath(fmt.Sprintf("/%s", currYear.Format("2006")), yearlyChart)
	}

	server.Serve(8081)

}
