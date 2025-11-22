package retry

import (
	"errors"
	"reflect"
	"time"

	"github.com/bash360/hupe/pkg/apperror"
	"github.com/bash360/hupe/pkg/hupe"
	"github.com/bash360/hupe/pkg/utils"
)

type Retry struct {
	delay time.Duration
	count uint
	fn    interface{}
	args  []any
}

func New(fn interface{}, args ...any) (*Retry, error) {

	err := validateFunc(fn)
	if err != nil {
		return nil, err
	}

	err = validateArgs(fn, args...)

	if err != nil {
		return nil, err
	}

	return &Retry{
		fn:    fn,
		delay: time.Millisecond * 500,
		count: 4,
		args:  args,
	}, nil
}

func (r *Retry) SetDelay(delay uint) hupe.IRetry {
	r.delay = time.Millisecond * time.Duration(delay)
	return r
}

func (r *Retry) SetCount(count uint) hupe.IRetry {
	r.count = count
	return r
}

func (r *Retry) Execute() ([]any, error) {

	operation := reflect.ValueOf(r.fn)
	args := make([]reflect.Value, 0)
	if r.args != nil {
		for _, v := range r.args {
			args = append(args, reflect.ValueOf(v))
		}
	}

	var err error
	returnValues := make([]reflect.Value, 0)

	for i := 0; i <= int(r.count); i++ {

		returnValues = operation.Call(args)
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
			time.Sleep(r.delay)
		}

	}

	payload := utils.ValueToInterface(returnValues[:len(returnValues)-1])

	return payload, err

}
