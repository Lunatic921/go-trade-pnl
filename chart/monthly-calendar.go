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
	Day             time.Time
	FormattedDay    string
	DayVal          int
	Profit          string
	TradeCount      int
	Wins            int
	Losses          int
	WinLossPct      string
	DayResultClass  string
	DollarsPerShare string
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
	monthlyProfit := 0.0
	monthlyWins, monthlyTrades, daysTraded := 0, 0, 0

	for _, calDay := range calDays {
		if len(newWeek.Days) == 7 {
			calWeeks = append(calWeeks, newWeek)
			newWeek = CalendarWeek{}
		}

		if calDay.Compare(c.Month) == -1 {
			newWeek.Days = append(newWeek.Days, TradingDay{
				Day:            calDay,
				DayVal:         calDay.Day(),
				DayResultClass: "no-trade-day",
				FormattedDay:   calDay.Format("2006-01-02"),
			})
			continue
		}

		dailyTrades := c.Portfolio.FilterTrades(calDay.Year(), int(calDay.Month()), calDay.Day())
		dailyProfit := c.Portfolio.GetProfit(calDay.Year(), int(calDay.Month()), calDay.Day())

		monthlyProfit += dailyProfit

		profitStr := fmt.Sprintf("$%.2f", dailyProfit)
		if calDay.Weekday() == time.Saturday || calDay.Weekday() == time.Sunday {
			profitStr = ""
		}

		dayResultClasses := ""
		if calDay.Compare(c.Month) < 0 || len(dailyTrades) == 0 {
			dayResultClasses = "no-trade-day"
		} else {
			daysTraded++

			if dailyProfit > 0 {
				dayResultClasses = "green-day"
			} else if dailyProfit < 0 {
				dayResultClasses = "red-day"
			}
		}

		dailyWins := 0
		for _, trade := range dailyTrades {
			if trade.GetProfit() >= 0 {
				dailyWins++
			}

			monthlyTrades++
		}

		monthlyWins += dailyWins

		winLossPct := 0.0
		if len(dailyTrades) > 0 {
			winLossPct = (float64(dailyWins) / float64(len(dailyTrades))) * 100.0
		}

		numberOfDailyTrades := float64(len(c.Portfolio.FilterTrades(calDay.Year(), int(calDay.Month()), calDay.Day())))

		sharesTradedOnDay := float64(c.Portfolio.GetSharesTraded(calDay.Year(), int(calDay.Month()), calDay.Day())) / numberOfDailyTrades

		dollarsPerShare := dailyProfit / sharesTradedOnDay

		newWeek.Days = append(newWeek.Days, TradingDay{
			Day:             calDay,
			DayVal:          calDay.Day(),
			Profit:          profitStr,
			TradeCount:      len(dailyTrades),
			Wins:            dailyWins,
			Losses:          len(dailyTrades) - dailyWins,
			WinLossPct:      fmt.Sprintf("%.2f", winLossPct),
			DayResultClass:  dayResultClasses,
			FormattedDay:    calDay.Format("2006-01-02"),
			DollarsPerShare: fmt.Sprintf("$%.2f", dollarsPerShare),
		})
	}

	if len(newWeek.Days) > 0 && len(newWeek.Days) < 7 {
		calWeeks = append(calWeeks, newWeek)
	}

	winRate := fmt.Sprintf("%.2f%%", (float64(monthlyWins)/float64(monthlyTrades))*100.0)
	dailyAvg := fmt.Sprintf("$%.2f", monthlyProfit/float64(daysTraded))

	sharesPerTrade := float64(c.Portfolio.GetSharesTraded(c.Month.Year(), int(c.Month.Month()), -1)) / float64(monthlyTrades)
	dollarsPerShareMonth := fmt.Sprintf("$%.2f", (monthlyProfit/sharesPerTrade)/float64(daysTraded))

	data := struct {
		Weeks                  []CalendarWeek
		MonthlyProfit          string
		MonthlyDollarsPerShare string
		WinRate                string
		WinningTrades          int
		TotalTrades            int
		DailyAvg               string
		DailyTrades            int
	}{
		Weeks:                  calWeeks,
		MonthlyProfit:          fmt.Sprintf("$%.2f", monthlyProfit),
		MonthlyDollarsPerShare: dollarsPerShareMonth,
		WinRate:                winRate,
		WinningTrades:          monthlyWins,
		TotalTrades:            monthlyTrades,
		DailyAvg:               dailyAvg,
		DailyTrades:            monthlyTrades / daysTraded,
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

func (c *MonthlyCalendar) SetTradeMode(includeSwings bool) {
	c.Portfolio.IncludeSwing = includeSwings
}
