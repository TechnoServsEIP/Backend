package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/TechnoServsEIP/Backend/tracking"

	"github.com/TechnoServsEIP/Backend/model"
	"github.com/TechnoServsEIP/Backend/utils"
	"github.com/docker/docker/client"
)

func InvitePlayer(w http.ResponseWriter, r *http.Request) {
	invitation := &model.Invitation{}

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
		tracking.LogErr("docker", err)
		log.Default().Println("error when creating docker client", err)
		utils.Respond(w, utils.Message(false, "Error when open client docker env"), 500)
		return
	}

	u64, err := strconv.ParseUint(invitation.UserId, 10, 32)

	OneDocker := model.OneUserServer(uint(u64), invitation.ContainerId)

	if OneDocker == nil {
		utils.Respond(w, utils.Message(false, "invalid user_id or container_id"), 500)
		return
	}

	info, err := cli.ContainerInspect(ctx, OneDocker.IdDocker)

	if err != nil {
		tracking.LogErr("docker", err)
		log.Default().Println(err.Error())
		utils.Respond(w, map[string]interface{}{
			"error": err.Error(),
		}, 500)
		return
	}

	adress := "x2021alsablue1371139462001.northeurope.cloudapp.azure.com:"

	adress += info.HostConfig.PortBindings["25565/tcp"][0].HostPort

	log.Default().Println(adress)

	user := &model.Account{}

	err = model.GetDB().Where("id = ?", u64).Find(user).Error
	if err != nil {
		utils.Respond(w, map[string]interface{}{
			"error": err.Error(),
		}, 400)
		return
	}

	log.Default().Println(user.Email)

	err = utils.SendInvitationEmail(user.Email, adress, invitation.Recipient)
	if err != nil {
		log.Default().Println("an error append when sending email, err: ", err)
		utils.Respond(w, utils.Message(false, "an error append when sending email"), 500)
	}
	utils.Respond(w, utils.Message(true, "The mail has been sended"), 200)
}
