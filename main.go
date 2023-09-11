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

	server.AddContentUpdater(&server.TradingJournalUpdater{Directory: *historyFilePathPtr, Server: &server.TradeServer{}})

	fmt.Printf("Starting TradeWiz Server\n")

	srv := &http.Server{Addr: fmt.Sprintf(":%d", 37777)}
	srv.ListenAndServe()
}
