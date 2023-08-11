package main

import (
	"flag"
	"fmt"
	"go-trade-pnl/server"
	"net/http"
)

func main() {

	historyFilePathPtr := flag.String("trades", "", "TOS trade history")
	flag.Parse()

	updateStream := make(chan any, 1)

	server.AddContentUpdater(&server.TradingJournalUpdater{Directory: *historyFilePathPtr, Server: &server.TradeServer{}}, updateStream)

	srv := &http.Server{Addr: fmt.Sprintf("localhost:%d", 37777)}
	srv.ListenAndServe()
}
