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

	"github.com/simrie/go_demo.git/store"
)

/*
GetStatsHandler returns the Stats on the items in the Store
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
GetHashHandler retrieves a hashed value of a string
and then stores it as an Item in the Store
*/
func GetHashHandler(itemStore *store.Store, response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "text/plain")
	var parentCtx = request.Context()

	_, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()

	// Only process Post requests
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
		response.Write([]byte(`"Error": Request data not readable "`))
		return
	}
	strBody := string(b)
	if strBody == "" {
		response.WriteHeader(http.StatusBadRequest)
		log.Printf("\nGetHashHandler %v: ", err.Error())
		response.Write([]byte(`"Error": Data is not a string "`))
		return
	}
	args := strings.Split(strBody, "assword=")
	if len(args) <= 0 {
		response.WriteHeader(http.StatusBadRequest)
		log.Printf("\nGetHashHandler %v: ", err.Error())
		response.Write([]byte(`"Error": Missing Password in Data "`))
		return
	}
	pwd = args[1]

	// Retrieve a created Item that contains the hashed value of pwd
	item, err := store.CreateItem(pwd)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		log.Printf("\nGetHashHandler %v: ", err.Error())
		response.Write([]byte(`"Error": Hashing Error "`))
		return
	}

	// Store the new Item
	var order int32
	order, err = itemStore.StoreItem(item)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("\nGetHashHandler %v: ", err.Error())
		response.Write([]byte(`"Error": Storing error "`))
		return
	}

	response.WriteHeader(http.StatusOK)
	strOrder := strconv.Itoa(int(order))
	response.Write([]byte(strOrder))
}

/*
HandlerDefault returns when an unknown endpoint is called
unless the url has a valid integer id as a hash arg: /hash/:id
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
				response.Write([]byte(`"Error": Invalid url format "`))
				return
			}
			int32id := int32(id)
			resp, err = itemStore.GetItemById(int32id)
			if err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				log.Printf(`\nHandlerDefault testId %s %v`, testId, err.Error())
				response.Write([]byte(`"Error": Unable to retrieve by id "`))
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
