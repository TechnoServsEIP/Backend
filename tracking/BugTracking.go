package tracking

import (
	"log"
	"os"
)

func writeError(terr string, err error) {
	log.Default().Println(err)
	path := "LogFile/" + terr + ".log"
	file, errFile := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if errFile != nil {
		log.Default().Println(errFile.Error())
	}
	if _, err := file.WriteString(err.Error() + "\n"); err != nil {
		log.Println(err)
	}
	file.Close()
}

func LogErr(terr string, err error) {
	PublishError(terr, err)
	writeError(terr, err)
}
