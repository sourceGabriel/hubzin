package main

import (
	"log"
	"net/http"

	"github.com/sourceGabriel/hubzin/backend"
)

func main() {
	srv := backend.NewServer()
	log.Println("DeskHub backend listening on :8080")
	if err := http.ListenAndServe(":8080", srv.Handler()); err != nil {
		log.Fatal(err)
	}
}
