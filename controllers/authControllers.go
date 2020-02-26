package controllers

import (
	"fmt"
	"net/http"
	"oauth2server/utils"
	"oauth2server/models"
	"encoding/json"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request /user/new")
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		fmt.Println("An error occurred while decoding request ", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), 400)
		return
	}

	resp := account.Create() //Create account
	utils.Respond(w, resp, 201)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid request"), 400)
		return
	}

	resp := models.Login(account.Email, account.Password)
	utils.Respond(w, resp, 200)
}
