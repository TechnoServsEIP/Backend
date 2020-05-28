package controllers

import (
	"gitlab.sysroot.ovh/technoservs/microservices/game-servers/utils"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	message := utils.Message(true, "Welcome to the api")
	utils.Respond(w, message, 200)
}