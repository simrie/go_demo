package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/simrie/go_demo.git/store"
)

/*
GetStatsHandler returns the Stats on the items in the store
*/
func GetStatsHandler(itemStore *store.Store, response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var parentCtx = request.Context()

	_, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()

	stats, err := itemStore.GetStats()

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Println(`{ "GetStatsHandler": "` + err.Error() + `" }`)
		response.Write([]byte(`{ "message": Error Retrieving Stats" }`))
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(stats)
}

/*
GetHashHandler posts an item to or gets an item from the item store
*/
func GetHashHandler(itemStore *store.Store, response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var parentCtx = request.Context()
	var resp store.Item

	_, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()

	if request.Method == http.MethodGet {
		//request.RequestURI
		//request.URL
		log.Printf("%v", request.RequestURI)
	}

	if request.Method == http.MethodPost {
		reqBody, _ := request.GetBody()
		log.Printf("%v", reqBody)

	}

	resp = store.Item{Value: "test"}
	var err error
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Println(`{ "GetHashHandler": "` + err.Error() + `" }`)
		response.Write([]byte(`{ "message": Error " }`))
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(resp.Value)
}

/*
HandlerDefault returns when an unknown endpoint is called
*/
func HandlerDefault(response http.ResponseWriter, request *http.Request) {
	log.Printf("HandlerDefault %s\n", request.RequestURI)
	response.Write([]byte("Invalid\n"))
}
