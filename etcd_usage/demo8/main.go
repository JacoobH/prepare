package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
		kv     clientv3.KV
		putOp  clientv3.Op
		getOp  clientv3.Op
		opResp clientv3.OpResponse
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

	//create op:operation
	putOp = clientv3.OpPut("/cron/jobs/job8", "job8")

	//execute op
	if opResp, err = kv.Do(context.TODO(), putOp); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Write Revision:", opResp.Put().Header.Revision)

	getOp = clientv3.OpGet("/cron/jobs/job8")

	if opResp, err = kv.Do(context.TODO(), getOp); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Data Revision:", opResp.Get().Kvs[0].ModRevision)
	fmt.Println("Data value:", string(opResp.Get().Kvs[0].Value))
}
