package utils

import (
	"reflect"

	"github.com/bash360/hupe/pkg/apperror"
)

func ValidateFunc(function interface{}) error {
	var functionValue reflect.Value = reflect.ValueOf(function)
	var functionType reflect.Type = functionValue.Type()

	if functionValue.IsValid() && functionValue.Kind() != reflect.Func {
		return apperror.ErrNotAFunction
	}

	numOut := functionType.NumOut()
	if numOut == 0 {
		return apperror.ErrNoReturn
	}

	var errorT reflect.Type = functionType.Out(numOut - 1)

	errorType := reflect.TypeOf(new(error)).Elem()

	if !errorT.Implements(errorType) {
		return apperror.ErrInvalidReturn
	}

	return nil
}

func ValidateArgs(operation any, args ...any) error {

	opV := reflect.ValueOf(operation)
	opT := opV.Type()
	if opV.IsValid() && opV.Kind() != reflect.Func {
		return apperror.ErrNotAFunction
	}
	if opT.NumIn() != len(args) {
		return apperror.ErrArgumentSize
	}

	for i, arg := range args {
		actual := reflect.TypeOf(arg)
		if !(actual.AssignableTo(opT.In(i))) {
			return apperror.ErrUnassignableArgument
		}

	}

	return nil
}
