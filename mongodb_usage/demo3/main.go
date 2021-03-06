package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// TimePoint Task execution time point
type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

// LogRecord A log
type LogRecord struct {
	JobName   string    `bson:"jobName"`   // job name
	Command   string    `bson:"command"`   // shell command
	Err       string    `bson:"err"`       // script error
	Content   string    `bson:"content"`   // script output
	TimePoint TimePoint `bson:"timePoint"` // execution time point
}

func main() {
	var (
		client     *mongo.Client
		clientOps  *options.ClientOptions
		err        error
		collection *mongo.Collection
		record     *LogRecord
		logArr     []interface{} // just like C void* | Java Object
		result     *mongo.InsertManyResult
		insertId   interface{} // objectId
		docId      primitive.ObjectID
	)
	clientOps = options.Client().ApplyURI("mongodb://localhost:27017")
	if client, err = mongo.Connect(context.TODO(), clientOps); err != nil {
		fmt.Println(err)
		return
	}
	collection = client.Database("cron").Collection("log")

	// Insert data
	record = &LogRecord{
		JobName: "job10",
		Command: "echo hello",
		Err:     "",
		Content: "hello",
		TimePoint: TimePoint{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix() + 10,
		},
	}

	logArr = []interface{}{
		record,
		record,
		record,
	}

	if result, err = collection.InsertMany(context.TODO(), logArr); err != nil {
		fmt.Println(err)
		return
	}

	// snowflake

	for _, insertId = range result.InsertedIDs {
		docId = insertId.(primitive.ObjectID)
		fmt.Println("autoincrement id:", docId.Hex())
	}
}
