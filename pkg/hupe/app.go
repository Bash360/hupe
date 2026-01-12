package hupe

import (
	"log"

	"github.com/bash360/hupe/internal/circuit"
	"github.com/bash360/hupe/internal/retry"
	hupeI "github.com/bash360/hupe/pkg/hupe/interface"
)

type Hupe struct {
	circuit hupeI.ICircuit
	retry   hupeI.IRetry
}

func New(options circuit.CircuitOptions) (*Hupe, error) {

	retry, err := retry.New(options.Operation)

	if err != nil {
		log.Fatalln("An Error occured", err.Error())
		return nil, err
	}

	circuit, err := circuit.New(&options)
	if err != nil {
		log.Fatalln("An Error occured", err.Error())
		return nil, err
	}

	return &Hupe{retry: retry, circuit: circuit}, nil
}

func (h *Hupe) Execute() ([]any, error) {
	return h.circuit.Execute()
}