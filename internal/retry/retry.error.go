package retry

import "errors"

var ErrNoReturn = errors.New("last or only return value must be an error type")

var ErrInvalidReturn = errors.New("last or only return value must be an error type")
