package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TechnoServsEIP/Backend/app"
	"github.com/TechnoServsEIP/Backend/models"
	"github.com/TechnoServsEIP/Backend/utils"
	"github.com/docker/docker/client"
)

var Invite = func(w http.ResponseWriter, r *http.Request) {
	invitation := &models.Invitation{}

	err := json.NewDecoder(r.Body).Decode(invitation)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	if invitation.UserId == "" {
		utils.Respond(w, utils.Message(false, "User ID is empty"), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		app.LogErr("docker", err)
		fmt.Println("error when creating docker client", err)
		utils.Respond(w, utils.Message(false, "Error when open client docker env"), 500)
		return
	}

	u64, err := strconv.ParseUint(invitation.UserId, 10, 32)

	OneDocker := models.OneUserServer(uint(u64), invitation.ContainerId)

	if OneDocker == nil {
		utils.Respond(w, utils.Message(false, "invalid user_id or container_id"), 500)
		return
	}

	info, err := cli.ContainerInspect(ctx, OneDocker.IdDocker)

	if err != nil {
		app.LogErr("docker", err)
		fmt.Println(err.Error())
		utils.Respond(w, map[string]interface{}{
			"error": err.Error(),
		}, 500)
		return
	}

	adress := "https://x2021alsablue1371139462001.northeurope.cloudapp.azure.com:"

	adress += info.HostConfig.PortBindings["25565/tcp"][0].HostPort

	fmt.Println(adress)

	user := &models.Account{}

	err = models.GetDB().Where("id = ?", u64).Find(user).Error
	if err != nil {
		utils.Respond(w, map[string]interface{}{
			"error": err.Error(),
		}, 400)
		return
	}

	fmt.Println(user.Email)

	utils.SendInvitationEmail(user.Email, adress, invitation.Recipient)

	utils.Respond(w, utils.Message(true, "The mail has been sended"), 200)
}
