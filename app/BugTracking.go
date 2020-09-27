package app

import (
	"fmt"
	"os"

	"github.com/TechnoServsEIP/Backend/models"
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
		path := log + ".log"
		_, err := os.Stat(path)
		checkIfExist(path, err)
	}
}
func writeError(terr string, err error) {
	fmt.Println(err)
	file, errFile := os.OpenFile("LogFile/"+terr+".log", os.O_CREATE|os.O_APPEND, 0644)
	if errFile != nil {
		fmt.Println(errFile.Error())
	}
	d1 := []byte(err.Error())
	_, err = file.Write(d1)
	// log.SetOutput(file)
	// log.Print(err.Error())
	file.Close()
}

func LogErr(terr string, err error) {
	PublishError(terr, err)
	writeError(terr, err)
}
