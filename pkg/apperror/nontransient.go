package apperror

type NonTransient struct {
	Err error
}

func (n NonTransient) Error() string {
	return n.Err.Error()
}

func (n NonTransient) Unwrap() error {
	return n.Err
}
