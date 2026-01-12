package circuit

import (
	"testing"
	"time"

	"github.com/bash360/hupe/internal/shared"
	hupeI "github.com/bash360/hupe/pkg/hupe/interface"
)

func TestCircuit(t *testing.T) {
type input struct {
		Threshold         float64
	Timeout           time.Duration
	Operation         *shared.Operation
	SlidingWindowSize uint
	Fallback          *shared.Operation
	Retry             *hupeI.IRetry
	}

	type want struct{
		err error
		circuit *CircuitBreaker
	}

	tests := []struct {
		name    string
		input   input
		want    want
		wantErr bool
	}{
		{
			"Test Circuit Constructor with valid inputs",
			input{
				Threshold:         0.6,
				Timeout:           time.Duration(time.Second * 5),
				Operation:         &shared.Operation{},
				SlidingWindowSize: 15,
				Fallback:          &shared.Operation{},
				Retry:             nil,
			},
			want{
				err:     nil,
				circuit: &CircuitBreaker{},
			},
			false,
		},
	}


}