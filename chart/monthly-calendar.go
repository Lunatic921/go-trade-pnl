package chart

import (
	"fmt"
	"go-trade-pnl/trading"
	"io"
	"text/template"
	"time"
)

type MonthlyCalendar struct {
	Month     time.Time
	Portfolio *trading.Portfolio
}

type TradingDay struct {
	Day            time.Time
	DayVal         int
	Profit         string
	TradeCount     int
	Wins           int
	Losses         int
	WinLossPct     string
	DayResultClass string
}

type CalendarWeek struct {
	Days []TradingDay
}

const tmplPath = "chart/templates/monthly-calendar.html"

func (c *MonthlyCalendar) Draw(w io.Writer) error {
	tmpl := template.Must(template.ParseFiles(tmplPath))

	calDays := c.getCalendarDays()
	var calWeeks []CalendarWeek
	newWeek := CalendarWeek{}

	for _, calDay := range calDays {
		if len(newWeek.Days) == 7 {
			calWeeks = append(calWeeks, newWeek)
			newWeek = CalendarWeek{}
		}

		dailyTrades := c.Portfolio.GetTradesByDay(calDay)
		dailyProfit := c.Portfolio.GetProfitByDay(calDay)

		profitStr := fmt.Sprintf("$%.2f", dailyProfit)
		if calDay.Weekday() == time.Saturday || calDay.Weekday() == time.Sunday {
			profitStr = ""
		}

		dayResultClasses := ""
		if calDay.Compare(c.Month) < 0 || len(dailyTrades) == 0 {
			dayResultClasses = "no-trade-day"
		} else if dailyProfit > 0 {
			dayResultClasses = "green-day"
		} else if dailyProfit < 0 {
			dayResultClasses = "red-day"
		}

		dailyWins := 0
		for _, trade := range dailyTrades {
			if trade.GetProfit() >= 0 {
				dailyWins++
			}
		}

		winLossPct := 0.0
		if len(dailyTrades) > 0 {
			winLossPct = (float64(dailyWins) / float64(len(dailyTrades))) * 100.0
		}

		newWeek.Days = append(newWeek.Days, TradingDay{
			Day:            calDay,
			DayVal:         calDay.Day(),
			Profit:         profitStr,
			TradeCount:     len(dailyTrades),
			Wins:           dailyWins,
			Losses:         len(dailyTrades) - dailyWins,
			WinLossPct:     fmt.Sprintf("%.2f", winLossPct),
			DayResultClass: dayResultClasses,
		})
	}

	if len(newWeek.Days) > 0 && len(newWeek.Days) < 7 {
		calWeeks = append(calWeeks, newWeek)
	}

	data := struct {
		Weeks []CalendarWeek
	}{
		Weeks: calWeeks,
	}

	tmpl.Execute(w, data)
	return nil
}

func (c *MonthlyCalendar) getCalendarDays() []time.Time {
	days := make([]time.Time, 0, 42)

	loc, _ := time.LoadLocation("America/Phoenix")

	startDay := time.Date(c.Month.Year(), c.Month.Month(), 1, 0, 0, 0, 0, loc)

	weekDayDiff := int(startDay.Weekday() - time.Sunday)
	for weekDayDiff != 0 {
		prevMonthDay := startDay.AddDate(0, 0, -int(weekDayDiff))
		weekDayDiff--
		days = append(days, prevMonthDay)
	}

	monthEndDate := time.Date(c.Month.Year(), c.Month.Month()+1, 1, 0, 0, 0, 0, loc)

	for startDay != monthEndDate {
		days = append(days, startDay)
		startDay = startDay.AddDate(0, 0, 1)
	}

	return days
}
