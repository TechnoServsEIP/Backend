package main

import (
<<<<<<< HEAD
	"context"
	"time"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/service/grpc"
	microlog "github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/broker/nsq"
	"github.com/micro/go-plugins/registry/consul"

	"gitlab.sysroot.ovh/technoservs/grpc-protos/go-game-servers"
	"gitlab.sysroot.ovh/technoservs/internal-libs/go-utils/logger"
	"gitlab.sysroot.ovh/technoservs/microservices/game-servers/handlers"
)

func main() {
	log := logger.Default()
	ctx := logger.ToCtx(context.Background(), log)
	microlog.SetLogger(logger.ToGoMicro(log))

	serviceName := "fr.technoservs.srv.game-servers"

	service := grpc.NewService(
		micro.Name(serviceName),
		micro.Context(ctx),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		micro.Registry(
			consul.NewRegistry(
				registry.Addrs("172.17.0.1:8500"),
			),
		),
		micro.Broker(
			nsq.NewBroker(
				broker.Addrs("172.17.0.1:4150"),
			),
		),
	)

	service.Init()

	// Register handlers
	game_servers.RegisterGameServerServiceHandler(service.Server(), new(handlers.GameServer))

	if err := service.Run(); err != nil {
		log.WithError(err).Fatal("fail to start service")
	}
=======
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
	router.HandleFunc("/docker/update", controllers.ModifyGameServer).Methods("POST")
	router.HandleFunc("/offers/list", controllers.ListOffers).Methods("GET")
	router.HandleFunc("/offers/", controllers.GetOffer).Methods("POST")
	router.HandleFunc("/offers/create", controllers.CreateOffer).Methods("POST")
	router.HandleFunc("/offers/update", controllers.UpdateOffer).Methods("POST")
	router.HandleFunc("/offers/delete", controllers.DeleteOffer).Methods("POST")

	router.Use(app.Cors)
	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	log.Fatal(http.ListenAndServe(":"+port, router))
>>>>>>> clientGRPCBilling
}
