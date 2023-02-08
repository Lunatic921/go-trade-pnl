package trading

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Portfolio struct {
	filePath string
	Trades   []*Trade
}

func NewPortfolio(filePath string) *Portfolio {
	p := &Portfolio{filePath: filePath, Trades: make([]*Trade, 0, 1000)}
	p.parseTradeFile()
	return p
}

func (p *Portfolio) GetTradeCount() int {
	return len(p.Trades)
}

func (p *Portfolio) GetTradesByDay(day time.Time) []*Trade {
	startIdx, endIdx := -1, -1

	for i, trade := range p.Trades {

		if trade.CloseTime.Year() == day.Year() && trade.CloseTime.Month() == day.Month() && trade.CloseTime.Day() == day.Day() {
			if startIdx == -1 {
				startIdx = i
			}
		} else {
			if startIdx != -1 {
				endIdx = i - 1
			}
		}
	}

	if startIdx != -1 && endIdx == -1 {
		endIdx = len(p.Trades) - 1
	}

	return p.Trades[startIdx:endIdx]
}

func (p *Portfolio) GetTradesByMonth(month time.Time) []Trade {
	trades := make([]Trade, 0, 3100)

	for _, trade := range p.Trades {

		if trade.CloseTime.Year() == month.Year() && trade.CloseTime.Month() == month.Month() {
			trades = append(trades, *trade)
		}
	}

	return trades
}

func (p *Portfolio) GetTradesByYear(year time.Time) []Trade {
	trades := make([]Trade, 0, 40000)
	for _, trade := range p.Trades {
		if trade.CloseTime.Year() == year.Year() {
			trades = append(trades, *trade)
		}
	}

	return trades
}

func (p *Portfolio) GetTradingDays() []time.Time {
	days := make([]time.Time, 0, 365)

	for _, trade := range p.Trades {
		foundDay := false

		loc, _ := time.LoadLocation("America/Phoenix")
		newYear, newMonth, newDay := trade.CloseTime.Date()
		newDate := time.Date(newYear, newMonth, newDay, 0, 0, 0, 0, loc)

		for _, day := range days {
			oldYear, oldMonth, oldDay := day.Date()

			if newYear == oldYear && newMonth == oldMonth && oldDay == newDay {
				foundDay = true
				break
			}
		}

		if !foundDay {
			days = append(days, newDate)
		}
	}

	sort.Slice(days, func(i, j int) bool {
		return days[i].Before(days[j])
	})

	return days
}

func (p *Portfolio) parseTradeFile() {

	file, err := os.Open(p.filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	startLine, endLine := 0, 0

	for i, line := range lines {
		if line == "Account Trade History" {
			startLine = i + 2
			for j, line2 := range lines[startLine:] {
				if line2 == "" {
					endLine = j + startLine
					break
				}
			}
		}
	}

	reader := csv.NewReader(strings.NewReader(strings.Join(lines[startLine:endLine], "\n")))
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV data:", err)
		return
	}

	var trades TradeExecutions
	for _, record := range records {
		layout := "1/2/06 15:04:05"

		execTime, _ := time.Parse(layout, record[1])
		qty, _ := strconv.Atoi(record[4])
		price, _ := strconv.ParseFloat(record[10], 64)
		netPrice, _ := strconv.ParseFloat(record[11], 64)

		trades = append(trades, TradeExecution{
			ExecTime:  execTime,
			Spread:    record[2],
			Side:      TradeSide(record[3]),
			Qty:       qty,
			PosEffect: ParseOperation(record[5]),
			Symbol:    record[6],
			Exp:       record[7],
			Strike:    record[8],
			Type:      record[9],
			Price:     price,
			NetPrice:  netPrice,
			OrderType: record[12],
		})
	}

	sort.Sort(trades)

	//Break trade executions into their respective Trades
	for _, tradeEx := range trades {
		foundTrade := false

		//Find an open trade
		for _, trade := range p.Trades {

			if trade.Ticker == tradeEx.Symbol && trade.isOpen() {
				foundTrade = true
				trade.execute(tradeEx)

				break
			}
		}

		if !foundTrade {
			trade := &Trade{Ticker: tradeEx.Symbol, Side: tradeEx.Side}
			trade.execute(tradeEx)

			p.Trades = append(p.Trades, trade)
		}
	}
}
