package controllers

import (
	"encoding/json"
	"fmt"
	"gitlab.sysroot.ovh/technoservs/microservices/game-servers/models"
	"gitlab.sysroot.ovh/technoservs/microservices/game-servers/utils"
	"net/http"
)

func ListOffers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request /offers/list")

	resp := models.GetOfferList()

	utils.Respond(w, resp, 200)
}

func GetOffer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request /offers/{uuid}")

	resp := models.GetOffer(r.URL.Query()["uuid"][0])

	utils.Respond(w, resp, 200)
}

func CreateOffer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request /offers/create")

	offer := &models.Offer{}
	err := json.NewDecoder(r.Body).Decode(offer)
	if err != nil {
		fmt.Println("An error occurred while decoding request ", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), 400)
		return
	}

	resp := offer.Create()

	utils.Respond(w, resp, 201)
}

func UpdateOffer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request /offers/update/{uuid}")

	offer := &models.Offer{}
	err := json.NewDecoder(r.Body).Decode(offer)
	if err != nil {
		fmt.Println("An error occurred while decoding request ", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), 400)
		return
	}

	resp := offer.Update(r.URL.Query()["uuid"][0])

	utils.Respond(w, resp, 200)
}

func DeleteOffer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request /offers/delete/{uuid}")

	resp := models.Delete(r.URL.Query()["uuid"][0])

	utils.Respond(w, resp, 204)
}