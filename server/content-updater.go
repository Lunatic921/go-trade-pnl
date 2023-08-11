package server

import (
	"fmt"
	"go-trade-pnl/chart"
	"go-trade-pnl/trading"

	"os"
	"time"

	"github.com/radovskyb/watcher"
)

type ContentUpdater interface {
	Monitor(updateStream chan any)
	Update()
}

func AddContentUpdater(updater ContentUpdater, updateStream chan any) {
	updater.Update()
	go updater.Monitor(updateStream)
}

type TradingJournalUpdater struct {
	Directory string
	Server    *TradeServer
}

func (tju *TradingJournalUpdater) Monitor(updateStream chan any) {
	dirWatcher := watcher.New()
	dirWatcher.Add(tju.Directory)
	go dirWatcher.Start(time.Second * 5)

	for {
		select {
		case <-dirWatcher.Event:
			tju.Update()
			updateStream <- true
		case error := <-dirWatcher.Error:
			fmt.Printf("Error watching: %+v\n", error)
			os.Exit(1)
		case close := <-dirWatcher.Closed:
			fmt.Printf("Watcher closed: %+v\n", close)
			os.Exit(1)
		}
	}
}

func (tju *TradingJournalUpdater) Update() {
	p := trading.NewPortfolio(tju.Directory)

	tradingDays := p.GetTradingDays(-1, -1, -1)

	srv := tju.Server

	for _, date := range tradingDays {
		intraDayChart := &chart.IntradayChart{Day: date, Portfolio: p}
		page := TradePage{Chart: intraDayChart}
		srv.CreatePath(fmt.Sprintf("/%s", date.Format("2006-01-02")), page)
	}

	loc, _ := time.LoadLocation("America/Phoenix")

	monthAfterEndMonth := time.Date(tradingDays[len(tradingDays)-1].Year(), tradingDays[len(tradingDays)-1].Month()+1, 0, 0, 0, 0, 0, loc)

	for currentMonth := time.Date(tradingDays[0].Year(), tradingDays[0].Month(), 1, 0, 0, 0, 0, loc); currentMonth.Before(monthAfterEndMonth); currentMonth = currentMonth.AddDate(0, 1, 0) {
		monthFormatted := currentMonth.Format("2006-01")
		monthlyChart := &chart.MonthlyChart{Month: currentMonth, Portfolio: p}
		page := TradePage{Chart: monthlyChart}
		srv.CreatePath(fmt.Sprintf("/%s", monthFormatted), page)

		monthlyCalendar := &chart.MonthlyCalendar{Month: currentMonth, Portfolio: p}
		page = TradePage{Chart: monthlyCalendar}
		srv.CreatePath(fmt.Sprintf("/calendar/%s", monthFormatted), page)

		monthlyProfitPrice := &chart.MonthlyPriceProfitChart{Month: currentMonth, Portfolio: p}
		page = TradePage{Chart: monthlyProfitPrice}
		srv.CreatePath(fmt.Sprintf("/profitprice/%s", monthFormatted), page)

		monthlySharePrice := &chart.MonthlyPriceShareChart{Month: currentMonth, Portfolio: p}
		page = TradePage{Chart: monthlySharePrice}
		srv.CreatePath(fmt.Sprintf("/shareprice/%s", monthFormatted), page)

		timeOfDayChart := &chart.TimeOfDayChart{Month: currentMonth, Portfolio: p}
		page = TradePage{Chart: timeOfDayChart}
		srv.CreatePath(fmt.Sprintf("/tod/%s", currentMonth.Format("2006-01")), page)

		durationChart := &chart.MonthlyDurationScatterChart{Month: currentMonth, Portfolio: p}
		page = TradePage{Chart: durationChart}
		srv.CreatePath(fmt.Sprintf("/duration/%s", currentMonth.Format("2006-01")), page)
	}

	firstYear := time.Date(tradingDays[0].Year(), 1, 1, 0, 0, 0, 0, loc)
	yearAfterLastYear := time.Date(tradingDays[len(tradingDays)-1].Year()+1, 1, 1, 0, 0, 0, 0, loc)

	for currYear := firstYear; currYear.Before(yearAfterLastYear); currYear = currYear.AddDate(1, 0, 0) {
		yearlyChart := &chart.YearlyChart{Year: currYear, Portfolio: p}
		page := TradePage{Chart: yearlyChart}
		srv.CreatePath(fmt.Sprintf("/%s", currYear.Format("2006")), page)
	}
}
