package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// handles requests to set a key-value pair
func handleSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value := vars["value"]

	// add key-value pair to the store
	store[key] = value

	// notify other servers in the cluster
	err := notifyOtherServers(fmt.Sprintf("set,%s,%s", key, value))
	if err != nil {
		fmt.Println(err)
	}

	// send response to client
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Key-value pair set successfully\n")
}

// handles requests to get the value for a given key
func handleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	// retrieve value for the given key
	value, ok := store[key]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Key not found\n")
		return
	}

	// send response to client
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Value: %s\n", value)
}

// handles requests to delete a key-value pair
func handleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	// delete key-value pair from the store
	delete(store, key)

	// notify other servers in the cluster
	err := notifyOtherServers(fmt.Sprintf("delete,%s", key))
	if err != nil {
		fmt.Println(err)
	}

	// send response to client
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Key-value pair deleted successfully\n")
}
