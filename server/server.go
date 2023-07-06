package server

import (
	"fmt"
	"go-trade-pnl/chart"
	"net/http"
)

type TradePage struct {
	Chart chart.Chart
}

func (t *TradePage) httpServer(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("swings") == "1" {
		t.Chart.SetTradeMode(true)
	} else {
		t.Chart.SetTradeMode(false)
	}
	t.Chart.Draw(w)
}

func (t *TradePage) CreatePath(path string, chart chart.Chart) {
	http.HandleFunc(path, t.httpServer)
}

func Serve(port int) {
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
