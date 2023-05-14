package main

import (
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
)

// notifyOtherServers sends a message to other servers in the cluster
// message format: "action,key,value" or "action,key"
func notifyOtherServers(msg string) error {
	// create a new etcd client
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	defer etcdClient.Close()

	// get the IDs of all servers in the cluster
	resp, err := etcdClient.Get(context.Background(), "servers/", clientv3.WithPrefix())
	if err != nil {
		return err
	}

	// send the message to each server in the cluster
	for _, kv := range resp.Kvs {
		id := string(kv.Key)[8:]
		_, err = etcdClient.Put(context.Background(), fmt.Sprintf("messages/%s", id), msg)
		if err != nil {
			return err
		}
	}

	return nil
}

// heartbeat sends periodic heartbeats to etcd to keep the lease alive
func heartbeat(etcdClient *clientv3.Client, leaseID clientv3.LeaseID) {
	for {
		// send a heartbeat
		_, err := etcdClient.KeepAliveOnce(context.Background(), leaseID)
		if err != nil {
			log.Println(err)
		}

		// wait for a short time before sending the next heartbeat
		time.Sleep(5 * time.Second)
	}
}
