package server

import (
	"fmt"
	"go-learn-pl/chart"
	"net/http"
)

type TradePage struct {
	Chart chart.Chart
}

func (t *TradePage) httpServer(w http.ResponseWriter, _ *http.Request) {
	t.Chart.Draw(w)
}

func (t *TradePage) CreatePath(path string, chart chart.Chart) {
	http.HandleFunc(path, t.httpServer)
}

func Serve(port int) {
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
