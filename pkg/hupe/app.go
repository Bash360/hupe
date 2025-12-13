package hupe

import (
	"log"

	"github.com/bash360/hupe/internal/circuit"
	"github.com/bash360/hupe/internal/retry"
	"github.com/bash360/hupe/internal/shared"
	hupeI "github.com/bash360/hupe/pkg/hupe/interface"
)

type Hupe struct {
	circuit hupeI.ICircuit
	retry   hupeI.IRetry
	
}

func New(fn interface{}, args ...any) *Hupe {

	retry, err:=retry.New(&shared.Operation{Fn: fn, Args: args})

	if err!=nil {
		log.Fatalln("An Error occured", err.Error())
	}

	return &Hupe{retry: retry, circuit: circuit.New()}
}

/*
timeout
operation
fallback

count
delay
threshold
slidingwindowsize
*/
