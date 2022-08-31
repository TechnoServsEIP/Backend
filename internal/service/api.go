package service

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	controllers "github.com/TechnoServsEIP/Backend/internal/controller"
	"github.com/TechnoServsEIP/Backend/internal/model"
	"github.com/TechnoServsEIP/Backend/internal/utils"
	"github.com/TechnoServsEIP/Backend/tracking"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type TechnoservsService struct {
	UserController *controllers.UserController
	AuthController *controllers.AuthController
}

func (s *TechnoservsService) Router() *mux.Router {
	router := mux.NewRouter()

	// r.HandleFunc("/", HomeHandler).Methods("GET")
	// r.Methods("GET").Name("GetFizzBuzz").Handler(s.FizzBuzzHandler()).Path("/fizzbuzz")
	// r.Methods("GET").Name("GetFizzBuzzStats").Handler(s.GetFizzBuzzStatsHandler()).Path("/fizzbuzz/stats")

	router.HandleFunc("/", HomeHandler).Methods("GET")
	router.Methods("POST").Name("RefreshToken").Handler(s.RefreshToken()).Path("/token/refresh")
	router.Methods("POST").Name("RevokeToken").Handler(s.RefreshToken()).Path("/token/revoke")
	router.HandleFunc("/token/revoke", s.AuthController.RevokeToken).Methods("POST")
	router.HandleFunc("/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/user/update", controllers.UpdateAccount).Methods("POST")
	router.HandleFunc("/user/confirm", controllers.Confirm).Methods("POST")
	router.HandleFunc("/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/user/currentUser", controllers.GetEmail).Methods("GET")
	router.HandleFunc("/user/", controllers.GetUsers).Methods("GET")
	router.HandleFunc("/user/activate", controllers.Activate).Methods("POST")
	router.HandleFunc("/user/deactivate", controllers.Deactivate).Methods("POST")
	router.HandleFunc("/user/verify", controllers.VerifyUser).Methods("POST")
	router.HandleFunc("/user/removeverification", controllers.RemoveVerification).Methods("POST")
	router.HandleFunc("/user/forgotpassword", controllers.SendPasswordReset).Methods("POST")
	router.HandleFunc("/user/resetpassword", controllers.ChangePassword).Methods("POST")
	// router.HandleFunc("/user/getactivitybyuser", controllers.GetActivityByUser).Methods("POST")
	// router.HandleFunc("/user/getbills", controllers.GetBillsByUser).Methods("POST")
	// router.HandleFunc("/user/insertbill", controllers.InsertBills).Methods("POST")
	// router.HandleFunc("/docker/create", controllers.CreateDocker).Methods("POST")
	// router.HandleFunc("/docker/start", controllers.StartDocker).Methods("POST")
	// router.HandleFunc("/docker/stop", controllers.StopDocker).Methods("POST")
	// router.HandleFunc("/docker/stopAll", controllers.StopDockerAll).Methods("POST")
	// router.HandleFunc("/docker/delete", controllers.DeleteDocker).Methods("POST")
	// router.HandleFunc("/docker/deleteAll", controllers.DeleteDockerAll).Methods("POST")
	// router.HandleFunc("/docker/logs", controllers.GetServerLogs).Methods("POST")
	// router.HandleFunc("/docker/list", controllers.ListUserServers).Methods("POST")
	// router.HandleFunc("/docker/infos", controllers.GetInfosUserServer).Methods("POST")
	// router.HandleFunc("/docker/playersonline", controllers.GetPlayersOnline).Methods("POST")
	// router.HandleFunc("/docker/limitnumberplayers", controllers.LimitNumberPlayers).Methods("POST")
	// router.HandleFunc("/docker/limitnumberplayersofuserservers", controllers.LimitNumberPlayersOfUserServers).Methods("POST")
	// router.HandleFunc("/docker/total", controllers.GetTotalServers).Methods("GET")
	// router.HandleFunc("/minecraft/getserverproperties", controllers.GetServerProperties).Methods("POST")
	// router.HandleFunc("/minecraft/updateserverproperties", controllers.UpdateServerProperties).Methods("POST")
	// router.HandleFunc("/docker/update", controllers.ModifyGameServer).Methods("POST")
	// router.HandleFunc("/payment/new", controllers.PaymentNew).Methods("POST")
	// router.HandleFunc("/payment/renew", controllers.PaymentRenew).Methods("POST")
	// router.HandleFunc("/invitation", controllers.InvitePlayer).Methods("POST")
	// router.HandleFunc("/Command", controllers.CommandRoute).Methods("POST")

	router.Use(JwtAuthentication)

	return router
}

func (a *TechnoservsService) RefreshToken() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Default().Println("request /refreshToken")

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

		a.AuthController.RefreshToken(refreshTokenRequest.RefreshToken, refreshTokenRequest.AccessToken)

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
			log.Default().Println("everything is fine until here lets test user")
			log.Default().Println("userid: ", rtk.UserId)
			log.Default().Println("test expiration")
			tmp := rtk.VerifyExpiresAt(time.Now().Unix(), true)
			if !tmp {
				log.Default().Println("token expired")
				return
			}
			log.Default().Println("token good !")
		}
		user := model.GetUserFromId(int(rtk.UserId))
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
	})
}

func (a *TechnoservsService) GetFizzBuzzStatsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Default().Println("request /fizzbuzz/stats, fetching stats")

		resp, err := a.FizzbuzzController.GetFizzbuzzStats()
		if err != nil {
			log.Default().Println("an error append in GetFizzBuzzStatsHandler, err :", err)
			utils.Respond(w, utils.Message(false, "Internal error"), 500)
		}

		utils.Respond(w, resp, 200)
	})
}
