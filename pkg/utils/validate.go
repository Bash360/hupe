package utils

import (
	"reflect"

	"github.com/bash360/hupe/pkg/apperror"
)

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
