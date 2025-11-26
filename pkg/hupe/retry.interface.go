package hupe

type IRetry interface {
	WithDelay(delay uint) IRetry

	WithCount(count uint) IRetry

	Execute() ([]any, error)
}
