package hupe

import "github.com/bash360/hupe/internal/retry"

type IRetry interface {
	SetInterval(interval uint) *retry.Retry

	SetCount(count uint) *retry.Retry

	Execute() ([]any, error)
}
