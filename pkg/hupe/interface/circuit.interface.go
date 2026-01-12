package hupeI

type ICircuit interface {
	Execute() ([]any, error)
}
