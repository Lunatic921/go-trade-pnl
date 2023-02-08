package chart

import (
	"go-trade-pnl/trading"
	"io"
	"time"
)

type YearlyChart struct {
	LineChart
	Year      time.Time
	Portfolio *trading.Portfolio
}

func (c *YearlyChart) Draw(w io.Writer) error {
	// TODO
	return nil
}
