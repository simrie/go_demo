package server

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/simrie/go_demo.git/hasher"
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
	response.Header().Set("content-type", "text/plain")
	var parentCtx = request.Context()

	_, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()

	// Fall through if request Method is unknown

	if request.Method != http.MethodPost {
		response.WriteHeader(http.StatusBadRequest)
		log.Printf("\nGetHashHandler : Unknown Method")
		response.Write([]byte(`"Error": Unknown Method "`))
		return
	}

	var pwd string
	b, err := ioutil.ReadAll(request.Body)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("\nGetHashHandler %v: ", err.Error())
		response.Write([]byte(`"Error": Hasher error "`))
		return
	}
	strBody := string(b)
	if strBody == "" {
		response.WriteHeader(http.StatusBadRequest)
		log.Printf("\nGetHashHandler %v: ", err.Error())
		response.Write([]byte(`"Error": Missing Password "`))
		return
	}
	args := strings.Split(strBody, "assword=")
	if len(args) <= 0 {
		response.WriteHeader(http.StatusBadRequest)
		log.Printf("\nGetHashHandler %v: ", err.Error())
		response.Write([]byte(`"Error": Missing Password "`))
		return
	}

	// encrypt and store the value
	var order int32
	hasher := hasher.Encode
	order, err = hasher(itemStore, pwd)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("\nGetHashHandler %v: ", err.Error())
		response.Write([]byte(`"Error": Hasher error "`))
		return
	}
	response.WriteHeader(http.StatusOK)
	strOrder := strconv.Itoa(int(order))
	response.Write([]byte(strOrder))
}

/*
HandlerDefault returns when an unknown endpoint is called
but is also used to extract and id from /hash/:id
since net/http doesn't seem to handle URL args
*/
func HandlerDefault(itemStore *store.Store, response http.ResponseWriter, request *http.Request) {
	// handle /hash/:id situation if :id is an integer
	if request.Method == http.MethodGet {
		var resp store.Item
		args := strings.Split(request.RequestURI, `/hash/`)
		if len(args) > 0 {
			testId := args[1]
			id, err := strconv.Atoi(testId)
			if err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				log.Printf(`\nHandlerDefault testId %s %v`, testId, err.Error())
				response.Write([]byte(err.Error()))
				return
			}
			int32id := int32(id)
			resp, err = itemStore.GetItem(int32id)
			if err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				log.Printf(`\nHandlerDefault testId %s %v`, testId, err.Error())
				response.Write([]byte(err.Error()))
				return
			}
			response.WriteHeader(http.StatusOK)
			json.NewEncoder(response).Encode(resp.Value)
			return
		}
	}

	// otherwise treat this as an invalid request
	log.Printf("HandlerDefault %s\n", request.RequestURI)
	response.Write([]byte("Invalid Request\n"))
}
