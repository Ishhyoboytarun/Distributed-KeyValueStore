package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

// store to hold key-value pairs
var store map[string]string

func initializeKVS() {
	// initialize the store
	store = make(map[string]string)

	// initialize the client to communicate with etcd
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:8080"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer etcdClient.Close()

	// create a lease for this server's ID
	resp, err := etcdClient.Grant(context.Background(), 10)
	if err != nil {
		log.Fatal(err)
	}
	leaseID := resp.ID

	// register this server's ID with etcd
	_, err = etcdClient.Put(context.Background(), fmt.Sprintf("servers/%d", leaseID), "", clientv3.WithLease(leaseID))
	if err != nil {
		log.Fatal(err)
	}

	// start the heartbeat loop
	go heartbeat(etcdClient, leaseID)
}

// main function - starts the HTTP server and registers handlers for different routes
func main() {

	initializeKVS()

	// initialize the router and register handlers for different routes
	router := mux.NewRouter()
	router.HandleFunc("/set/{key}/{value}", handleSet).Methods("POST")
	router.HandleFunc("/get/{key}", handleGet).Methods("GET")
	router.HandleFunc("/delete/{key}", handleDelete).Methods("DELETE")

	// start the HTTP server
	fmt.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
