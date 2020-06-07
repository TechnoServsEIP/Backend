package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TechnoServsEIP/Backend/models"
	"github.com/TechnoServsEIP/Backend/utils"
)

func ListOffers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request /offers/list")

	resp := models.GetOfferList()

	utils.Respond(w, resp, 200)
}

func GetOffer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request /offers/")
	data := &models.GetUuid{}

	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		fmt.Println("An error occurred while decoding request ", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), 400)
		return
	}
	fmt.Println(data.UUID)
	res := models.GetOffer(data.UUID)
	utils.Respond(w, res, 200)
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
	fmt.Println("request /offers/update")

	offer := &models.Offer{}
	err := json.NewDecoder(r.Body).Decode(offer)
	if err != nil {
		fmt.Println("An error occurred while decoding request ", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), 400)
		return
	}

	resp := offer.Update(offer.UUID)

	utils.Respond(w, resp, 200)
}

func DeleteOffer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request /offers/delete")

	data := &models.GetUuid{}

	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		fmt.Println("An error occurred while decoding request ", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), 400)
		return
	}
	fmt.Println(data.UUID)

	resp := models.Delete(data.UUID)

	utils.Respond(w, resp, 204)
}
