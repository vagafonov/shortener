package customerror

import "errors"

// ErrURLAlreadyExists error for already exists url.
var ErrURLAlreadyExists = errors.New("url already exists")
