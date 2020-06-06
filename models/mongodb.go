package models

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	//"go.mongodb.org/mongo-driver/mongo/readpref"
	"context"
	"time"
)

var mongoDB *mongo.Database //database
var ctx = context.Background()

func initMongoDb() {
	err := godotenv.Load() //Load .env file
	if err != nil {
		fmt.Print("error when opening .env file", err)
	}

	dbHost := os.Getenv("mongodb_host")
	dbPort := os.Getenv("mongodb_port")
	dbName := os.Getenv("mongodb_name")

	dbURI := fmt.Sprintf("mongodb://%s:%s/%s", dbHost, dbPort, dbName)

	// Connect to the mongo database
	mongoCtx, _ := context.WithTimeout(ctx, 10*time.Second)
	mongoClient, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(dbURI))
	if err != nil {
		fmt.Print("fail to connect to the mongo database")
	}

	mongoDB = mongoClient.Database("technoservs-billing")
}

//returns a handle to the DB object
func GetMongoDB() *mongo.Database {
	return mongoDB
}

func GetCtx() context.Context {
	return ctx
}
