package controllers

import (
	"encoding/json"
	"net/http"

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