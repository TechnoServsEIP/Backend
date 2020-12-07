package controllers

import (
	"github.com/TechnoServsEIP/Backend/utils"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	message := utils.Message(true, "Welcome to the api")
	utils.Respond(w, message, 200)
}
