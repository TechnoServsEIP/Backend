package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.sysroot.ovh/technoservs/microservices/game-servers/app"
	"gitlab.sysroot.ovh/technoservs/microservices/game-servers/controllers"
	"gitlab.sysroot.ovh/technoservs/microservices/game-servers/utils"
	"log"
	"net/http"
	"os"
)

func main() {

	router := mux.NewRouter()

	port := os.Getenv("server_port") //Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println("listen on port", port)

	//TODO load database + pass to app struct

	router.HandleFunc("/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/user/update", controllers.UpdateAccount).Methods("POST")
	router.HandleFunc("/user/confirm", controllers.Confirm).Methods("POST")
	router.HandleFunc("/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/docker/create", controllers.CreateDocker).Methods("POST")
	router.HandleFunc("/docker/start", controllers.StartDocker).Methods("POST")
	router.HandleFunc("/docker/delete", controllers.StopDocker).Methods("GET")
	router.HandleFunc("/admin/Test", func(writer http.ResponseWriter, request *http.Request) {
		resp := map[string]interface{}{"message": "coucou"}
		utils.Respond(writer, resp, 200)
	}).Methods("GET")

	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	log.Fatal(http.ListenAndServe(":"+port, router))
}
