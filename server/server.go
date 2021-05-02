package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/simrie/go_demo.git/store"
)

/*
StartRouter defines the endpoints that use the db_pool for database connections
*/
func StartRouter(itemStore *store.Store) {
	// enable graceful shutdown per http documentation

	var srv http.Server
	srv.Addr = ":8080"
	// Overriding some limits and timeouts
	srv.ReadTimeout = 10 * time.Second
	srv.WriteTimeout = 10 * time.Second
	srv.MaxHeaderBytes = 1 << 20

	router := http.NewServeMux()
	router.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		GetStatsHandler(itemStore, w, r)
	})
	router.HandleFunc("/hash", func(w http.ResponseWriter, r *http.Request) {
		GetHashHandler(itemStore, w, r)
	})
	router.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Shutting Down"))
		ctx, cancel := context.WithCancel(context.Background())
		log.Printf("shutting down server")
		cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	})
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HandlerDefault(itemStore, w, r)
	})

	srv.Handler = router

	log.Printf("starting server")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

}
