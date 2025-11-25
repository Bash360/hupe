package retry

import (
	"time"

	"github.com/bash360/hupe/internal/shared"
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

	err = utils.ValidateArgs(fn, args...)

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

func (r *Retry) WithDelay(millisecond uint) hupe.IRetry {
	r.delay = time.Millisecond * time.Duration(millisecond)
	return r
}

func (r *Retry) WithCount(count uint) hupe.IRetry {
	r.count = count
	return r
}

func (r *Retry) Execute() ([]any, error) {

	payload, err := shared.Execute(shared.RetryPolicy{Delay: r.delay, Count: int(r.count)}, r.fn, r.args...)

	return payload, err

}
