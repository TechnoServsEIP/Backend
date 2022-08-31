package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/TechnoServsEIP/Backend/internal/model"
	"github.com/TechnoServsEIP/Backend/internal/utils"
	"github.com/TechnoServsEIP/Backend/tracking"
	"github.com/golang-jwt/jwt"
)

func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{"/token/revoke", "/loggedin", "/login/github/callback", "/login/github/", "/user/new", "/user/login", "/user/confirm", "/offers/list", "/offers/", "/docker/list", "/docker/create", "/user/forgotpassword", "/user/resetpassword", "/", "/token/refresh"} //List of endpoints that doesn't require auth
		adminOnlyPath := []string{"/user/verify", "/user/removeverification", "/docker/total", "/user/update", "/offers/create", "/offers/delete", "/offers/update", "/user/activate", "/user/deactivate", "/docker/deleteAll", "/docker/stopAll", "/docker/limitnumberplayers", "/docker/limitnumberplayersofuserservers"}
		requestPath := r.URL.Path //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			log.Default().Println("token is empty")
			response = utils.Message(false, "Missing auth token")
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response, http.StatusForbidden)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			response = utils.Message(false, "Invalid/Malformed auth token")
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response, http.StatusForbidden)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
		tk := &model.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("token_password")), nil
			})
		if err != nil { //Malformed token, returns with http code 403 as usual
			errorLog := errors.New("Malformed or expired token, err: " +
				err.Error())
			tracking.LogErr("jwt", errorLog)
			response = utils.Message(false, "Malformed authentication token")
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response, http.StatusForbidden)
			return
		}

		if !token.Valid || tk.IsRevoke == true { //Token is invalid, maybe not signed on this server
			response = utils.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response, http.StatusForbidden)
			return
		}
		tmp := tk.VerifyExpiresAt(time.Now().Unix(), true)
		fmt.Print("Value expireAt: ", tmp)
		if !tmp {
			response = utils.Message(false, "Token is expired.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response, http.StatusForbidden)
			return
		}
		for _, value := range adminOnlyPath {
			if value == requestPath {
				if tk.Role != "admin" {
					response = utils.Message(false, "request path Forbidden.")
					w.WriteHeader(http.StatusForbidden)
					w.Header().Add("Content-Type", "application/json")
					utils.Respond(w, response, http.StatusForbidden)
					return
				}
			}
		}
		//everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		log.Default().Println("User ", tk.UserId) //Useful for monitoring
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}

func DecryptToken(tokenString string) (jwt.Claims, bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Token{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})
	if err != nil { //Malformed token, returns with http code 403 as usual
		tracking.LogErr("jwt", err)
		log.Default().Println("Malformed authentication token ", err)
		return token.Claims, token.Valid, err
	}

	log.Default().Println(token.Claims.(*model.Token).UserId)
	return token.Claims, token.Valid, nil
}