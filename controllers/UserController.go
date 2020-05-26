package controllers

import (
	"encoding/json"
	"fmt"
	"gitlab.sysroot.ovh/technoservs/microservices/game-servers/models"
	"gitlab.sysroot.ovh/technoservs/microservices/game-servers/utils"
	"net/http"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	res := models.GetUsers()
	utils.Respond(w, map[string]interface{}{"res": res}, 200)
}

func Activate(w http.ResponseWriter, r *http.Request) {
	idJson := &struct{Id int `json:"Id,string,omitempty"`}{}
	err := json.NewDecoder(r.Body).Decode(idJson)
	if err != nil {
		utils.Respond(w, utils.Message(false, "malformed request"), 400)
		return
	}
	id := int(idJson.Id)
	res := models.ActivateUser(id)
	utils.Respond(w, map[string]interface{}{"res": res}, 200)
}

func Deactivate(w http.ResponseWriter, r *http.Request) {
	idJson := &struct{Id int `json:"Id,string,omitempty"`}{}
	err := json.NewDecoder(r.Body).Decode(idJson)
	if err != nil {
		utils.Respond(w, utils.Message(false, "malformed request"), 400)
		return
	}
	id := int(idJson.Id)
	fmt.Println(id)
	res := models.DeactivateUser(id)
	utils.Respond(w, map[string]interface{}{"res": res}, 200)
}
