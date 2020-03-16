package handlers

import (
	"context"

	"gitlab.sysroot.ovh/technoservs/grpc-protos/go-game-servers"
)

func (g GameServer) ListServerEnvironment(context.Context, *game_servers.ListServerEnvironmentRequest, *game_servers.ListServerEnvironmentResponse) error {
	panic("implement me")
}
