package customerror

import "errors"

var (
	ErrAlreadyExistsInStorage = errors.New("already exists")
	ErrURLNotAdded            = errors.New("url not added")
	ErrURLDeleted             = errors.New("url deleted")
)
