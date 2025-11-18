package circuit

import "time"

const (
	Closed = iota
	Open
	HalfOpen
)

type CircuitBreaker struct {
	state         int
	threshold     float64
	slidingWindow []error
	timeOut       time.Duration
	fallBack      interface{}
}

func New(fallBack interface{}, args ...interface{}) *CircuitBreaker {

	return &CircuitBreaker{
		state:         Open,
		threshold:     0.5,
		slidingWindow: make([]error, 0, 10),
		timeOut:       time.Duration(time.Second * 10),
		fallBack:      fallBack,
	}

}
