package apperror

type Transient struct {
	Err error
}

func (t Transient) Error() string {
	return t.Err.Error()
}

func (t Transient) Unwrap() error {
	return t.Err
}
