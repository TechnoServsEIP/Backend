package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-gomail/gomail"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var ports = [...]string{"25575", "25576", "25577", "25578", "25579", "26000"}

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}, httpCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(data)
}

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GetPort() string {
	return ports[seededRand.Intn(len(ports))]
}

func SendConfirmationEmail(url, to string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "jonathan.frickert@epitech.eu")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Confirm Email")
	m.SetBody("text/html", "Please click <a href="+url+">here</a> to confirm your email")

	d := gomail.NewDialer("in-v3.mailjet.com", 587, "56a430fb737fca0b6c5d33a449c6206e", "097628559f09a9bab73a6fab8b2d357d")

	return d.DialAndSend(m)
}

//func DecryptToken(token string) (jwt.StandardClaims, bool, error)  {
//	decryptToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, errors.New("error")
//		}
//		return []byte(os.Getenv("token_password")), nil
//	})
//	claims, valid := decryptToken.Claims.(jwt.StandardClaims)
//	if !decryptToken.Valid {
//		valid = false
//	}
//	return claims, valid, err
//}

func DecryptToken(tokenString string) (jwt.StandardClaims, bool, error) {
	decryptToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})
	if err != nil { //Malformed token, returns with http code 403 as usual
		fmt.Println("Malformed authentication token", err)
		return jwt.StandardClaims{}, false, errors.New("error")
	}
	claims, valid := decryptToken.Claims.(jwt.StandardClaims)
	if !decryptToken.Valid {
		valid = false
	}
	return claims, valid, nil

	//tokenParsed, err := jwt.ParseWithClaims(tokenString, jwt.StandardClaims{},
	//func(token *jwt.Token) (interface{}, error) {
	//	return []byte(os.Getenv("token_password")), nil
	//})
	//if err != nil { //Malformed token, returns with http code 403 as usual
	//	fmt.Println("Malformed authentication token")
	//	return tokenParsed.Claims.(jwt.StandardClaims), false, errors.New("error")
	//}
	//return tokenParsed.Claims.(jwt.StandardClaims), tokenParsed.Valid, nil
}
