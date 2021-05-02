package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func hola(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hola\n"))
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello\n"))
}

/*
StartRouter defines the endpoints that use the db_pool for database connections
*/
func StartRouter() {
	// enable graceful shutdown per http documentation

	var srv http.Server
	srv.Addr = ":8080"
	// Overriding some limits and timeouts
	srv.ReadTimeout = 10 * time.Second
	srv.WriteTimeout = 10 * time.Second
	srv.MaxHeaderBytes = 1 << 20

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("/hola", hola)
	mux.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		srv.Shutdown(context.Background())
	})

	srv.Handler = mux

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	<-idleConnsClosed

	log.Printf("server started")

}
