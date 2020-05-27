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

	fmt.Println("listen on port", port)

	//TODO load database + pass to app struct

	router.HandleFunc("/user/new", controllers.CreateAccount).Methods("POST, OPTIONS")
	router.HandleFunc("/user/update", controllers.UpdateAccount).Methods("POST, OPTIONS")
	router.HandleFunc("/user/confirm", controllers.Confirm).Methods("POST, OPTIONS")
	router.HandleFunc("/user/login", controllers.Authenticate).Methods("POST, OPTIONS")
	router.HandleFunc("/user/", controllers.GetUsers).Methods("GET, OPTIONS")
	router.HandleFunc("/user/activate", controllers.Activate).Methods("POST, OPTIONS")
	router.HandleFunc("/user/deactivate", controllers.Deactivate).Methods("POST, OPTIONS")
	router.HandleFunc("/docker/create", controllers.CreateDocker).Methods("POST, OPTIONS")
	router.HandleFunc("/docker/start", controllers.StartDocker).Methods("POST, OPTIONS")
	router.HandleFunc("/docker/stop", controllers.StopDocker).Methods("GET, OPTIONS")
	router.HandleFunc("/docker/delete", controllers.DeleteDocker).Methods("POST, OPTIONS")
	router.HandleFunc("/docker/logs", controllers.GetServerLogs).Methods("POST, OPTIONS")
	router.HandleFunc("/docker/list", controllers.ListUserServers).Methods("POST, OPTIONS")
	router.HandleFunc("/offers/list", controllers.ListOffers).Methods("GET, OPTIONS")
	router.HandleFunc("/offers/", controllers.GetOffer).Methods("POST, OPTIONS")
	router.HandleFunc("/offers/create", controllers.CreateOffer).Methods("POST, OPTIONS")
	router.HandleFunc("/offers/update", controllers.UpdateOffer).Methods("POST")
	router.HandleFunc("/offers/delete", controllers.DeleteOffer).Methods("POST, OPTIONS")

	router.Use(app.Cors)
	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	log.Fatal(http.ListenAndServe(":"+port, router))
}
