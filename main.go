package main

import (
	"fmt"

	"github.com/simrie/go_demo.git/server"
	"github.com/simrie/go_demo.git/store"
)

func main() {
	fmt.Println("Go Demo with net/http server and graceful shutdown")

	// initialize empty itemStore
	itemStore := store.InitializeStore()
	server.StartRouter(itemStore)
}
