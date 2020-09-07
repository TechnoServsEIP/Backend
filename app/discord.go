package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)


func PublishError(logFile string, errorMessage error) {
	requestBody, err := json.Marshal(map[string]string {
		"message" : errorMessage.Error(),
		"logFile" : logFile,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = http.Post("http://localhost:3000/", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println(err.Error())
	}
}
