package hupeI

type ICircuit interface {
	SetFallback() []any
	SetState(state int)
}
