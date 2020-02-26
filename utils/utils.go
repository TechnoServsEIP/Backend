package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"time"
)

var ports = [...]string{"25575", "25576", "25577", "25578", "25579", "26000"}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Message(status bool, message string) (map[string]interface{}) {
	return map[string]interface{} {"status" : status, "message" : message}
}

func Respond(w http.ResponseWriter, data map[string] interface{}, httpCode int)  {
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

func TestHostPort(host string) (string, error) {
	for _, port := range ports {
		timeout := time.Second

		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
		if err != nil {
			fmt.Println("Connecting error:", err)
		}
		if conn != nil {
			defer conn.Close()
			fmt.Println("Opened", net.JoinHostPort(host, port))
			return port, nil
		}
	}
	return "", errors.New("Can't find available port on host, " + host)
}