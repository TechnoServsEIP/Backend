package tests

import (
	"testing"
	"fmt"
	"time"
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/TechnoServsEIP/Backend/utils"
	"github.com/docker/go-connections/nat"
	"github.com/TechnoServsEIP/Backend/controllers"
)

func CreateServer() string {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("error when creating docker client")
		return ""
	}

	port := utils.GetPort()
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: port,
	}
	containerPort, err := nat.NewPort("tcp", "25565")
	if err != nil {
		fmt.Println("error when creating container port")
		return ""
	}
	portBinding := nat.PortMap{
		containerPort: []nat.PortBinding{hostBinding},
	}

	contName := "technoservers_unitTest_" + utils.GenerateRandomString(6)
	cont, err := cli.ContainerCreate(
		ctx,
		&container.Config{Image: "docker.io/itzg/minecraft-server", Env: []string{"EULA=TRUE"}},
		&container.HostConfig{PortBindings: portBinding},
		nil,
		contName)
	if err != nil {
		fmt.Println("error when creating container")
		return ""
	}

	err = cli.ContainerStart(ctx, cont.ID, types.ContainerStartOptions{})
	if err != nil {
		fmt.Println("an error occurred when starting container")
		return ""
	}

	return cont.ID
}

func StopServer(contID string) {
	ctx := context.Background()

	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("An error occurred when stopping container")
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	fmt.Print(containers)
	if err != nil {
		fmt.Println("An error occurred when stopping container")
	}

	err = cli.ContainerStop(ctx, contID, nil)
	if err != nil {
		fmt.Println("An error occurred when stopping container")
		return
	}
}

func DelServer(contID string) {
	ctx := context.Background()
	
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Error failed to contact docker api")
		return
	}

	err = cli.ContainerRemove(ctx, contID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})

	if err != nil {
		fmt.Println("An error appear when removing container")
	}
}

func TestPlayersOnline(t *testing.T) {
	contID := CreateServer()

	if contID == "" {
		t.Errorf("Error when creating container")
	}

	time.Sleep(120 * time.Second)

	playersOnline := controllers.GetNumberPlayers(contID)


	if playersOnline["connectedPlayers"] != 0 {
		t.Errorf("Error number of connected players")
	}
	if playersOnline["maxPlayers"] != 20 {
		t.Errorf("Error number of max players")
	}
	if len(playersOnline["listPlayers"].([]string)) != 0 {
		t.Errorf("Error list players")
	}

	StopServer(contID)
	DelServer(contID)
}