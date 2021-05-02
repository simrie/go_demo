package main

import (
	"fmt"

	"github.com/simrie/go_demo.git/server"
)

func main() {
	fmt.Println("How ya' doin'?")

	server.StartRouter()
}
