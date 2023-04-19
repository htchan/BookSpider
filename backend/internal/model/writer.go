package model

import (
	"fmt"
	"log"
	"strings"
)

type Writer struct {
	ID   int
	Name string
}

func NewWriter(name string) Writer {
	return Writer{
		ID:   -1,
		Name: name,
	}
}

func (w Writer) Checksum() string {
	if len(w.Name) > 100 {
		log.Printf("id: %v; name: %v is too long", w.ID, w.Name)
		return ""
	}

	writerName := strings.ReplaceAll(w.Name, " ", "")
	return strToShortHex(simplified(fmt.Sprintf("%s", writerName)))
}
