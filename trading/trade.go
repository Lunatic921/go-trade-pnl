package trading

import (
	"math"
	"time"
)

type TradeSide string

const (
	UNKNOWN TradeSide = "UNKNOWN"
	LONG    TradeSide = "LONG"
	SHORT   TradeSide = "SHORT"
)

type TradeOperation string

const (
	TO_OPEN  = "TO OPEN"
	TO_CLOSE = "TO CLOSE"
)

func ParseOperation(s string) (t TradeOperation) {
	if s == TO_OPEN {
		return TO_OPEN
	} else {
		return TO_CLOSE
	}
}

type Trade struct {
	Ticker            string
	Side              TradeSide
	CurrentShareCount int
	TotalShareCount   int
	OpenTime          time.Time
	CloseTime         time.Time

	OpenExecutions  TradeExecutions
	CloseExecutions TradeExecutions
}

type Trades []Trade

func (t *Trade) execute(e TradeExecution) {

	if t.Side == UNKNOWN {
		if e.Qty > 0 {
			t.Side = LONG
		} else {
			t.Side = SHORT
		}
	}
	t.CurrentShareCount += e.Qty

	if t.CurrentShareCount == 0 {
		t.CloseTime = e.ExecTime
	} else if t.OpenTime == (time.Time{}) {
		t.OpenTime = e.ExecTime
	}

	if e.PosEffect == TO_OPEN {
		t.OpenExecutions = append(t.OpenExecutions, e)
		t.TotalShareCount += e.Qty
	} else {
		t.CloseExecutions = append(t.CloseExecutions, e)
	}
}

func (t *Trade) isOpen() (_ bool) {
	return t.CurrentShareCount != 0
}

func (t *Trade) GetProfit() (profit float64) {
	if t.isOpen() {
		return 0.0
	}

	openPrice, closePrice := 0.0, 0.0

	for _, e := range t.OpenExecutions {
		openPrice += math.Abs(float64(e.Qty) * e.Price)
	}

	for _, e := range t.CloseExecutions {
		closePrice += math.Abs(float64(e.Qty) * e.Price)
	}

	profit = float64(closePrice) - float64(openPrice)

	if t.Side == SHORT {
		profit *= -1
	}

	return profit
}

func (t *Trade) GetOpeningPriceAvg() (avgPrice float64) {
	return t.getPriceAvg(t.OpenExecutions)
}

func (t *Trade) GetClosingPriceAvg() (avgPrice float64) {
	return t.getPriceAvg(t.CloseExecutions)
}

func (t *Trade) getPriceAvg(executions TradeExecutions) (avgPrice float64) {
	avgPrice = 0.0

	for _, e := range executions {
		avgPrice += e.NetPrice
	}

	avgPrice /= float64(len(executions))

	return avgPrice
}

func (t *Trade) GetDuration() time.Duration {
	return time.Time.Sub(t.CloseTime, t.OpenTime)
}
