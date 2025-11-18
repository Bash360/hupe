package retry

import (
	"errors"
	"reflect"
)

func validateFunc(function interface{}) error {
	var functionValue reflect.Value = reflect.ValueOf(function)
	var functionType reflect.Type = functionValue.Type()

	if functionValue.Kind() != reflect.Func {
		return errors.New("not a function")
	}

	numOut := functionType.NumOut()
	if numOut == 0 {
		return errors.New("last or only return value must be an error type")
	}

	var errorT reflect.Type = functionType.Out(numOut - 1)

	errorType := reflect.TypeOf(new(error)).Elem()

	if !errorT.Implements(errorType) {
		return errors.New("last or only return value must be an error type")
	}

	return nil
}
