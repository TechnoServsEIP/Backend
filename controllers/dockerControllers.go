package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/TechnoServsEIP/Backend/app"
	"github.com/TechnoServsEIP/Backend/models"
	"github.com/TechnoServsEIP/Backend/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func stopContainer(ctx context.Context, cli *client.Client, containerId string) error {
	return cli.ContainerStop(ctx, containerId, nil)
}

func startContainer(ctx context.Context, cli *client.Client, containerId string) error {
	return cli.ContainerStart(ctx, containerId, types.ContainerStartOptions{})
}

func GetAllPortBinded() []string {
	ctx := context.Background()

	cli, err := client.NewEnvClient()
	if err != nil {
		app.LogErr("docker", err)
		fmt.Printf("Impossible to use the docker client\n")
		return []string{}
	}

	containers := models.ListAllDockers()
	var tmpPortsBinded []string

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

func CreateDocker(w http.ResponseWriter, r *http.Request) {
	docker := &models.Docker{}

	err := json.NewDecoder(r.Body).Decode(docker)
	if err != nil {
		fmt.Println("an error occurred when decoding body, err :", err)
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	if docker.UserId == "" {
		utils.Respond(w, utils.Message(false, "User ID is empty"), http.StatusBadRequest)
		return
	}

	if !checkIfUserCanCreate(docker.UserId) {
		utils.Respond(w, utils.Message(false, "User already have max server reached"), 413)
		return
	}

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		errorLog := errors.New("error when creating docker client, err: " +
			err.Error())
		app.LogErr("docker", errorLog)
		utils.Respond(w, utils.Message(false, "An error occurred when creating container"), 500)
		return
	}

	port := utils.GetPort()

	if port == "no port available" {
		fmt.Println("An error occurred, No port available")
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
		errorLog := errors.New("Error while assigning container port, err: " +
			err.Error())
		app.LogErr("docker", errorLog)
		utils.Respond(w, utils.Message(false, "An error occurred while creating server"), 500)
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
		errorLog := errors.New("Error while creating container, err: " +
			err.Error())
		app.LogErr("docker", errorLog)
		utils.Respond(w, utils.Message(false, "An error occurred while creating server"), 500)
		return
	}

	err = startContainer(ctx, cli, cont.ID)
	if err != nil {
		errorLog := errors.New("Error while starting container, err: " +
			err.Error())
		app.LogErr("docker", errorLog)
		utils.Respond(w, utils.Message(false, "An error occurred while creating server"), 500)
		return
	}
	fmt.Println("Starting container ", cont.ID)

	u64, err := strconv.ParseUint(docker.UserId, 10, 32)
	if err != nil {
		fmt.Println("error while parsing docker id, err: ", err)
		utils.Respond(w, utils.Message(false, "An error occurred while creating server"), 500)
		return
	}

	dockerStore := &models.DockerStore{
		Game:         docker.Game,
		IdDocker:     cont.ID,
		UserId:       uint(u64),
		ServerName:   docker.ServerName,
		ServerStatus: "Started",
		LimitPlayers: 20,
	}
	resp := dockerStore.Create()

	info, err := cli.ContainerInspect(ctx, cont.ID)
	if err != nil {
		errorLog := errors.New("Error while listing docker env, err: " +
			err.Error())
		app.LogErr("docker", errorLog)
		utils.Respond(w, utils.Message(false, "An error occurred while creating server"), 500)
		return
	}

	resp["settings"] = &info
	utils.Respond(w, resp, http.StatusCreated)
}

func StartDocker(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	userId := r.Context().Value("user").(uint) //Grab the id of the user that send the request
	fmt.Println("userId: (", userId, ")")

	docker := &models.DockerDelete{}

	err := json.NewDecoder(r.Body).Decode(docker)
	if err != nil {
		fmt.Println("error while decoding body, err: ", err)
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		errorLog := errors.New("Error while creating docker env, err: " +
			err.Error())
		app.LogErr("docker", errorLog)
		utils.Respond(w, utils.Message(false, "An error occurred when starting container"), 500)
		return
	}
	err = startContainer(ctx, cli, docker.ContainerId)
	if err != nil {
		errorLog := errors.New("an error occurred when starting container" +
			"container_id=" + docker.ContainerId + "err" + err.Error())
		app.LogErr("docker", errorLog)
		utils.Respond(w, utils.Message(false, "An error occurred when starting container"), 500)
		return
	}
	fmt.Println("Starting container ", docker.ContainerId)

	info, err := cli.ContainerInspect(ctx, docker.ContainerId)
	if err != nil {
		errorLog := errors.New("an error occurred when inspecting container, err :" +
			err.Error())
		app.LogErr("docker", errorLog)
		utils.Respond(w, utils.Message(false, "An error occurred when starting container"), 500)
		return
	}

	dockerStore := &models.DockerStore{
		IdDocker: docker.ContainerId,
		UserId:   userId,
	}
	err = dockerStore.UpdateServerStatus("Started")
	if err != nil {
		_ = stopContainer(ctx, cli, docker.ContainerId)
		errorLog := errors.New("Error while updating status server, err: " +
			err.Error())
		app.LogErr("postgres", errorLog)
		utils.Respond(w, utils.Message(false, "An error occurred when starting container"), 500)
		return
	}

	dockerHistory := &models.DockerHistory{
		IdDocker:          docker.ContainerId,
		UserId:            userId,
		ActivityTimeStart: time.Now(),
	}
	fmt.Println("Insert start activity for user,  ", userId)
	_ = dockerHistory.InsertStartActivityContainer()

	resp := map[string]interface{}{}
	resp["settings"] = info
	utils.Respond(w, resp, 200)
}

func StopDocker(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	userId := r.Context().Value("user").(uint) //Grab the id of the user that send the request
	fmt.Println("userId: (", userId, ")")

	docker := &models.DockerDelete{}

	err := json.NewDecoder(r.Body).Decode(docker)
	if err != nil {
		fmt.Println("An error occurred while decoding body, err: ", err)
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		errorLog := errors.New("an error occurred when creating docker env, err: " +
			err.Error())
		app.LogErr("docker", errorLog)
		utils.Respond(w, utils.Message(false, "An error append while stopping container"), 500)
		return
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		errorLog := errors.New("an error occurred when listing docker container, err: " +
			err.Error())
		app.LogErr("docker", errorLog)
		utils.Respond(w, utils.Message(false, "An error append while stopping container"), 500)
		return
	}

	for _, _container := range containers {
		fmt.Printf("%s %s\n", _container.ID[:10], _container.Image)
	}
	fmt.Println("Stop container " + docker.ContainerId)

	err = stopContainer(ctx, cli, docker.ContainerId)
	if err != nil {
		errorLog := errors.New("An error occurred when stopping container, err: " +
			err.Error())
		app.LogErr("docker", errorLog)
		utils.Respond(w, utils.Message(false, "An error append while stopping container"), 500)
		return
	}
	dockerStore := &models.DockerStore{
		IdDocker: docker.ContainerId,
		UserId:   userId,
	}

	err = dockerStore.UpdateServerStatus("Stopped")
	if err != nil {
		errorLog := errors.New("Error while updating status server, err: " +
			err.Error())
		app.LogErr("postgres", errorLog)
		utils.Respond(w, utils.Message(false, "An error append while update server status"), 500)
		return
	}

	resp := models.InsertStopActivityContainer(userId, docker.ContainerId)
	if resp["status"] == false {
		errorMsg := errors.New(resp["message"].(string))
		app.LogErr("docker", errorMsg)
	}
	fmt.Println("Insert stop activity for user,  ", userId)

	utils.Respond(w, map[string]interface{}{"status": 200, "message": "Container Stop successfully"}, http.StatusOK)
}

func StopDockerAll(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	userId := r.Context().Value("user").(uint) //Grab the id of the user that send the request
	fmt.Println("userId: (", userId, ")")

	cli, err := client.NewEnvClient()
	if err != nil {
		errorLog := errors.New("an error occurred when creating docker env, err: " +
			err.Error())
		app.LogErr("docker", errorLog)
		utils.Respond(w, utils.Message(false, "An error append stopping containers"), 500)
		return
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		errorLog := errors.New("an error occurred when listing containers, err: " +
			err.Error())
		app.LogErr("docker", errorLog)
		utils.Respond(w, utils.Message(false, "An error append while update server status"), 500)
		return
	}

	for _, _container := range containers {
		fmt.Printf("%s %s\n", _container.ID[:10], _container.Image)

		err = cli.ContainerStop(ctx, _container.ID, nil)
		if err != nil {
			fmt.Println("An error occurred when stopping container ", _container.ID)
			utils.Respond(w, utils.Message(false, "An error append while stopping server"), 500)
			return
		}

		dockerStore := &models.DockerStore{
			IdDocker: _container.ID,
			UserId:   userId,
		}
		err = dockerStore.UpdateServerStatus("Stopped")
		if err != nil {
			errorLog := errors.New("an error occurred when updating server status, err: " +
				err.Error())
			app.LogErr("postgres", errorLog)
			utils.Respond(w, utils.Message(false, "An error append while update server status"), 500)
			return
		}
		fmt.Println("Stop container " + _container.ID)
	}

	utils.Respond(w, map[string]interface{}{"status": 200, "message": "Container Stop successfully"}, http.StatusOK)
}

func GetServerLogs(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		errorLog := errors.New("an error occurred when creating docker env, err: " +
			err.Error())
		app.LogErr("docker", errorLog)
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
		errorLog := errors.New("an error occurred when reading docker logs, err: " +
			err.Error())
		app.LogErr("docker", errorLog)
		utils.Respond(w, utils.Message(false, "Error bad container_id"), http.StatusBadRequest)
		return
	}

	var logs bytes.Buffer

	_, _ = io.Copy(&logs, out)

	resp := map[string]interface{}{
		"logs": logs.String(),
	}

	utils.Respond(w, resp, 200)
}

func DeleteDocker(w http.ResponseWriter, r *http.Request) {
	// userId := r.Context().Value("user").(uint)
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		errorLog := errors.New("an error occurred when creating docker env, err: " +
			err.Error())
		app.LogErr("docker", errorLog)
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
		errorLog := errors.New("an error occurred when inspecting container, err: " +
			err.Error())
		app.LogErr("docker", errorLog)
		utils.Respond(w, utils.Message(false, "Error while retrieving port of the container"), 500)
		return
	}

	err = cli.ContainerRemove(ctx, docker.ContainerId, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	if err != nil {
		errorLog := errors.New("an error occurred when removing container, err: " +
			err.Error())
		app.LogErr("docker", errorLog)
		fmt.Println("An error appear when removing container: ", docker.ContainerId, "err ", err)
		utils.Respond(w, utils.Message(false, "Error while removing container"), http.StatusBadRequest)
		return
	}

	utils.FreeThePort(info.HostConfig.PortBindings["25565/tcp"][0].HostPort)

	userIdUint, err := strconv.ParseUint(docker.UserId, 10, 32)
	dockerStore := &models.DockerStore{
		IdDocker: docker.ContainerId,
		UserId:   uint(userIdUint),
	}
	err = dockerStore.UpdateServerStatus("Deleted")
	if err != nil {
		errorLog := errors.New("Error while updating status server, err " +
			err.Error())
		app.LogErr("postgres", errorLog)
		utils.Respond(w, utils.Message(false, "An error append while update server status"), 500)
	}
	fmt.Println("Delete container " + docker.ContainerId)

	resp := models.RemoveContainer(uint(userIdUint), docker.ContainerId)
	utils.Respond(w, resp, 204)
}

func DeleteDockerAll(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(uint)
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		errorLog := errors.New("an error occurred when creating docker env, err: " +
			err.Error())
		app.LogErr("docker", errorLog)
		utils.Respond(w, utils.Message(false, "Error failed to contact docker api"), http.StatusBadRequest)
		return
	}

	docker := &models.DockerDelete{}
	err = json.NewDecoder(r.Body).Decode(docker)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		errorLog := errors.New("an error occurred when listing containers, err: " +
			err.Error())
		app.LogErr("docker", errorLog)
		panic(err)
	}
	for _, _container := range containers {
		fmt.Printf("%s %s\n", _container.ID[:10], _container.Image)

		// retrieve the port to free
		info, err := cli.ContainerInspect(ctx, docker.ContainerId)
		if err != nil {
			errorLog := errors.New("An error occurred when inspecting container" +
				",err: " + err.Error())
			app.LogErr("docker", errorLog)
			utils.Respond(w, utils.Message(false, "Error while retrieving port of the container"), 500)
			return
		}

		err = cli.ContainerStop(ctx, _container.ID, nil)
		if err != nil {
			errorLog := errors.New("An error occurred when stopping container " +
				_container.ID + ", err " + err.Error())
			app.LogErr("docker", errorLog)
			utils.Respond(w, utils.Message(false, "Error while stopping container"), 500)
			return
		}

		utils.FreeThePort(info.HostConfig.PortBindings["25565/tcp"][0].HostPort)

		dockerStore := &models.DockerStore{
			IdDocker: _container.ID,
			UserId:   userId,
		}
		err = dockerStore.UpdateServerStatus("Stopped")
		if err != nil {
			errorLog := errors.New("Error while updating status server:" +
				" err: " + err.Error())
			app.LogErr("postgres", errorLog)
			utils.Respond(w, utils.Message(false, "An error append while update server status"), 500)
		}
		fmt.Println("Stop container " + _container.ID)
	}

	err = cli.ContainerRemove(ctx, docker.ContainerId, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	if err != nil {
		errorLog := errors.New("An error appear when removing container: " +
			docker.ContainerId + "err " + err.Error())
		app.LogErr("docker", errorLog)
		utils.Respond(w, utils.Message(false, "Error while removing container"), http.StatusBadRequest)
		return
	}

	resp := models.RemoveContainer(userId, docker.ContainerId)
	utils.Respond(w, resp, 204)
}

func ListUserServers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	docker := &models.DockerList{}
	// userId := r.Context().Value("user").(uint)

	cli, err := client.NewEnvClient()
	if err != nil {
		errorLog := errors.New("an error occurred when creating docker env," +
			" err: " + err.Error())
		app.LogErr("docker", errorLog)
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
			errorLog := errors.New("an error occurred when inspecting container," +
				" err: " + err.Error())
			app.LogErr("docker", errorLog)
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

func GetInfosUserServer(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	docker := &models.DockerDelete{}

	cli, err := client.NewEnvClient()
	if err != nil {
		errorLog := errors.New("error when creating docker client, err: " + err.Error())
		app.LogErr("docker", errorLog)
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
		errorLog := errors.New("an error occurred when inspecting container," +
			" err: " + err.Error())
		app.LogErr("docker", errorLog)
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
		errorLog := errors.New("an error occurred when reading container logs, " +
			"err: " + err.Error())
		app.LogErr("docker", errorLog)
		utils.Respond(w, utils.Message(false, "Error bad container_id"), http.StatusBadRequest)
		return
	}

	var logs bytes.Buffer

	_, _ = io.Copy(&logs, out)
	serverInfo["logs"] = logs.String()

	serverInfo["playersOnline"] = GetNumberPlayers(docker.ContainerId)

	utils.Respond(w, serverInfo, 200)
}

func GetNumberPlayers(containerId string) map[string]interface{} {
	cmd := exec.Command("docker", "exec", containerId, "rcon-cli", "list")
	outputListPlayer, err := cmd.CombinedOutput()

	if err != nil {
		errorLog := errors.New("an error occurred when executing list user command, " +
			"err: " + err.Error())
		app.LogErr("docker", errorLog)
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
		"maxPlayers":       maxP,
		"listPlayers":      listPlayers,
	}
}

func GetPlayersOnline(w http.ResponseWriter, r *http.Request) {
	docker := &models.DockerDelete{}

	err := json.NewDecoder(r.Body).Decode(docker)
	if err != nil {
		errorLog := errors.New("an error occurred when decoding body, err: " + err.Error())
		app.LogErr("docker", errorLog)
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	playersOnline := map[string]interface{}{}
	playersOnline["playersOnline"] = GetNumberPlayers(docker.ContainerId)

	utils.Respond(w, playersOnline, 200)
}

func ModifyGameServer(w http.ResponseWriter, r *http.Request) {
	docker := &models.GameServer{}

	err := json.NewDecoder(r.Body).Decode(docker)
	if err != nil {
		errorLog := errors.New("an error occurred when decoding body, err: " + err.Error())
		app.LogErr("docker", errorLog)
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
		errorLog := errors.New("an error occured when converting userId, err: " + err.Error())
		app.LogErr("docker", errorLog)
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

func GetTotalServers(w http.ResponseWriter, r *http.Request) {
	containers := models.ListAllDockers()
	total := len(*containers)
	resp := map[string]interface{}{}

	fmt.Println(total)

	resp["total"] = total

	utils.Respond(w, resp, 200)
}

/*
 * Update the limit of number players of a minecraft server into the DB
 * Update the max-players attribut into server.properties of the concerned minecraft
 * The concerned minecraft server is not restarted
 */
func LimitNumberPlayers(w http.ResponseWriter, r *http.Request) {
	docker := &models.DockerLimitPlayers{}

	// Decode the body
	err := json.NewDecoder(r.Body).Decode(docker)
	if err != nil {
		errorLog := errors.New("an error occurred when decoding body, err: " + err.Error())
		app.LogErr("docker", errorLog)
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	// Convert userId in uint
	userID, err := strconv.ParseUint(docker.UserId, 10, 32)

	// Create a dockerStore struct
	dockerStore := &models.DockerStore{
		IdDocker: docker.ContainerId,
		UserId:   uint(userID),
	}

	// Change the limit number player into the DB
	if err := dockerStore.ChangeLimitPlayer(docker.LimitPlayers); err != nil {
		errorLog := errors.New("An error occurred when changing limit number players into the DB: " + err.Error())
		app.LogErr("postgres", errorLog)
		utils.Respond(w, utils.Message(false, "Error while changing limit number players into the DB"), 500)
		return
	}

	// Update the max players into the server.properties
	if err := models.UpdateMaxPlayers(docker.LimitPlayers, docker.ContainerId); err != nil {
		errorLog := errors.New("An error occurred when updating players into the server.properties: " + err.Error())
		app.LogErr("docker", errorLog)
		utils.Respond(w, utils.Message(false, "Error while updating players into the server.properties"), 500)
		return
	}

	// Return the limit players in response
	resp := map[string]interface{}{}
	resp["limit_players"] = docker.LimitPlayers
	utils.Respond(w, resp, 200)
}

/*
 * Update the limit of number players of a minecraft user servers into the DB
 * Update the max-players attribut into server.properties of the minecraft user servers
 * The concerned minecraft servers is not restarted
 */
func LimitNumberPlayersOfUserServers(w http.ResponseWriter, r *http.Request) {
	docker := &models.DockerLimitPlayersUserServers{}

	// Decode the body
	err := json.NewDecoder(r.Body).Decode(docker)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	// Check if userId exist
	if docker.UserId == "" {
		utils.Respond(w, utils.Message(false, "User ID is empty"), http.StatusBadRequest)
		return
	}

	// Convert userId in uint
	userID, err := strconv.ParseUint(docker.UserId, 10, 32)

	// Get all user servers
	allDocker := models.UserServers(uint(userID))

	// Trigger a invalid userId
	if allDocker == nil {
		utils.Respond(w, utils.Message(false, "invalid user_id"), 500)
		return
	}

	// Create a empty list of docker store
	list := make([]models.DockerStore, 0)

	// Change limit number player into the DB and update max players into the server.properties for each server
	for _, element := range *allDocker {
		// Change the limit number player into the DB
		if err := element.ChangeLimitPlayer(docker.LimitPlayers); err != nil {
			errorLog := errors.New("An error occurred when changing limit number players into the DB: " + err.Error())
			app.LogErr("postgres", errorLog)
			utils.Respond(w, utils.Message(false, "Error while changing limit number players into the DB"), 500)
			return
		}

		// Update the max players into the server.properties
		if err := models.UpdateMaxPlayers(docker.LimitPlayers, element.IdDocker); err != nil {
			errorLog := errors.New("An error occurred when updating players into the server.properties: " + err.Error())
			app.LogErr("docker", errorLog)
			utils.Respond(w, utils.Message(false, "Error while updating players into the server.properties"), 500)
			return
		}

		// Append a server into the list
		list = append(list, element)
	}

	// Return user servers updated
	resp := map[string]interface{}{}
	resp["user_servers"] = list
	utils.Respond(w, resp, 200)
}
