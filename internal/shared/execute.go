package shared

import (
	"errors"
	"reflect"
	"time"

	"github.com/bash360/hupe/pkg/apperror"
	"github.com/bash360/hupe/pkg/utils"
)

func Execute(retryPolicy RetryPolicy, fn interface{}, args ...any) ([]any, error) {

	operation := reflect.ValueOf(fn)
	argsV := make([]reflect.Value, len(args))

	for i, v := range args {
		argsV[i] = reflect.ValueOf(v)
	}

	var err error
	returnValues := make([]reflect.Value, 0)

	for i := 0; i <= int(retryPolicy.Count); i++ {

		returnValues = operation.Call(argsV)
		returnedErr := returnValues[len(returnValues)-1].Interface()

		if returnedErr != nil {
			err = returnedErr.(error)
		} else {
			err = nil
		}

		if err == nil || errors.As(err, &apperror.NonTransient{}) {
			break
		}

		if errors.As(err, &apperror.Transient{}) {
			time.Sleep(retryPolicy.Delay)
		}

	}

	payload := utils.ValueToInterface(returnValues[:len(returnValues)-1])

	return payload, err

}
