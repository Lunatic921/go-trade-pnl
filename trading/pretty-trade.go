package trading

import "fmt"

type PrettyTrade struct {
	Trade
}

func (t *PrettyTrade) GetPercentGain() string {
	return fmt.Sprintf("%.2f%%", t.Trade.GetPercentGain()*100.0)
}

func (t *PrettyTrade) GetProfit() string {
	profit := t.Trade.GetProfit()

	negativeSign := ""
	if profit < 0.0 {
		negativeSign = "-"
		profit *= -1.0
	}

	return fmt.Sprintf("%s$%.2f", negativeSign, profit)
}
