package apperror

type Transient struct {
	err error
}

func (t Transient) Error() string {
	return t.err.Error()
}

func (t Transient) Unwrap() error {
	return t.err
}
