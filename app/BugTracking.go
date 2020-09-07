package app

import (
	"fmt"
	"github.com/TechnoServsEIP/Backend/models"
	"log"
	"os"
)
func checkIfExist(path string, err error) {
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return
		}
		defer file.Close()
	}
}
func CreateLogFile() {
	for _, log := range models.LogType() {
		path := log +".log"
		_, err := os.Stat(path)
		checkIfExist(path, err)
	}
}
func writeError(terr string, err error) {
	file, err := os.OpenFile("LogFile/"+terr+".log", os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()
	log.SetOutput(file)
	log.Print(err.Error())

}

func LogErr(terr string, err error) {
	PublishError(terr, err)
	CreateLogFile()
	writeError(terr, err)
}
