package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TechnoServsEIP/Backend/app"
	"github.com/TechnoServsEIP/Backend/models"
	"github.com/TechnoServsEIP/Backend/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func GetServerProperties(w http.ResponseWriter, r *http.Request) {
	docker := &models.DockerDelete{}

	err := json.NewDecoder(r.Body).Decode(docker)
	if err != nil {
		app.LogErr("docker", err)
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	data := models.ServerPropertiesByServerId(docker.ContainerId)

	if data == nil {
		utils.Respond(w, utils.Message(false, "Error while retrieve server properties"), 500)
		return
	}

	resp := map[string]interface{}{
		"properties": data,
	}

	utils.Respond(w, resp, 200)
}

func RestartServer(containerId string, userId string) error {
	ctx := context.Background()

	cli, err := client.NewEnvClient()
	if err != nil {
		app.LogErr("docker", err)
		fmt.Println("An error occurred when stopping container ", containerId)
		return err
	}

	err = cli.ContainerStop(ctx, containerId, nil)
	if err != nil {
		fmt.Println("An error occurred when stopping container or the container is already stopped", containerId)
	}

	u64, err := strconv.ParseUint(userId, 10, 32)

	dockerStore := &models.DockerStore{
		IdDocker: containerId,
		UserId:   uint(u64),
	}

	dockerStore.UpdateServerStatus("Stoped")

	err = cli.ContainerStart(ctx, containerId, types.ContainerStartOptions{})
	if err != nil {
		app.LogErr("docker", err)
		fmt.Println("An error occurred when starting container", "container_id=", containerId, "err", err)
		return err
	}

	dockerStore.UpdateServerStatus("Started")

	return nil
}

func UpdateServerProperties(w http.ResponseWriter, r *http.Request) {
	data := &models.UpdateServerProperties{}

	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		app.LogErr("docker", err)
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	if data.ContainerId == "" {
		utils.Respond(w, utils.Message(false, "Error, container_id is empty"), 401)
		return
	}

	// Convert userId in uint
	userID, err := strconv.ParseUint(data.UserId, 10, 32)

	// Create a dockerStore struct
	dockerStore := &models.DockerStore{
		IdDocker: data.ContainerId,
		UserId:   uint(userID),
	}

	maxPlayers, err := dockerStore.GetLimitPlayer()

	if err != nil {
		errorLog := errors.New("An error occurred when get limit number players from the DB: " + err.Error())
		app.LogErr("postgres", errorLog)
	}

	if err := models.CreateNewServerProperties(*data, data.ContainerId, maxPlayers); err != nil {
		app.LogErr("docker", err)
		utils.Respond(w, utils.Message(false, "Error while updating server properties"), 500)
		return
	}

	if err := RestartServer(data.ContainerId, data.UserId); err != nil {
		app.LogErr("docker", err)
		utils.Respond(w, utils.Message(false, "Error while updating server properties"), 500)
		return
	}

	utils.Respond(w, map[string]interface{}{}, 204)
}
