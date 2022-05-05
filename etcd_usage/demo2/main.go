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
		putResp *clientv3.PutResponse
	)
	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	//Establish a client
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	//Use to read or write KV of etcd
	kv = clientv3.NewKV(client)

	if putResp, err = kv.Put(context.TODO(), "cron/jobs/job1", "hi", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Revision:", putResp.Header.Revision)
		if putResp.PrevKv != nil {
			fmt.Println("PrevValue:", string(putResp.PrevKv.Value))
		}
	}
}
