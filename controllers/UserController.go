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

func SendPasswordReset(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := struct{
		Email string
	}{}
	url := "localhost:8080"
	msgSuccess := utils.Message(true, "password change email send")
	msgFailure := utils.Message(false, "request failed")
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		utils.Respond(w, msgFailure, 400)
		return
	}
	err = utils.SendConfirmationEmail(url, data.Email)
	if err != nil {
		utils.Respond(w, msgFailure, 400)
		return
	}

	utils.Respond(w, msgSuccess, 400)
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := struct{
		Email string
		Password string
	}{}
	user := &models.Account{}
	msgSuccess := utils.Message(true, "password change")
	msgFailure := utils.Message(false, "request failed")
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		utils.Respond(w, msgFailure, 400)
		return
	}
	err = models.GetDB().Where("email = ?", data.Email).Find(user).Error
	if err != nil {
		utils.Respond(w, msgFailure, 400)
		return
	}
	err = models.ChangePassword(data.Password, user.ID)
	if err != nil {
		utils.Respond(w, msgFailure, 400)
		return
	}

	utils.Respond(w, msgSuccess, 400)
}

func GetEmail(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userId := r.Context().Value("user")
	user := &models.Account{}
	msgFailure := utils.Message(false, "request failed")

	err := models.GetDB().Where("id = ?", userId).Find(user).Error
	if err != nil {
		utils.Respond(w, msgFailure, 400)
		return
	}
	msgSuccess := utils.Message(true, user.Email)
	utils.Respond(w, msgSuccess, 200)
}