package storage

import "errors"

// var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")

type Storage interface {
	GetByHash(key string) string
	GetByValue(key string) string
	Set(key string, val string) error
}
