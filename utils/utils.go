package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
	"net"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-gomail/gomail"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var ports = []string{"25576", "25577", "25578", "25579", "26000", "26001"}
var portsBinded = []string{}

func checkBindedPort(port string) bool {
	_, err := net.Listen("tcp", ":" + port)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't listen on port %q: %s\n", port, err)
		return false
	}

	fmt.Printf("TCP Port %q is available", port)
	return true
}

func ReOrderPorts(allPortsBinded []string) {
	portsBinded = allPortsBinded

	if (len(portsBinded) > 0) {
		for i := 0; i < len(ports); i++ {
			_, res := Find(portsBinded, ports[i])

			if (res) {
				ports[i] = ports[len(ports)-1]
				ports[len(ports)-1] = ""
				ports = ports[:len(ports)-1]
				i--
			}
		}
	}
	fmt.Println("ports binded: ", portsBinded)
	fmt.Println("ports: ", ports)
}

func Find(slice []string, val string) (int, bool) {
    for i, item := range slice {
        if item == val {
            return i, true
        }
    }
    return -1, false
}

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}, httpCode int) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods",
		"GET, POST, PATCH, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, "+
			"X-CSRF-Token, Authorization, Origin, X-Auth-Token")
	w.Header().Set("Access-Control-Expose-Headers",
		"Authorization")
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

func FreeThePort(portToFree string) {
	i, res := Find(portsBinded, portToFree)

	if (res) {
		ports = append(ports, portToFree)

		portsBinded[i] = portsBinded[len(portsBinded)-1]
		portsBinded[len(portsBinded)-1] = ""
		portsBinded = portsBinded[:len(portsBinded)-1]
	}

	fmt.Println("ports binded: ", portsBinded)
	fmt.Println("ports: ", ports)
}

func GetPort() string {
	if (len(ports) > 0) {
		portState := checkBindedPort(ports[0])

		if (portState) {
			/*
			* Get the first port of the ports slice
			*/
			portsBinded = append(portsBinded, ports[0])
			fmt.Println("ports binded: ", portsBinded)
		
			/*
			* Delete the first port of the ports slice
			*/
			portToSend := ports[0]
			ports = ports[1:len(ports)]
			fmt.Println("ports available: ", ports)

			return portToSend;
		} else {
			fmt.Println("no available ports")
			return "no port available";
		}
	}

	return "no port available";
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
