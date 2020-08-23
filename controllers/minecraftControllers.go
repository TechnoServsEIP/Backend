package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/TechnoServsEIP/Backend/models"
	"github.com/TechnoServsEIP/Backend/utils"
)

func GetServerProperties(w http.ResponseWriter, r *http.Request) {
	docker := &models.DockerDelete{}

	err := json.NewDecoder(r.Body).Decode(docker)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	data := models.ServerPropertiesByServerId(docker.ContainerId)

	if (data == nil) {
		utils.Respond(w, utils.Message(false, "Error while retrieve server properties"), 500)
		return
	}

	resp := map[string]interface{}{
		"properties": data,
	}

	utils.Respond(w, resp, 200)
}

func restartServer(containerId string, userId string) error {
	ctx := context.Background()

	cli, err := client.NewEnvClient()
	if err != nil {
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
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	if err := models.CreateNewServerProperties(*data, data.ContainerId); err != nil {
		utils.Respond(w, utils.Message(false, "Error while updating server properties"), 500)
		return
	}

	if err := restartServer(data.ContainerId, data.UserId); err != nil {
		utils.Respond(w, utils.Message(false, "Error while updating server properties"), 500)
		return
	}

	utils.Respond(w, map[string]interface{}{}, 204)
}