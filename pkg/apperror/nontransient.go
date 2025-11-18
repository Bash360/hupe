package apperror

type NonTransient struct {
	err error
}

func (n NonTransient) Error() string {
	return n.err.Error()
}

func (n NonTransient) Unwrap() error {
	return n.err
}
