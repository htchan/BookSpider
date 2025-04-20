package service

import "errors"

var (
	ErrUnavailable           = errors.New("unavailable")
	ErrBookStatusNotError    = errors.New("book status is not error")
	ErrBookStatusNotEnd      = errors.New("book status is not end")
	ErrBookNotDownload       = errors.New("book not downloaded")
	ErrBookAlreadyDownloaded = errors.New("book was downloaded")
	ErrBookFileNotFound      = errors.New("book file not found")
	ErrInvalidBookID         = errors.New("invalid book id")
	ErrInvalidHashCode       = errors.New("invalid hash code")
	ErrTooManyFailedChapters = errors.New("too many failed chapters")
)
