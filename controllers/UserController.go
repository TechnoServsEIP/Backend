package controllers

import (
	"gitlab.sysroot.ovh/technoservs/microservices/game-servers/models"
	"gitlab.sysroot.ovh/technoservs/microservices/game-servers/utils"
	"net/http"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	res := models.GetUsers()
	utils.Respond(w, map[string]interface{}{"res": res}, 200)
}