package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"gitlab.sysroot.ovh/technoservs/microservices/game-servers/models"
	"gitlab.sysroot.ovh/technoservs/microservices/game-servers/utils"
	"net/http"
)

var CreateDocker = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint) //Grab the id of the user that send the request
	fmt.Println("user: (", user, ")")

	docker := &models.Docker{}

	err := json.NewDecoder(r.Body).Decode(docker)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("error when creating docker client", err)
		return
	}

	port:= utils.GetPort()
	fmt.Println("port" + port)
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: port,
	}
	containerPort, err := nat.NewPort("tcp", port)
	if err != nil {
		fmt.Println("error when creating container port", err)
		return
	}
	portBinding := nat.PortMap{
		containerPort: []nat.PortBinding{hostBinding},
	}
	
	contName := "technoservers_test_" + docker.Game + "_" + utils.GenerateRandomString(6)
	fmt.Println("containeur name: " + contName)
	cont, err := cli.ContainerCreate(
		ctx,
		&container.Config{Image: "docker.io/itzg/minecraft-server", Env: []string{"EULA=TRUE"}},
		&container.HostConfig{PortBindings: portBinding},
		nil,
		contName)
	if err != nil {
		fmt.Println("error when creating container", err)
		return
	}

	err = cli.ContainerStart(ctx, cont.ID, types.ContainerStartOptions{})
	if err != nil {
		fmt.Println("an error occurred when starting container", "container_id=", cont.ID, "err", err)
		return
	}
	fmt.Println("Starting container ", cont.ID)

	dockerStore := &models.DockerStore{
		Game:   docker.Game,
		Id:     cont.ID,
		UserId: user,
	}
	resp := dockerStore.Create()
	utils.Respond(w, resp, http.StatusCreated)
}

var StartDocker = func(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	userId := r.Context().Value("user").(uint) //Grab the id of the user that send the request
	fmt.Println("userId: (", userId, ")")

	docker := &models.DockerDelete{}

	err := json.NewDecoder(r.Body).Decode(docker)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	err = cli.ContainerStart(ctx, docker.ContainerId, types.ContainerStartOptions{})
	if err != nil {
		fmt.Println("an error occurred when starting container", "container_id=", docker.ContainerId, "err", err)
		return
	}
	fmt.Println("Starting container ", docker.ContainerId)
}

var StopDocker = func(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	userId := r.Context().Value("user").(uint) //Grab the id of the user that send the request
	fmt.Println("userId: (", userId, ")")

	docker := &models.DockerDelete{}

	err := json.NewDecoder(r.Body).Decode(docker)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
	fmt.Println("Stop container " + docker.ContainerId)

	err = cli.ContainerStop(ctx, docker.ContainerId, nil)
	if err != nil {
		fmt.Println("An error occurred when stopping container ", docker.ContainerId)
		return
	}
	dockerStore := &models.DockerStore{
		Id:     docker.ContainerId,
		UserId: userId,
	}
	dockerStore.Update()
	utils.Respond(w,  map[string]interface{} {"status" : 200, "message" : "Container Stop successfully"}, http.StatusOK)
}