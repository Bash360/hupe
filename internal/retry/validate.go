package retry

import (
	"log"
	"reflect"
)

func validateFunc(function interface{}) error {
	var functionValue reflect.Value = reflect.ValueOf(function)
	var functionType reflect.Type = functionValue.Type()

	if functionValue.Kind() != reflect.Func {
		return ErrNotAFunction
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

func validateArgs(operation any, args ...any) error {

	opV := reflect.ValueOf(operation)
	opT := opV.Type()
	if opT.NumIn() != len(args) {
		return ErrArgumentSize
	}

	for i, arg := range args {
		actual := reflect.TypeOf(arg)
		if !(actual.AssignableTo(opT.In(i))) {
			return ErrUnassignableArgument
		}

	}

	return nil
}
