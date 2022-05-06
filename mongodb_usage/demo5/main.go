package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// TimeBeforeCond {"$lt":timestamp}
type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

// DeleteCond {"timePoint.startTime":{"$lt":timestamp}}
type DeleteCond struct {
	BeforeCond TimeBeforeCond `bson:"timePoint.startTime"`
}

func main() {

	var (
		client     *mongo.Client
		clientOps  *options.ClientOptions
		err        error
		db         *mongo.Database
		collection *mongo.Collection
		delCond    *DeleteCond
		delResult  *mongo.DeleteResult
	)
	clientOps = options.Client().ApplyURI("mongodb://localhost:27017")
	//1.establish connection
	if client, err = mongo.Connect(context.TODO(), clientOps); err != nil {
		fmt.Println(err)
		return
	}
	//2.select database my_db
	db = client.Database("cron")
	//3.select table my_collection
	collection = db.Collection("log")

	delCond = &DeleteCond{
		BeforeCond: TimeBeforeCond{
			Before: time.Now().Unix(),
		},
	}

	if delResult, err = collection.DeleteMany(context.TODO(), delCond); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Number of rows deleted:", delResult.DeletedCount)
}
