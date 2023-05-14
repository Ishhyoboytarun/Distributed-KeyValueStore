package main

import (
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
)

type KVStore struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

func NewKVStore(endpoints []string) (*KVStore, error) {
	config := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	}
	client, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}
	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)
	return &KVStore{client, kv, lease}, nil
}

func (kvs *KVStore) Get(key string) (string, error) {
	resp, err := kvs.kv.Get(context.Background(), key)
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("Key not found: %s", key)
	}
	return string(resp.Kvs[0].Value), nil
}

func (kvs *KVStore) Put(key string, value string) error {
	_, err := kvs.kv.Put(context.Background(), key, value)
	return err
}

func (kvs *KVStore) Delete(key string) error {
	_, err := kvs.kv.Delete(context.Background(), key)
	return err
}
