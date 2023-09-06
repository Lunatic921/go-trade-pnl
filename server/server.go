package server

import (
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

type TradeServer struct {
	pages map[string]*TradePage
}

func (srv *TradeServer) CreatePath(path string, page *TradePage) {

	if srv.pages == nil {
		srv.pages = make(map[string]*TradePage)
	}

	existingPage, ok := srv.pages[path]
	if !ok {
		srv.pages[path] = page
		http.HandleFunc(path, page.httpServer)
	} else {
		existingPage.Chart = page.Chart
	}
}
