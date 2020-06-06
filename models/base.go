package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db *gorm.DB //database

func init() {
	initMongoDb()

	err := godotenv.Load() //Load .env file
	if err != nil {
		fmt.Print("error when opening .env file", err)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Build connection string

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		panic(err.Error())
	}

	db = conn
	db.Debug().AutoMigrate(&Account{}, &DockerStore{}, DockerDelete{}) //Database migration
}

//returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}
