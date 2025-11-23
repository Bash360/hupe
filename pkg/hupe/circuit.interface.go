package hupe

type ICircuit interface {
	CheckErrRate()
	SetThreshold() ICircuit
	SetTimeout() ICircuit
}
