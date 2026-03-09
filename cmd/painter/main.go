package main

import (
	"log"
	"net/http"
	"os"

	"github.com/flaboy/painter/internal/httpserver"
)

func main() {
	addr := ":" + portFromEnv()
	if err := http.ListenAndServe(addr, httpserver.NewHandler()); err != nil {
		log.Fatal(err)
	}
}

func portFromEnv() string {
	if v := os.Getenv("PAINTER_PORT"); v != "" {
		return v
	}
	return "7013"
}
