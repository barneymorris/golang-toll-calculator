package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/betelgeusexru/golang-toll-calculator/types"
)

func main() {
	listenAddr := flag.String("listenaddr", ":3000", "the listen address of the HTTP server")
	flag.Parse()
	
	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)

	makeHTTPTransport(*listenAddr, svc)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	fmt.Println("HTTP Transport running on port ", listenAddr)

	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.ListenAndServe(listenAddr, nil)
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}		
	}
}