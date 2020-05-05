package controllers

import (
	"encoding/json"
	"fmt"
	"gitlab.sysroot.ovh/technoservs/microservices/game-servers/app"
	"gitlab.sysroot.ovh/technoservs/microservices/game-servers/models"
	"gitlab.sysroot.ovh/technoservs/microservices/game-servers/utils"
	"net/http"
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

var Confirm = func(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request /user/confirm")
	token := r.URL.Query()["token"][0]
	fmt.Println(token)
	claims, valid, err := app.DecryptToken(token)
	if !valid {
		fmt.Println("invalid token")
		if err != nil {
			fmt.Println("error ", err)
		}
		return
	}

	fmt.Println("after decrypt")

	user := models.GetUserFromId(int(claims.(*models.Token).UserId))
	fmt.Println(user)
	user.Verified = true
	models.Update(int(user.ID), map[string]interface{}{
	"verified": true,
	})
	//c.Redirect(http.StatusPermanentRedirect, "https://localhost:8000/#/login")
}
