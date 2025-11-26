package retry

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/bash360/hupe/internal/shared"
	"github.com/bash360/hupe/pkg/apperror"
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
			want{apperror.ErrNotAFunction, nil},
			true,
		},
		{
			"Test constructor with correct implementing function but no return type",
			input{[]any{"mark bashir"}, func(name string) {
				fmt.Println(name)
			}},
			want{apperror.ErrNoReturn, nil},
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

		{
			"Test constructor with correct implementing function and return type but mismatching argument and parameters",
			input{[]any{"tobi", 16}, func(name string) error {

				return errors.New("returns an error")
			}},
			want{apperror.ErrArgumentSize, &Retry{}},
			true,
		},

		{
			"Test constructor with correct implementing function and return type but mismatching argument and parameters type",
			input{[]any{"tobi", 16}, func(a int, b int) error {

				return errors.New("returns an error")
			}},
			want{apperror.ErrUnassignableArgument, &Retry{}},
			true,
		},
	}

	for _, v := range tests {
		_, err := New(&shared.Operation{Fn: v.input.operation, Args: v.input.args})

		if !errors.Is(v.want.err, err) {
			t.Errorf("constructor Error %s", v.name)
			continue
		}

	}
}

func TestWithDelay(t *testing.T) {
	r, err := New(&shared.Operation{Fn: func() error { return errors.New("dummy error") }})

	if err != nil {
		t.Errorf("Set Interval test failed %s", err.Error())
		return
	}
	r.WithDelay(300)

	if r.delay != time.Millisecond*time.Duration(300) {
		t.Error("Set interval test method is not working properly ")
		return
	}

}

func TestWithCount(t *testing.T) {
	r, err := New(&shared.Operation{Fn: func() error { return errors.New("dummy error") }})

	if err != nil {
		t.Errorf("Set Count test failed %s", err.Error())
		return
	}

	r.WithCount(5)

	if r.count != 5 {
		t.Error("Set Count error: count specified with set count and value in struct differ")
		return
	}

}

func TestExecute(t *testing.T) {

	fn := func() func() (string, error) {
		count := 0
		var err error = nil
		var result string = ""
		return func() (string, error) {

			switch count {
			case 0:
				err = apperror.Transient{Err: errors.New("Server not ready 1")}
			case 1:
				err = apperror.Transient{Err: errors.New("Server not ready 2")}
			case 2:
				err = apperror.Transient{Err: errors.New("Server not ready 3")}
			case 3:
				err = apperror.Transient{Err: errors.New("Server not ready 4")}
			default:
				err = nil
				result = "success"

			}
			count++

			return result, err

		}
	}

	r, err := New(&shared.Operation{Fn: fn()})

	if err != nil {
		t.Log("Test Execute failing: ", err.Error())
	}

	payload, err := r.Execute()

	if err != nil {

		t.Errorf("Test Execute failing: %s", err.Error())
	}

	t.Log(payload)

}
