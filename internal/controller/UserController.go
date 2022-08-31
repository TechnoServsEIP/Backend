package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/now"

	"github.com/TechnoServsEIP/Backend/model"
	"github.com/TechnoServsEIP/Backend/utils"
)

type UserController struct {
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	res := model.GetUsers()
	utils.Respond(w, map[string]interface{}{"res": res}, 200)
}

func Activate(w http.ResponseWriter, r *http.Request) {
	idJson := &struct {
		Id int `json:"Id,string,omitempty"`
	}{}
	err := json.NewDecoder(r.Body).Decode(idJson)
	if err != nil {
		utils.Respond(w, utils.Message(false, "malformed request"), 400)
		return
	}
	id := int(idJson.Id)
	res := model.ActivateUser(id)
	utils.Respond(w, map[string]interface{}{"res": res}, 200)
}

func Deactivate(w http.ResponseWriter, r *http.Request) {
	idJson := &struct {
		Id int `json:"Id,string,omitempty"`
	}{}
	err := json.NewDecoder(r.Body).Decode(idJson)
	if err != nil {
		utils.Respond(w, utils.Message(false, "malformed request"), 400)
		return
	}
	id := int(idJson.Id)
	log.Default().Println(id)
	res := model.DeactivateUser(id)
	utils.Respond(w, map[string]interface{}{"res": res}, 200)
}

func VerifyUser(w http.ResponseWriter, r *http.Request) {
	idJson := &struct {
		Id int `json:"Id,string,omitempty"`
	}{}
	err := json.NewDecoder(r.Body).Decode(idJson)
	if err != nil {
		utils.Respond(w, utils.Message(false, "malformed request"), 400)
		return
	}
	id := int(idJson.Id)
	res := model.VerifyUser(id)
	utils.Respond(w, map[string]interface{}{"res": res}, 200)
}

func RemoveVerification(w http.ResponseWriter, r *http.Request) {
	idJson := &struct {
		Id int `json:"Id,string,omitempty"`
	}{}
	err := json.NewDecoder(r.Body).Decode(idJson)
	if err != nil {
		utils.Respond(w, utils.Message(false, "malformed request"), 400)
		return
	}
	id := int(idJson.Id)
	log.Default().Println(id)
	res := model.RemoveVerification(id)
	utils.Respond(w, map[string]interface{}{"res": res}, 200)
}

func SendPasswordReset(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := struct {
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
	data := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	user := &model.Account{}

	msgSuccess := utils.Message(true, "password change")
	msgFailure := utils.Message(false, "request failed")
	missingPassword := utils.Message(false, "Missing password")

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		utils.Respond(w, msgFailure, 400)
		return
	}
	if data.Password == "" {
		utils.Respond(w, missingPassword, 400)
		return
	}

	err = model.GetDB().Where("email = ?", data.Email).Find(user).Error
	if err != nil {
		utils.Respond(w, msgFailure, 400)
		return
	}
	err = model.ChangePassword(data.Password, user.ID)
	if err != nil {
		utils.Respond(w, msgFailure, 400)
		return
	}

	utils.Respond(w, msgSuccess, 200)
}

func GetEmail(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userId := r.Context().Value("user")
	msgFailure := utils.Message(false, "request failed")

	user := &model.Account{}
	err := model.GetDB().Where("id = ?", userId).Find(user).Error
	if err != nil {
		utils.Respond(w, msgFailure, 400)
		return
	}
	msgSuccess := utils.Message(true, user.Email)
	utils.Respond(w, msgSuccess, 200)
}

func GetActivityByUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userId := r.Context().Value("user").(uint)

	resp := model.GetUserActivity(userId)

	resp["total_time"] = getTotalTimeActivityPerMonth(userId,
		now.BeginningOfMonth())
	log.Default().Println(resp)

	utils.Respond(w, resp, 200)
}

func getTotalTimeActivity(userId uint) time.Duration {
	resp := model.GetUserActivity(userId)
	dockers := resp["docker"].([]model.DockerHistory)
	var totalDuration time.Duration

	for _, docker := range dockers {
		totalDuration += docker.ActivityTimeStop.
			Sub(docker.ActivityTimeStart)
	}
	return totalDuration
}

func getTotalTimeActivityPerMonth(userId uint, currentMonth time.Time) float64 {
	resp := model.GetUserActivity(userId)
	dockers := resp["docker"].([]model.DockerHistory)
	var currentMonthDuration time.Duration
	endOfCurrentMonth := currentMonth.Month() * 1

	log.Default().Println("currentMonth: ", currentMonth)
	log.Default().Println("endof current month: ", endOfCurrentMonth)

	for _, docker := range dockers {
		currentMonthDuration += docker.ActivityTimeStop.
			Sub(docker.ActivityTimeStart)
	}
	return float64(currentMonthDuration / time.Hour)
}

func GetTotalToPaidPerMonthByUser(userId uint) float64 {
	currentMonth := now.BeginningOfMonth()
	totalTimeActivity := getTotalTimeActivityPerMonth(userId, currentMonth)
	pricePerHour := 0.35

	return totalTimeActivity * pricePerHour
}

func GetBillsByUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userId := r.Context().Value("user").(uint)

	resp := model.GetBillsByUser(userId)
	log.Default().Println(resp)
	if resp["status"] == false {
		utils.Respond(w, resp, 404)
	}

	utils.Respond(w, resp, 200)
}

func InsertBills(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userId := r.Context().Value("user").(uint)

	req := struct {
		Email       string `json:"email"`
		Product     string `json:"product"`
		PriceToPaid string `json:"price"`
	}{}
	msgFailure := utils.Message(false, "request failed")
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.Respond(w, msgFailure, 400)
		return
	}

	currentTime := time.Now()
	bill := &model.Bill{
		UserId:       userId,
		Email:        req.Email,
		Price:        req.PriceToPaid,
		Product:      req.Product,
		StartSubDate: currentTime,
		EndSubDate:   currentTime.AddDate(0, 1, 0),
	}
	bill.InsertBill()
	log.Default().Println(*bill)

	resp := model.GetBillsByUser(userId)
	log.Default().Println(resp)
	if resp["status"] == false {
		utils.Respond(w, resp, 404)
	}

	utils.Respond(w, resp, 200)
}
