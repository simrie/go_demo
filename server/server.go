package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/simrie/go_demo.git/store"
)

/*
StartRouter defines the exposed endpoints, some of which
make sure of itemStore as we are not using a persistent database
*/
func StartRouter(itemStore *store.Store) {

	ctx, cancel := context.WithCancel(context.Background())

	// net/http server
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
		//ctx, cancel := context.WithCancel(context.Background())
		log.Printf("shutting down server")
		cancel()
		if err := srv.Shutdown(ctx); err != nil {
			//log.Fatal(err)
			log.Print(err)
			defer os.Exit(0)
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

	// sigint shutdown
	// how to references:
	// https://rafallorenz.com/go/handle-signals-to-graceful-shutdown-http-server/
	// https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a

	// expecting more than one signal?
	signalChan := make(chan os.Signal)

	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
		syscall.SIGKILL,
		//syscall.CTRL_C_EVENT,
	)

	sig := <-signalChan
	log.Printf("os.Interrupt sig %v - shutting down...\n", sig)

	if sig != nil {
		cancel()
		defer os.Exit(0)
	}

}
