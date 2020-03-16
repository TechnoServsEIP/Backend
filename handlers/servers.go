package handlers

import (
	"context"
	"fmt"

	"github.com/micro/go-micro/metadata"

	"gitlab.sysroot.ovh/technoservs/grpc-protos/go-game-servers"
)

func (g GameServer) ListServer(ctx context.Context, req *game_servers.ListServerRequest, res *game_servers.ListServerResponse) error {
	md, _ := metadata.FromContext(ctx)

	// local ip of service
	fmt.Println("local ip is", md["Local"])

	// remote ip of caller
	fmt.Println("remote ip is", md["Remote"])

	return nil
}

func (g GameServer) GetServer(context.Context, *game_servers.GetServerRequest, *game_servers.GetServerResponse) error {
	panic("implement me")
}

func (g GameServer) CreateServer(context.Context, *game_servers.CreateServerRequest, *game_servers.CreateServerResponse) error {
	panic("implement me")
}

func (g GameServer) DeleteServer(context.Context, *game_servers.DeleteServerRequest, *game_servers.DeleteServerResponse) error {
	panic("implement me")
}

func (g GameServer) UpdateServer(context.Context, *game_servers.UpdateServerRequest, *game_servers.UpdateServerResponse) error {
	panic("implement me")
}
