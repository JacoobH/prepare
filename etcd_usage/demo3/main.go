package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		err     error
		kv      clientv3.KV
		getResp *clientv3.GetResponse
	)
	config = clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	kv = clientv3.NewKV(client)

	if getResp, err = kv.Get(context.TODO(), "name" /*clientv3.WithCountOnly()*/); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(getResp.Kvs, getResp.Count)
	}
}
