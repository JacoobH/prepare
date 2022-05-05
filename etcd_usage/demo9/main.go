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
		ctx            context.Context
		cancelFunc     context.CancelFunc
		kv             clientv3.KV
		txn            clientv3.Txn
		txnResp        *clientv3.TxnResponse
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

	ctx, cancelFunc = context.WithCancel(context.TODO())

	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseId)

	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
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

	// IF not exist key, THEN set it, ELSE fail
	kv = clientv3.NewKV(client)
	// create transaction
	txn = kv.Txn(context.TODO())
	// Define transaction
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/job9"), "=", 0)).
		Then(clientv3.OpPut("/cron/lock/job9", "xxx", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/lock/job9"))
	// Commit transaction
	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println(err)
		return
	}

	// Determine if you grabbed the lock
	if !txnResp.Succeeded {
		fmt.Println("The lock is occupied", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
	}
}
