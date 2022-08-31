package service

import (
	"net/http"

	"github.com/TechnoServsEIP/Backend/internal/utils"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	message := utils.Message(true, "Welcome to Technoservs")
	utils.Respond(w, message, 200)
}
