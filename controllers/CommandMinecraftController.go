package controllers

import (
	"encoding/json"
	"errors"
	"github.com/TechnoServsEIP/Backend/models"
	"github.com/TechnoServsEIP/Backend/tracking"
	"github.com/TechnoServsEIP/Backend/utils"
	"net/http"
	"os/exec"
	"strings"
)

type Command struct {
	UserID   uint   `json:"user_id"`
	DockerID string `json:"docker_id"`
	Command  string `json:"command"`
}

func CommandRoute(w http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	command := &Command{}
	err := json.NewDecoder(request.Body).Decode(command)
	tmp := strings.ToLower(command.Command)
	command.Command = tmp
	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid request"), 400)
		return
	}
	if command.Command == "stop" {
		utils.Respond(w, utils.Message(false, "This command is forbidden"), 404)
		return
	}
	dockers := models.UserServers(command.UserID)
	for _, docker := range *dockers {
		if docker.IdDocker == command.DockerID {
			cmd := exec.Command("docker", "exec", command.DockerID, "rcon-cli", command.Command)
			output, err := cmd.CombinedOutput()
			if err != nil {
				errorLog := errors.New("an error occurred when executing" + command.Command + "user command, " +
					"err: " + err.Error())
				tracking.LogErr("docker", errorLog)
				utils.Respond(w, utils.Message(false, "Invalid Request"), 400)
				return
			}
			utils.Respond(w, utils.Message(true, string(output)), 200)
		}
	}
}
