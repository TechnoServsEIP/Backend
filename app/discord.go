package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func PublishError(logFile string, errorMessage error) {
	requestBody, err := json.Marshal(map[string]string{
		"message": errorMessage.Error(),
		"logFile": logFile,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = http.Post("http://serverlog:3000/log", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println(err.Error())
	}
}
