package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TechnoServsEIP/Backend/tracking"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"time"

	"github.com/TechnoServsEIP/Backend/app"
	"github.com/TechnoServsEIP/Backend/models"
	"github.com/TechnoServsEIP/Backend/utils"
)

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request /user/new")
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		errorLog := errors.New("An error occurred while decoding request, err: " +
			err.Error())
		tracking.LogErr("jwt", errorLog)
		utils.Respond(w, utils.Message(false, "Invalid request"), 400)
		return
	}

	resp := account.Create() //Create account
	utils.Respond(w, resp, 201)
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	type tokenReqBody struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	refreshTokenRequest := &tokenReqBody{}

	err := json.NewDecoder(r.Body).Decode(refreshTokenRequest)
	if err != nil {
		errorLog := errors.New("An error occurred while decoding request, err: " +
			err.Error())
		tracking.LogErr("jwt", errorLog)
		utils.Respond(w, utils.Message(false, "Invalid request"), 400)
		return
	}
	rtk := &models.RefreshToken{}

	response := make(map[string]interface{})
	refreshToken, err := jwt.ParseWithClaims(refreshTokenRequest.RefreshToken, rtk,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})
	if err != nil { //Malformed token, returns with http code 403 as usual
		errorLog := errors.New("Malformed or expired refresh token, err: " +
			err.Error())
		tracking.LogErr("jwt", errorLog)
		response = utils.Message(false, "Malformed authentication token")
		w.Header().Add("Content-Type", "application/json")
		utils.Respond(w, response, http.StatusForbidden)
		return
	}

	if refreshToken.Valid && refreshToken.Claims.Valid() == nil {
		fmt.Println("everything is fine until here lets test user")
		fmt.Println("userid: ", rtk.UserId)
		fmt.Println("test expiration")
		tmp := rtk.VerifyExpiresAt(time.Now().Unix(), true)
		if !tmp {
			fmt.Println("token expired")
			return
		}
		fmt.Println("token good !")
	}
	user := models.GetUserFromId(int(rtk.UserId))
	resp, err := user.GenerateJWT()
	if err != nil {
		errorLog := errors.New("Error append when generating refresh token, err: " +
			err.Error())
		tracking.LogErr("jwt", errorLog)
		utils.Respond(w, utils.Message(false, "An error append when generating refresh token"),
			500)
		return
	}

	utils.Respond(w, map[string]interface{}{
		"access_token":  resp["access_token"],
		"refresh_token": resp["refresh_token"],
	}, 200)
}

func RevokeToken(w http.ResponseWriter, r *http.Request) {
	type tokenReqBody struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	TokenRequest := &tokenReqBody{}
	act := &models.Token{}

	err := json.NewDecoder(r.Body).Decode(TokenRequest)
	if err != nil {
		errorLog := errors.New("An error occurred while decoding request, err: " +
			err.Error())
		tracking.LogErr("jwt", errorLog)
		utils.Respond(w, utils.Message(false, "Invalid request"), 400)
		return
	}
	resp := make(map[string]interface{})

	_, err = jwt.ParseWithClaims(TokenRequest.AccessToken, act,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})
	if err != nil { //Malformed token, returns with http code 403 as usual
		errorLog := errors.New("Malformed or expired refresh token, err: " +
			err.Error())
		tracking.LogErr("jwt", errorLog)
		resp = utils.Message(false, "Malformed authentication token")
		w.Header().Add("Content-Type", "application/json")
		utils.Respond(w, resp, http.StatusForbidden)
		return
	}

	user := models.GetUserFromId(int(act.UserId))

	ack := &models.Token{
		UserId:   user.ID,
		Role:     user.Role,
		IsRevoke: true,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, ack)
	//Generating refresh_token with user_id, and exp duration
	rtk := &models.RefreshToken{
		UserId:   user.ID,
		IsRevoke: true,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtk)

	tokenString, err := accessToken.SignedString([]byte(os.Getenv("token_password")))
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("token_password")))

	utils.Respond(w, map[string]interface{}{
		"access_token":  tokenString,
		"refresh_token": refreshTokenString,
	}, 200)
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		tracking.LogErr("jwt", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), 400)
		return
	}

	resp := models.Login(account.Email, account.Password)
	if resp["status"] == false {
		utils.Respond(w, resp, 400)
		return
	}
	utils.Respond(w, resp, 200)
}

func Confirm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request /user/confirm")
	token := r.URL.Query()["token"][0]
	fmt.Println("len(token)")
	fmt.Println(len(token))
	fmt.Println(token)
	claims, valid, err := app.DecryptToken(token)
	if !valid {
		fmt.Println("invalid token")
		if err != nil {
			tracking.LogErr("jwt", err)
			fmt.Println("error ", err)
		}
		return
	}

	user := models.GetUserFromId(int(claims.(*models.Token).UserId))
	fmt.Println(user)
	user.Verified = true
	models.Update(int(user.ID), map[string]interface{}{
		"verified": true,
	})
}

func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fmt.Println("request /user/update")

	idJson := &struct {
		Id   int `json:"Id,string,omitempty"`
		Role string
	}{}
	err := json.NewDecoder(r.Body).Decode(idJson)
	if err != nil {
		tracking.LogErr("jwt", err)
		println(err.Error())
		utils.Respond(w, utils.Message(false, "malformed request"), 400)
		return
	}
	id := int(idJson.Id)
	fmt.Println(id)
	fmt.Println(idJson.Role)
	models.Update(id, map[string]interface{}{
		"role": idJson.Role,
	})
	response := utils.Message(true, "role update")
	utils.Respond(w, response, 200)
}
