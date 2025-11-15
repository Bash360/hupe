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
	slidingWindow []bool
	timeOut       time.Duration
	fallBack func 
	
}

func NewCircuitBreaker() *CircuitBreaker {

	return &CircuitBreaker{
		state:         Open,
		threshold:     0.5,
		slidingWindow: make([]bool, 10),
		timeOut:       time.Duration(time.Second * 10),
	}

}
