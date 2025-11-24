package circuit

import (
	"time"

	"github.com/bash360/hupe/pkg/utils"
)

const (
	Closed = iota
	Open
	HalfOpen
)

type CircuitBreaker struct {
	state            int
	threshold        float64
	slidingWindow    []error
	timeOut          time.Duration
	fallBack         interface{}
	HalfOpenAttempts uint
	LastTrialAt      time.Time
	args             []any
}

func New() *CircuitBreaker {

	return &CircuitBreaker{
		state:            Open,
		threshold:        0.5,
		slidingWindow:    make([]error, 0, 10),
		timeOut:          time.Duration(time.Second * 10),
		HalfOpenAttempts: 2,
	}

}

func (c *CircuitBreaker) SetState(state int) *CircuitBreaker {

	c.state = state
	return c

}

func (c *CircuitBreaker) SetThreshold(threshold float64) error {
	if threshold > 0 || threshold < 1 {
		return ErrThreshold
	}
	c.threshold = threshold
	return nil
}

func (c *CircuitBreaker) SetFallback(fallback interface{}, args ...any) error {
	err := utils.ValidateArgs(fallback, args...)
	if err != nil {
		return err
	}

	return nil
}
