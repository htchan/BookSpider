package vendor

import (
	"errors"
	"fmt"
)

type StatusCodeError struct {
	StatusCode int
}

var _ error = StatusCodeError{}

func (err StatusCodeError) Error() string {
	return fmt.Sprintf("status code is %d", err.StatusCode)
}

var (
	ErrTimeout = errors.New("request timeout")

	ErrFieldsNotFound = errors.New("book fields not found")
	// book fields not found error
	ErrBookTitleNotFound   = errors.New("title not found")
	ErrBookWriterNotFound  = errors.New("writer not found")
	ErrBookTypeNotFound    = errors.New("type not found")
	ErrBookDateNotFound    = errors.New("date not found")
	ErrBookChapterNotFound = errors.New("chapter not found")
	// chapter list not found error
	ErrChapterListUrlNotFound   = errors.New("url not found")
	ErrChapterListTitleNotFound = errors.New("title not found")
	ErrChapterListEmpty         = errors.New("empty chapter list")
	// chapter fields not found error
	ErrChapterTitleNotFound   = errors.New("chapter title not found")
	ErrChapterContentNotFound = errors.New("chapter content not found")
)
