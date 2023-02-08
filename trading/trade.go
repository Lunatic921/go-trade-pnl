package trading

import (
	"math"
	"time"
)

type TradeSide string

const (
	UNKNOWN TradeSide = "UNKNOWN"
	LONG              = "LONG"
	SHORT             = "SHORT"
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
	Ticker     string
	Side       TradeSide
	ShareCount int
	OpenTime   time.Time
	CloseTime  time.Time

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

		t.OpenTime = e.ExecTime
	}
	t.ShareCount += e.Qty

	if t.ShareCount == 0 {
		t.CloseTime = e.ExecTime
	}

	if e.PosEffect == TO_OPEN {
		t.OpenExecutions = append(t.OpenExecutions, e)
	} else {
		t.CloseExecutions = append(t.CloseExecutions, e)
	}
}

func (t *Trade) isOpen() (_ bool) {
	return t.ShareCount != 0
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
