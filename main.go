package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/httpapi"
	"github.com/jb843051627/fjord-resonance/internal/service"
	"github.com/jb843051627/fjord-resonance/internal/sqlite"
)

func main() {
	databasePath := os.Getenv("FJORD_RESONANCE_DB")
	store, err := sqlite.Open(databasePath)
	if err != nil {
		panic(err)
	}
	defer store.Close()
	app := serviceApplication(store)
	defer app.Close()
	if len(os.Args) > 1 && os.Args[1] == "--smoke-test" {
		fmt.Println("fjord-resonance smoke test: ok")
		return
	}
	apiServer := httpapi.New(app)
	mux := http.NewServeMux()
	mux.Handle("/api/", apiServer.Handler())
	mux.Handle("/healthz", apiServer.Handler())
	mux.Handle("/", http.FileServer(http.Dir("web")))
	server := &http.Server{Addr: envOr("FJORD_RESONANCE_ADDR", ":8080"), Handler: mux, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func serviceApplication(store *sqlite.Store) *service.Application {
	return service.NewApplication(store)
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
