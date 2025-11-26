package circuit

import (
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bash360/hupe/internal/shared"
	"github.com/bash360/hupe/pkg/apperror"
	"github.com/bash360/hupe/pkg/hupe"
	"github.com/bash360/hupe/pkg/utils"
)

const (
	Closed = iota
	Open
	HalfOpen
)

var mux sync.Mutex

const defaultWindowSize = 10

type CircuitBreaker struct {
	state             int32
	threshold         float64
	slidingWindow     []error
	timeOut           time.Duration
	fallBack          *shared.Operation
	operation         *shared.Operation
	lastTrialAt       time.Time
	slidingWindowSize uint
	halfOpenAttempts  int32
	retry             *hupe.IRetry
}

type CircuitOptions struct {
	Threshold         float64
	Timeout           time.Duration
	Operation         *shared.Operation
	SlidingWindowSize uint
	Fallback          *shared.Operation
	Retry             *hupe.IRetry
}

func New() *CircuitBreaker {

	return &CircuitBreaker{
		state:             Open,
		threshold:         0.5,
		slidingWindow:     make([]error, 0, defaultWindowSize),
		timeOut:           time.Duration(time.Second * 10),
		slidingWindowSize: defaultWindowSize,
		halfOpenAttempts:  0,
	}

}

func (c *CircuitBreaker) setState(state int) {

	atomic.StoreInt32(&c.state, int32(state))

}

func (c *CircuitBreaker) SetThreshold(threshold float64) error {
	if threshold > 0 || threshold < 1 {
		return ErrThreshold
	}
	c.threshold = threshold
	return nil
}

// func (c *CircuitBreaker) SetFallback(fallback interface{}, args ...any) error {
// 	err := utils.ValidateArgs(fallback, args...)
// 	if err != nil {
// 		return err
// 	}
// 	c.fallBack = fallback
// 	c.args = args
// 	return nil
// }

// func (c *CircuitBreaker) SetTimeout(millisecond uint) {
// 	c.timeOut = time.Millisecond * time.Duration(millisecond)
// }

func (c *CircuitBreaker) runFallback() []any {
	fallbackV := reflect.ValueOf(c.fallBack.Fn)
	args := make([]reflect.Value, len(c.fallBack.Args))

	for i, v := range c.fallBack.Args {
		args[i] = reflect.ValueOf(v)
	}

	out := fallbackV.Call(args)

	output := utils.ValueToInterface(out)

	return output
}

func (c *CircuitBreaker) checkThreshold() bool {
	errCount := 0

	for _, v := range c.slidingWindow {
		if v != nil {
			errCount += 1
		}

	}

	return errCount/len(c.slidingWindow) > int(c.threshold)

}

func (c *CircuitBreaker) AddError(err error) {
	if errors.As(err, &apperror.Transient{}) {
		addToSlidingWindow(&c.slidingWindow, err, int(c.slidingWindowSize))
	} else {
		addToSlidingWindow(&c.slidingWindow, nil, int(c.slidingWindowSize))
	}

	if threshold := c.checkThreshold(); threshold {
		c.setState(Open)
	}

}

func (c *CircuitBreaker) Execute() ([]any, error) {
	var payload []any
	var err error
	if c.state == Open && time.Since(c.lastTrialAt) > c.timeOut {
		c.setState(HalfOpen)
		atomic.StoreInt32(&c.halfOpenAttempts, 0)
	}

	switch c.state {
	case Closed:
		rety := *c.retry
		payload, err = rety.Execute()
		c.AddError(err)
		return payload, err
	case Open:
		payload, err = c.halfOpen()

	case HalfOpen:
		if time.Since(c.lastTrialAt) > c.timeOut && c.halfOpenAttempts < 1 {
			payload, err = shared.Execute(shared.RetryPolicy{Count: 0, Delay: 0}, c.operation.Fn, c.operation.Args...)
			atomic.AddInt32(&c.halfOpenAttempts, 1)
			if err == nil {
				c.setState(Closed)
				mux.Lock()
				defer mux.Unlock()
				c.slidingWindow = make([]error, 0, c.slidingWindowSize)
				c.halfOpenAttempts = 0

			} else {
				c.setState(Open)
				mux.Lock()
				defer mux.Unlock()
				c.lastTrialAt = time.Now().Add(c.timeOut)

				payload, err = c.halfOpen()

			}

		}

	}
	return payload, err
}

func (c *CircuitBreaker) halfOpen() ([]any, error) {
	if c.fallBack != nil {
		payload := c.runFallback()
		return payload, nil
	} else {
		return nil, c.slidingWindow[len(c.slidingWindow)-1]
	}
}
func addToSlidingWindow(slidingWindow *[]error, err error, windowSize int) {
	mux.Lock()
	defer mux.Unlock()
	window := *slidingWindow
	if len(window) < windowSize {
		*slidingWindow = append(window, err)

	} else {
		copy(window[:], window[1:])
		window[windowSize-1] = err

	}

}
