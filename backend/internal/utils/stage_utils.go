package utils

import (
	"log"
	"os"
)

var StageFileName string

func WriteStage(s string) {
	file, err := os.OpenFile(StageFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		log.Println(err, "\n", StageFileName)
	}
	file.WriteString(s + "\n")
	file.Close()
}
