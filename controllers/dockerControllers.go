package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"os/exec"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/TechnoServsEIP/Backend/models"
	"github.com/TechnoServsEIP/Backend/utils"
)

var GetAllPortBinded = func() []string {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	
	if err != nil {
		fmt.Printf("Impossible to use the docker client\n")
		return []string{}
	}

	containers := models.ListAllDockers()
	if err != nil {
		fmt.Printf("Error when try to retrieve containers ports\n")
		return []string{}
	}

	tmpPortsBinded := []string{}

	if containers != nil {
		for _, container := range *containers {
			info, err := cli.ContainerInspect(ctx, container.IdDocker)

			if err != nil {
				fmt.Println(err.Error())
				return []string{}
			}
			tmpPortsBinded = append(tmpPortsBinded, info.HostConfig.PortBindings["25565/tcp"][0].HostPort)
		}
	}

	return tmpPortsBinded
}

var CreateDocker = func(w http.ResponseWriter, r *http.Request) {
	//user := r.Context().Value("user").(uint) //Grab the id of the user that send the request
	//fmt.Println("user: (", user, ")")

	docker := &models.Docker{}

	err := json.NewDecoder(r.Body).Decode(docker)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	if docker.UserId == "" {
		utils.Respond(w, utils.Message(false, "User ID is empty"), http.StatusBadRequest)
		return
	}

	if !checkIfUserCanCreate(docker.UserId) {
		utils.Respond(w, utils.Message(false, "User already have max server reached"), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("error when creating docker client", err)
		return
	}

	port := utils.GetPort()

	if port == "no port available" {
		utils.Respond(w, utils.Message(false, "No port available"), 413)
		return
	}
	fmt.Println("port" + port)
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: port,
	}
	containerPort, err := nat.NewPort("tcp", "25565")
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

	u64, err := strconv.ParseUint(docker.UserId, 10, 32)

	dockerStore := &models.DockerStore{
		Game:         docker.Game,
		IdDocker:     cont.ID,
		UserId:       uint(u64),
		ServerName:   docker.ServerName,
		ServerStatus: "Started",
	}
	resp := dockerStore.Create()

	info, err := cli.ContainerInspect(ctx, cont.ID)

	if err != nil {
		fmt.Println(err)
		utils.Respond(w, resp, http.StatusCreated)
		return
	}

	resp["settings"] = &info

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

	info, err := cli.ContainerInspect(ctx, docker.ContainerId)

	if err != nil {
		fmt.Println(err)
		return
	}

	dockerStore := &models.DockerStore{
		IdDocker: docker.ContainerId,
		UserId:   userId,
	}

	dockerStore.UpdateServerStatus("Started")

	if err != nil {
		fmt.Println("Error while updating status server")
	}

	resp := map[string]interface{}{}

	resp["settings"] = info

	utils.Respond(w, resp, 200)
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
		IdDocker: docker.ContainerId,
		UserId:   userId,
	}

	dockerStore.UpdateServerStatus("Stoped")

	if err != nil {
		fmt.Println("Error while updating status server")
	}

	utils.Respond(w, map[string]interface{}{"status": 200, "message": "Container Stop successfully"}, http.StatusOK)
}

var GetServerLogs = func(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error failed to contact docker api"), http.StatusBadRequest)
		return
	}

	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Details:    true,
	}

	docker := &models.DockerDelete{}

	err = json.NewDecoder(r.Body).Decode(docker)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	out, err := cli.ContainerLogs(ctx, docker.ContainerId, options)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error bad container_id"), http.StatusBadRequest)
		return
	}

	var logs bytes.Buffer

	io.Copy(&logs, out)

	resp := map[string]interface{}{
		"logs": logs.String(),
	}

	utils.Respond(w, resp, 200)
}

var DeleteDocker = func(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(uint)
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error failed to contact docker api"), http.StatusBadRequest)
		return
	}

	docker := &models.DockerDelete{}

	err = json.NewDecoder(r.Body).Decode(docker)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	// retrieve the port to free
	info, err := cli.ContainerInspect(ctx, docker.ContainerId)

	if err != nil {
		fmt.Println(err.Error())
		utils.Respond(w, utils.Message(false, "Error while retrieving port of the container"), 500)
		return 
	}

	err = cli.ContainerRemove(ctx, docker.ContainerId, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	if err != nil {
		fmt.Println("An error appear when removing container: ", docker.ContainerId, "err ", err)
		utils.Respond(w, utils.Message(false, "Error while removing container"), http.StatusBadRequest)
		return
	}

	utils.FreeThePort(info.HostConfig.PortBindings["25565/tcp"][0].HostPort)

	resp := models.RemoveContainer(userId, docker.ContainerId)

	utils.Respond(w, resp, 204)
}

var ListUserServers = func(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	docker := &models.DockerList{}
	// userId := r.Context().Value("user").(uint)

	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("error when creating docker client", err)
		return
	}

	err = json.NewDecoder(r.Body).Decode(docker)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	if docker.UserId == "" {
		utils.Respond(w, utils.Message(false, "User ID is empty"), http.StatusBadRequest)
		return
	}

	u64, err := strconv.ParseUint(docker.UserId, 10, 32)

	allDocker := models.UserServers(uint(u64))

	if allDocker == nil {
		utils.Respond(w, utils.Message(false, "invalid user_id"), 500)
		return
	}

	list := make([]models.DockerStore, 0)

	for _, element := range *allDocker {
		info, err := cli.ContainerInspect(ctx, element.IdDocker)

		if err != nil {
			fmt.Println(err.Error())
			utils.Respond(w, map[string]interface{}{
				"error": err.Error(),
			}, 500)
			return
		}

		element.Settings = &info
		list = append(list, element)
	}

	resp := map[string]interface{}{
		"list": list,
	}

	utils.Respond(w, resp, 200)
}

var GetInfosUserServer = func(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	docker := &models.DockerDelete{}

	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("error when creating docker client", err)
		return
	}

	err = json.NewDecoder(r.Body).Decode(docker)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	u64, err := strconv.ParseUint(docker.UserId, 10, 32)

	OneDocker := models.OneUserServer(uint(u64), docker.ContainerId)

	if OneDocker == nil {
		utils.Respond(w, utils.Message(false, "invalid user_id or container_id"), 500)
		return
	}

	serverInfo := map[string]interface{}{}

	serverInfo["server_infos"] = OneDocker

	info, err := cli.ContainerInspect(ctx, OneDocker.IdDocker)

	if err != nil {
		fmt.Println(err.Error())
		utils.Respond(w, map[string]interface{}{
			"error": err.Error(),
		}, 500)
		return
	}

	serverInfo["settings"] = &info

	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Details:    true,
	}

	out, err := cli.ContainerLogs(ctx, docker.ContainerId, options)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error bad container_id"), http.StatusBadRequest)
		return
	}

	var logs bytes.Buffer

	io.Copy(&logs, out)
	serverInfo["logs"] = logs.String()

	serverInfo["playersOnline"] = GetNumberPlayers(docker.ContainerId)

	utils.Respond(w, serverInfo, 200)
}

var GetNumberPlayers = func(containerId string) map[string]interface{} {
	cmd := exec.Command("docker", "exec", containerId, "rcon-cli", "list")
	outputListPlayer, err := cmd.CombinedOutput()

	if err != nil {
		return map[string]interface{}{}
	}

	re := regexp.MustCompile("[0-9]+")	
	result := strings.Split(string(outputListPlayer), ":")
	nbPlayers := result[0]
	listPlayers := strings.Split(result[1][:len(result[1])-1], ",")
	
	players := re.FindAllString(nbPlayers, -1)
	connectedPlayers := players[0]
	maxPlayers := players[1]

	// fmt.Println("connected players: " + connectedPlayers + "/" + maxPlayers)
	// fmt.Print(listPlayers)

	conP, _ := strconv.ParseInt(connectedPlayers, 10, 64)
	maxP, _ := strconv.ParseInt(maxPlayers, 10, 64)

	if conP == 0 {
		listPlayers = []string{}
	}

	return map[string]interface{}{
		"connectedPlayers": conP,
		"maxPlayers": maxP,
		"listPlayers": listPlayers,
	}
}

var GetPlayersOnline = func(w http.ResponseWriter, r *http.Request) {
	docker := &models.DockerDelete{}

	err := json.NewDecoder(r.Body).Decode(docker)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	playersOnline := map[string]interface{}{}
	playersOnline["playersOnline"] = GetNumberPlayers(docker.ContainerId)

	utils.Respond(w, playersOnline, 200)
}

var ModifyGameServer = func(w http.ResponseWriter, r *http.Request) {
	docker := &models.GameServer{}

	err := json.NewDecoder(r.Body).Decode(docker)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseUint(docker.UserId, 10, 32)

	dockerStore := &models.DockerStore{
		IdDocker: docker.ContainerId,
		UserId:   uint(userID),
	}

	err = dockerStore.UpdateGameServer(docker)

	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while updating game server"), 500)
		return
	}

	resp := map[string]interface{}{
		"game_server_updated": docker,
	}

	utils.Respond(w, resp, 200)
}


func checkIfUserCanCreate(UserId string) bool {
	userID, err := strconv.Atoi(UserId)
	user := models.GetUserFromId(int(userID))

	if user.Role == "admin" {
		return true
	}

	u64, err := strconv.ParseUint(UserId, 10, 32)

	if err != nil {
		return false
	}

	allDocker := models.UserServers(uint(u64))

	if allDocker == nil {
		return false
	}

	if len(*allDocker) >= 1 {
		return false
	}

	return true
}