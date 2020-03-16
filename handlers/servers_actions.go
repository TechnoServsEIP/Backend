package handlers

import (
	"context"

	"gitlab.sysroot.ovh/technoservs/grpc-protos/go-game-servers"
)

func (g GameServer) StartServerAction(context.Context, *game_servers.StartServerActionRequest, *game_servers.StartServerActionResponse) error {
	panic("implement me")
}

func (g GameServer) StopServerAction(context.Context, *game_servers.StopServerActionRequest, *game_servers.StopServerActionResponse) error {
	panic("implement me")
}
