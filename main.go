package main

import (
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
}
