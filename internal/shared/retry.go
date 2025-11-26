package shared

import "time"

type RetryPolicy struct {
	Delay time.Duration
	Count int
}
