package retry

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bash360/hupe/pkg/hupe"
)

func TestConstructor(t *testing.T) {
	type input struct {
		args      []any
		operation interface{}
	}

	type want struct {
		err   error
		retry hupe.IRetry
	}
	tests := []struct {
		name    string
		input   input
		want    want
		wantErr bool
	}{
		{
			"Test constructor with wrong type input",
			input{nil, "wrong input"},
			want{ErrNotAFunction, nil},
			true,
		},
		{
			"Test constructor with correct implementing function but no return type",
			input{[]any{"mark bashir"}, func(name string) {
				fmt.Println(name)
			}},
			want{ErrNoReturn, nil},
			true,
		},
		{
			"Test constructor with correct implementing function and return type",
			input{nil, func() error {

				return errors.New("returns an error")
			}},
			want{nil, &Retry{}},
			false,
		},
	}

	for _, v := range tests {
		_, err := New(v.input.operation, v.input.args...)

		if !errors.Is(v.want.err, err) {
			t.Fatalf("constructor  %s", v.name)
		}

		t.Logf("constructor  %s", v.name)

	}
}
