package service

import "errors"

var (
	ErrUnavailable      = errors.New("unavailable")
	ErrBookNotDownload  = errors.New("book not downloaded")
	ErrBookFileNotFound = errors.New("book file not found")
	ErrInvalidBookID    = errors.New("invalid book id")
	ErrInvalidHashCode  = errors.New("invalid hash code")
)
