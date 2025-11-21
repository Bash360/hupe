package retry

import (
	"errors"
	"fmt"
	"testing"
	"time"

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

		{
			"Test constructor with correct implementing function and return type but mismatching argument and parameters",
			input{[]any{"tobi", 16}, func(name string) error {

				return errors.New("returns an error")
			}},
			want{ErrArgumentSize, &Retry{}},
			true,
		},

		{
			"Test constructor with correct implementing function and return type but mismatching argument and parameters type",
			input{[]any{"tobi", 16}, func(a int, b int) error {

				return errors.New("returns an error")
			}},
			want{ErrUnassignableArgument, &Retry{}},
			true,
		},
	}

	for _, v := range tests {
		_, err := New(v.input.operation, v.input.args...)

		if !errors.Is(v.want.err, err) {
			t.Errorf("constructor Error %s", v.name)
			continue
		}

	}
}

func TestSetInterval(t *testing.T) {
	r, err := New(func() error { return errors.New("dummy error") })

	if err != nil {
		t.Errorf("Set Interval test failed %s", err.Error())
		return
	}
	r.SetInterval(300)

	if r.interval != time.Millisecond*time.Duration(300) {
		t.Error("Set interval test method is not working properly ")
		return
	}

}

func TestSetCount(t *testing.T) {
	r, err := New(func() error { return errors.New("dummy error") })

	if err != nil {
		t.Errorf("Set Count test failed %s", err.Error())
		return
	}

	r.SetCount(5)

	if r.count != 5 {
		t.Error("Set Count error: count specified with set count and value in struct differ")
		return
	}

}

func TestExecute(t *testing.T) {}
