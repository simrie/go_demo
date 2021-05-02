package server

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hola")
}

/*
StartRouter defines the endpoints that use the db_pool for database connections
*/
func StartRouter() {
	http.HandleFunc("/hola", hello)
	http.ListenAndServe(":8080", nil)
}
