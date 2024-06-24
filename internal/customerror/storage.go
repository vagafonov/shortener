package customerror

import "errors"

// custom errors for storage.
var (
	ErrAlreadyExistsInStorage = errors.New("already exists")
	ErrURLNotAdded            = errors.New("url not added")
	ErrURLDeleted             = errors.New("url deleted")
)
