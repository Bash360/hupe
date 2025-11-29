package hupeI

type ICircuit interface {
	CheckErrRate()
	Fallback() []any
	SetState(state int)
	AddError(err error)
}
