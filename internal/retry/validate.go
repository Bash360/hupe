package retry

import (
	"reflect"

	"github.com/bash360/hupe/pkg/apperror"
)

func validateFunc(function interface{}) error {
	var functionValue reflect.Value = reflect.ValueOf(function)
	var functionType reflect.Type = functionValue.Type()

	if functionValue.IsValid() && functionValue.Kind() != reflect.Func {
		return apperror.ErrNotAFunction
	}

	numOut := functionType.NumOut()
	if numOut == 0 {
		return ErrNoReturn
	}

	var errorT reflect.Type = functionType.Out(numOut - 1)

	errorType := reflect.TypeOf(new(error)).Elem()

	if !errorT.Implements(errorType) {
		return ErrInvalidReturn
	}

	return nil
}
