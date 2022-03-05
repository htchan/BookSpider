package logging

import (
	"errors"
	"fmt"
	"log"
	"os"
	"github.com/htchan/BookSpider/internal/utils"
)

var logLevel = 0

var f *os.File

func write(level int, content string) {
	if level < logLevel { return }
	//TODO: write the string to screen
	log.Print(content)
}

func UseFile(filename string) {
	//TODO: open file with filename
	var err error
	f, err = os.OpenFile(filename, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	utils.CheckError(err)
	log.SetOutput(f)
}

func SetLogLevel(level int) error {
	if level < 0 || level > 3 {
		return errors.New("invalid log level, only accept interger between 0 to 3")
	}
	logLevel = level
	return nil
}

func Debug(content string, values ...interface{}) {
	content = fmt.Sprintf(content, values...)
	write(0, content)
}

func Info(content string, values ...interface{}) {
	content = fmt.Sprintf(content, values...)
	write(1, content)
}

func Warn(content string, values ...interface{}) {
	content = fmt.Sprintf(content, values...)
	write(2, content)
}

func Error(content string, values ...interface{}) {
	content = fmt.Sprintf(content, values...)
	write(3, content)
}