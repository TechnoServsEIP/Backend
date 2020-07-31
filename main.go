package main

import (
	"fmt"
	"github.com/TechnoServsEIP/Backend/models"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/TechnoServsEIP/Backend/app"
	"github.com/TechnoServsEIP/Backend/controllers"
	"github.com/rs/cors"
)

func main() {
	models.Initialization()
	router := mux.NewRouter()

	port := os.Getenv("server_port") //Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println("listen on port", port)

	//TODO load database + pass to app struct
	router.HandleFunc("/", controllers.Home).Methods("GET")
	router.HandleFunc("/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/user/update", controllers.UpdateAccount).Methods("POST")
	router.HandleFunc("/user/confirm", controllers.Confirm).Methods("POST")
	router.HandleFunc("/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/user/currentUser", controllers.GetEmail).Methods("GET")
	router.HandleFunc("/user/", controllers.GetUsers).Methods("GET")
	router.HandleFunc("/user/activate", controllers.Activate).Methods("POST")
	router.HandleFunc("/user/deactivate", controllers.Deactivate).Methods("POST")
	router.HandleFunc("/user/forgotpassword", controllers.SendPasswordReset).Methods("POST")
	router.HandleFunc("/user/resetpassword", controllers.ChangePassword).Methods("POST")
	router.HandleFunc("/docker/create", controllers.CreateDocker).Methods("POST")
	router.HandleFunc("/docker/start", controllers.StartDocker).Methods("POST")
	router.HandleFunc("/docker/stop", controllers.StopDocker).Methods("GET")
	router.HandleFunc("/docker/delete", controllers.DeleteDocker).Methods("POST")
	router.HandleFunc("/docker/logs", controllers.GetServerLogs).Methods("POST")
	router.HandleFunc("/docker/list", controllers.ListUserServers).Methods("POST")
	router.HandleFunc("/docker/infos", controllers.GetInfosUserServer).Methods("POST")
	router.HandleFunc("/docker/playersonline", controllers.GetPlayersOnline).Methods("GET")
	router.HandleFunc("/docker/update", controllers.ModifyGameServer).Methods("POST")
	router.HandleFunc("/offers/list", controllers.ListOffers).Methods("GET")
	router.HandleFunc("/offers/", controllers.GetOffer).Methods("POST")
	router.HandleFunc("/offers/create", controllers.CreateOffer).Methods("POST")
	router.HandleFunc("/offers/update", controllers.UpdateOffer).Methods("POST")
	router.HandleFunc("/offers/delete", controllers.DeleteOffer).Methods("POST")

	//router.Use(app.Cors)
	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowCredentials: true,
		AllowedHeaders: []string{ "Accept", "Accept-Encoding", "Authorization", "Content-Type", "Content-Length", "X-CSRF-Token" },
		Debug: true,
	})

	handler := c.Handler(router)

	log.Fatal(http.ListenAndServeTLS(":"+port, "/go/src/app/certs/fullchain.pem", "/go/src/app/certs/privkey.pem", handler))
}
