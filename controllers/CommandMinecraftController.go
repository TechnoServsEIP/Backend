package controllers

import (
	"encoding/json"
	"errors"
	"github.com/TechnoServsEIP/Backend/app"
	"github.com/TechnoServsEIP/Backend/models"
	"github.com/TechnoServsEIP/Backend/utils"
	"net/http"
	"os/exec"
)

type Command struct {
	UserID uint
	DockerID string
	Command string
}

func CommandRoute(w http.ResponseWriter, request *http.Request ) {
	defer request.Body.Close()
	command := &Command{}
	err := json.NewDecoder(request.Body).Decode(command)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid request"), 400)
		return
	}
	if command.Command == "stop" {
		utils.Respond(w, utils.Message(false, "Invalid Request"), 400)
	}
	dockers := models.UserServers(command.UserID)
	for _, docker := range(*dockers) {
		if (docker.IdDocker == command.DockerID) {
			cmd := exec.Command("docker", "exec", command.DockerID, "rcon-cli", command.Command)
			output, err := cmd.CombinedOutput()
			if err != nil {
				errorLog := errors.New("an error occurred when executing" + command.Command+ "user command, " +
					"err: " + err.Error())
				app.LogErr("docker", errorLog)
				utils.Respond(w, utils.Message(false, "Invalid Request"), 400)
				return
			}
			utils.Respond(w, utils.Message(true, string(output)), 200)
		}
	}
}
