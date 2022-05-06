package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// FindByJobName jobName filter
type FindByJobName struct {
	JobName string `bson:"jobName"`
}

func main() {
	var (
		client     *mongo.Client
		clientOps  *options.ClientOptions
		err        error
		collection *mongo.Collection
		cond       *FindByJobName
		findOpt    *options.FindOptions
		record     *LogRecord
		cur        *mongo.Cursor
	)
	clientOps = options.Client().ApplyURI("mongodb://localhost:27017")
	if client, err = mongo.Connect(context.TODO(), clientOps); err != nil {
		fmt.Println(err)
		return
	}
	collection = client.Database("cron").Collection("log")
	cond = &FindByJobName{
		JobName: "job10",
	}
	findOpt = options.Find()
	findOpt.SetSkip(0)
	findOpt.SetLimit(2)

	if cur, err = collection.Find(context.TODO(), cond, findOpt); err != nil {
		fmt.Println(err)
		return
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		record = &LogRecord{}
		if err = cur.Decode(record); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(*record)
	}
}
