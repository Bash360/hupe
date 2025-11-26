package retry

import (
	"time"

	"github.com/bash360/hupe/internal/shared"
	"github.com/bash360/hupe/pkg/hupe"
	"github.com/bash360/hupe/pkg/utils"
)

type Retry struct {
	delay     time.Duration
	count     uint
	operation *shared.Operation
}

func New(operation *shared.Operation) (*Retry, error) {

	err := utils.ValidateFunc(operation.Fn)
	if err != nil {
		return nil, err
	}

	err = utils.ValidateArgs(operation.Fn, operation.Args...)

	if err != nil {
		return nil, err
	}

	return &Retry{
		operation: operation,
		delay:     time.Millisecond * 500,
		count:     4,
	}, nil
}

func (r *Retry) WithDelay(millisecond uint) hupe.IRetry {
	r.delay = time.Millisecond * time.Duration(millisecond)
	return r
}

func (r *Retry) WithCount(count uint) hupe.IRetry {
	r.count = count
	return r
}

func (r *Retry) Execute() ([]any, error) {

	payload, err := shared.Execute(shared.RetryPolicy{Delay: r.delay, Count: int(r.count)}, r.operation.Fn, r.operation.Args...)

	return payload, err

}
