package hupeI

type IHupe interface {
	WithTimeout(timeout int) IHupe
	WithFallback(fn interface{}, args ...any) error
	WithDelay(delay int) IHupe
	WithMaxRetries(retryCount int) IHupe
	WithErrThreshold(threshold float64) IHupe
	WithWindowSize(size int) IHupe
}

/*
timeout
operation
fallback

count
delay
threshold
slidingwindowsize
*/
