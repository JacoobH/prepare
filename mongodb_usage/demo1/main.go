package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	var (
		client     *mongo.Client
		clientOps  *options.ClientOptions
		err        error
		db         *mongo.Database
		collection *mongo.Collection
	)
	clientOps = options.Client().ApplyURI("mongodb://localhost:27017")
	//1.establish connection
	if client, err = mongo.Connect(context.TODO(), clientOps); err != nil {
		fmt.Println(err)
		return
	}
	//2.select database my_db
	db = client.Database("my_db")
	//3.select table my_collection
	collection = db.Collection("my_collection")

	collection = collection
}
