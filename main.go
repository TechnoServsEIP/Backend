package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gitlab.sysroot.ovh/technoservs/microservices/game-servers/app"
	"gitlab.sysroot.ovh/technoservs/microservices/game-servers/controllers"
)

func main() {

	router := mux.NewRouter()

	port := os.Getenv("server_port") //Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	//TODO load database + pass to app struct

	router.HandleFunc("/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/user/confirm", controllers.Confirm).Methods("POST")
	router.HandleFunc("/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/docker/create", controllers.CreateDocker).Methods("POST")
	router.HandleFunc("/docker/start", controllers.StartDocker).Methods("POST")
	router.HandleFunc("/docker/delete", controllers.StopDocker).Methods("GET")
	router.HandleFunc("/offers/list", controllers.ListOffers).Methods("GET")
	router.HandleFunc("/offers/", controllers.GetOffer).Methods("POST")
	router.HandleFunc("/offers/create", controllers.CreateOffer).Methods("POST")
	router.HandleFunc("/offers/update", controllers.UpdateOffer).Methods("POST")
	router.HandleFunc("/offers/delete", controllers.DeleteOffer).Methods("POST")

	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	log.Fatal(http.ListenAndServe(":"+port, router))
}
