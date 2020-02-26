package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"oauth2server/app"
	"oauth2server/controllers"
	"os"
)

func main() {

	router := mux.NewRouter()
	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	port := os.Getenv("server_port") //Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	//TODO load database + pass to app struct
	
	router.HandleFunc("/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/docker/create", controllers.CreateDocker).Methods("POST")
	router.HandleFunc("/docker/start", controllers.StartDocker).Methods("POST")
	router.HandleFunc("/docker/delete", controllers.StopDocker).Methods("GET")

	log.Fatal(http.ListenAndServe(":" + port, router))
}
