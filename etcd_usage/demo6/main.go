package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	var (
		config         clientv3.Config
		client         *clientv3.Client
		err            error
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
		keepResp       *clientv3.LeaseKeepAliveResponse
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		putResp        *clientv3.PutResponse
		getResp        *clientv3.GetResponse
		kv             clientv3.KV
	)
	config = clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	//Apply for a lease
	lease = clientv3.NewLease(client)

	//Apply for a 10-second lease
	if leaseGrantResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}

	//Get id of lease
	leaseId = leaseGrantResp.ID

	//ctx, _ := context.WithTimeout(context.TODO(), 5 * time.Second)

	if keepRespChan, err = lease.KeepAlive(context.TODO(), leaseId); err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					fmt.Println("lease expired")
					goto END
				} else { // 每秒会续租一次, 所以就会受到一次应答
					fmt.Println("lease keep alive:", keepResp.ID)
				}
			}
		}
	END:
	}()

	kv = clientv3.NewKV(client)

	//Put a KV and associate it with lease to expire in 10 seconds
	if putResp, err = kv.Put(context.TODO(), "/cron/lock/job1", "", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Write successes", putResp.Header.Revision)

	//Check periodically to see if KV has expired
	for {
		if getResp, err = kv.Get(context.TODO(), "/cron/lock/job1"); err != nil {
			fmt.Println(err)
			return
		}
		if getResp.Count == 0 {
			fmt.Println("kv expired")
		}
		fmt.Println("Haven't expired:", getResp.Kvs)

		time.Sleep(2 * time.Second)
	}
}
