package main

import (
	"fmt"

	"github.com/simrie/go_demo.git/server"
	"github.com/simrie/go_demo.git/store"
)

/*
Main initializes an empty Store and passes it to
the function that starts the routed server
so the router handler functions can access the itemStore
*/
func main() {
	fmt.Println("Go Demo with net/http server and graceful shutdown")

	// initialize empty itemStore
	itemStore := store.InitializeStore()
	server.StartRouter(itemStore)
}
