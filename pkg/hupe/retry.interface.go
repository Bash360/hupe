package hupe

type IRetry interface {
	SetInterval(interval uint) IRetry

	SetCount(count uint) IRetry

	Execute() ([]any, error)
}
