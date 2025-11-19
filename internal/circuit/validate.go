package circuit

import "errors"

func validateFallback(callback interface{}) error {
	return errors.New("error")
}
