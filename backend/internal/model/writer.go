package model

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
