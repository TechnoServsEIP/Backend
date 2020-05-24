package models

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"

	//"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
	"context"
)

var mongoDB *mongo.Database //database
var ctx = context.Background()

func initMongoDb() {
	dbName := os.Getenv("mongodb_name")
	dbHost := os.Getenv("mongodb_host")
	dbPort := os.Getenv("mongodb_port")
	// Connect to the mongo database
	mongoCtx, _ := context.WithTimeout(ctx, 10*time.Second)
	mongoClient, err := mongo.Connect(mongoCtx, options.Client().ApplyURI("mongodb://"+dbHost+":"+dbPort+"/"+dbName))
	if err != nil {
		fmt.Print("fail to connect to the mongo database")
	}

	// Check if the mongo database is pinging
	//mongoCtx, _ = context.WithTimeout(ctx, 2*time.Second)
	//err = mongoClient.Ping(mongoCtx, readpref.Primary())
	//if err != nil {
	//	fmt.Print("fail to ping the mongo database")
	//}

	mongoDB = mongoClient.Database("technoservs-billing")
}

//returns a handle to the DB object
func GetMongoDB() *mongo.Database {
	return mongoDB
}

func GetCtx() context.Context {
	return ctx
}
