package main

import (
	"github.com/k-orolevsk-y/go-metricts-tpl/cmd/server/handlers"
	"github.com/k-orolevsk-y/go-metricts-tpl/cmd/server/storage"
	"net/http"
)

func main() {
	storage := stor.NewMem()

	mux := http.NewServeMux()
	mux.HandleFunc("/update/gauge/", handlers.UpdateGauge(&storage))
	mux.HandleFunc("/update/counter/", handlers.UpdateCounter(&storage))
	mux.HandleFunc("/", handlers.BadRequest)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
