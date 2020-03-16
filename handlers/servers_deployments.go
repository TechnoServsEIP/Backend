package handlers

import (
	"context"

	"gitlab.sysroot.ovh/technoservs/grpc-protos/go-game-servers"
)

func (g GameServer) ListServerDeployment(context.Context, *game_servers.ListServerDeploymentRequest, *game_servers.ListServerDeploymentResponse) error {
	panic("implement me")
}

func (g GameServer) GetServerDeployment(context.Context, *game_servers.GetServerDeploymentRequest, *game_servers.GetServerDeploymentRequest) error {
	panic("implement me")
}
