package hupe

type IRetry interface {
	SetDelay(delay uint) IRetry

	SetCount(count uint) IRetry

	Execute() ([]any, error)
}
