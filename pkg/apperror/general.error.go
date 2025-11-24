package apperror

import "errors"

var ErrNotAFunction = errors.New("not a function")

var ErrArgumentSize = errors.New("operation parameter and argument num mismatch")

var ErrUnassignableArgument = errors.New("argument is unassignable to parameter of the function")
