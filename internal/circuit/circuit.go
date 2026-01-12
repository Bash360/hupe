package circuit

import (
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bash360/hupe/internal/shared"
	"github.com/bash360/hupe/pkg/apperror"
	hupeI "github.com/bash360/hupe/pkg/hupe/interface"
	"github.com/bash360/hupe/pkg/utils"
)

const (
	Closed = iota
	Open
	HalfOpen
)

const defaultWindowSize = 10

type CircuitBreaker struct {
	state             int32
	threshold         float64
	slidingWindow     []error
	timeOut           time.Duration
	fallBack          *shared.Operation
	operation         *shared.Operation
	lastTrialAt       time.Time
	slidingWindowSize int
	halfOpenAttempts  int32
	retry             *hupeI.IRetry
	mux               sync.RWMutex
}

type CircuitOptions struct {
	Threshold         float64
	Timeout           time.Duration
	Operation         *shared.Operation
	SlidingWindowSize uint
	Fallback          *shared.Operation
	Retry             *hupeI.IRetry
}

func New(options *CircuitOptions) (*CircuitBreaker, error) {

	c := &CircuitBreaker{
		state:            Open,
		halfOpenAttempts: 0,
	}

	if options.SlidingWindowSize == 0 {
		options.SlidingWindowSize = defaultWindowSize
	} else {
		c.slidingWindowSize = int(options.SlidingWindowSize)
	}

	if options.Timeout == 0 {
		options.Timeout = time.Second * 10
	} else {
		c.setTimeout(int(options.Timeout))
	}

	if options.Threshold == 0 {
		options.Threshold = 0.5
	}
	if options.Operation == nil {
		return nil, errors.New("operation cannot be nil")
	}

	if options.Fallback != nil {
		err := c.setFallback(options.Fallback.Fn, options.Fallback.Args...)
		if err != nil {
			return nil, err
		}
	}

	err := c.setThreshold(options.Threshold)
	if err != nil {
		return nil, err
	}
	c.operation = options.Operation
	c.slidingWindow = make([]error, 0, options.SlidingWindowSize)

	return c, nil

}

func (c *CircuitBreaker) setState(state int) {

	atomic.StoreInt32(&c.state, int32(state))

}

func (c *CircuitBreaker) setThreshold(threshold float64) error {
	if threshold > 0 || threshold < 1 {
		return ErrThreshold
	}
	c.threshold = threshold
	return nil
}

func (c *CircuitBreaker) setFallback(fallback interface{}, args ...any) error {
	err := utils.ValidateArgs(fallback, args...)
	if err != nil {
		return err
	}
	c.fallBack.Fn = fallback
	c.fallBack.Args = args
	return nil
}

func (c *CircuitBreaker) setTimeout(millisecond int) {
	c.timeOut = time.Millisecond * time.Duration(millisecond)
}

func (c *CircuitBreaker) SetWindowSize(size int) {
	c.slidingWindowSize = size
}

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
	c.mux.RLock()
	defer c.mux.RUnlock()
	errCount := 0

	for _, v := range c.slidingWindow {
		if v != nil {
			errCount += 1
		}

	}
	return errCount/len(c.slidingWindow) > int(c.threshold)

}

func (c *CircuitBreaker) addError(err error) {
	if errors.Is(err, apperror.Transient{}) {
		c.addToSlidingWindow(err)
	} else {
		c.addToSlidingWindow(nil)
	}

	if c.checkThreshold() {
		c.setState(Open)
	}

}

func (c *CircuitBreaker) Execute() ([]any, error) {
	var payload []any
	var err error

	if time.Since(c.lastTrialAt) > c.timeOut && atomic.CompareAndSwapInt32(&c.state, int32(Open), HalfOpen) {
		atomic.StoreInt32(&c.halfOpenAttempts, 0)
	}

	switch atomic.LoadInt32(&c.state) {
	case Closed:
		retry := *c.retry
		payload, err = retry.Execute()
		c.addError(err)
		return payload, err
	case Open:
		payload, err = c.halfOpen()

	case HalfOpen:
		if time.Since(c.lastTrialAt) > c.timeOut && atomic.LoadInt32(&c.halfOpenAttempts) < 1 {
			payload, err = shared.Execute(shared.RetryPolicy{Count: 0, Delay: 0}, c.operation.Fn, c.operation.Args...)
			atomic.AddInt32(&c.halfOpenAttempts, 1)
			c.mux.Lock()
			defer c.mux.Unlock()
			if err == nil {
				c.setState(Closed)
				c.slidingWindow = make([]error, 0, c.slidingWindowSize)
				atomic.StoreInt32(&c.halfOpenAttempts, 0)

			} else {
				c.setState(Open)
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
func (c *CircuitBreaker) addToSlidingWindow(err error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if len(c.slidingWindow) < c.slidingWindowSize {
		c.slidingWindow = append(c.slidingWindow, err)

	} else {
		copy(c.slidingWindow[:], c.slidingWindow[1:])
		c.slidingWindow[c.slidingWindowSize-1] = err

	}

}
