package apperror

import "errors"

var ErrNoReturn = errors.New("last or only return value must be an error type")

var ErrInvalidReturn = errors.New("last or only return value must be an error type")
var ErrNotAFunction = errors.New("not a function")

var ErrArgumentSize = errors.New("operation parameter and argument num mismatch")

var ErrUnassignableArgument = errors.New("argument is unassignable to parameter of the function")
