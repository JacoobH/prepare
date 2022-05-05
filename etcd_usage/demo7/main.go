package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	var (
		config             clientv3.Config
		client             *clientv3.Client
		err                error
		kv                 clientv3.KV
		getResp            *clientv3.GetResponse
		watchStartRevision int64
		watcher            clientv3.Watcher
		watchChan          <-chan clientv3.WatchResponse
		watchResp          clientv3.WatchResponse
		event              *clientv3.Event
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

	//Simulate the change of KV in ETCD
	go func() {
		for {
			kv.Put(context.TODO(), "/cron/jobs/job7", "I am job7")
			kv.Delete(context.TODO(), "/cron/jobs/job7")
			time.Sleep(1 * time.Second)
		}
	}()

	//Get the current value and observe subsequent changes
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job7"); err != nil {
		fmt.Println(err)
		return
	}

	if len(getResp.Kvs) != 0 {
		fmt.Println("Current value:", string(getResp.Kvs[0].Value))
	}

	//Etcd Cluster transaction ID, monotonically increasing
	watchStartRevision = getResp.Header.Revision

	//Create a watcher
	watcher = clientv3.NewWatcher(client)

	//Start watching
	fmt.Println("Start listening from this version:", watchStartRevision)
	watchChan = watcher.Watch(context.TODO(), "/cron/jobs/job7", clientv3.WithRev(watchStartRevision))

	//Handle event about the changes of kv
	for watchResp = range watchChan {
		for _, event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("Modified to:", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision)
			case mvccpb.DELETE:
				fmt.Println("Deleted", "Revision:", event.Kv.ModRevision)

			}
		}
	}
}
