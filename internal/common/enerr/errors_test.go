// Package provides error implementation for engine

package enerr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_errsToMap(t *testing.T) {
	type args struct {
		errs []string
	}
	tests := []struct {
		name string
		args args
		want map[string]error
	}{
		{name: "Simple", args: args{
			errs: []string{"Имя", "Должно содержать буквы"},
		}, want: map[string]error{"Имя": errors.New("Должно содержать буквы")}},
		{name: "More inputs", args: args{
			errs: []string{"Имя", "Должно содержать буквы", "Код", "Не должен быть пустым"},
		}, want: map[string]error{"Имя": errors.New("Должно содержать буквы"), "Код": errors.New("Не должен быть пустым")}},
		{name: "Non-valid input", args: args{
			errs: []string{"Имя", "Должно содержать буквы", "Код"},
		}, want: map[string]error{"Имя": errors.New("Должно содержать буквы"), "Код": errors.New("%not-defined%")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errsToMap(tt.args.errs)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEMNoArgs(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("E() did not panic")
		}
	}()
	_ = EM()
}

func TestEM(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		args args
		want *ApplicationError
	}{
		{
			name: "Simple Create",
			args: args{
				args: []interface{}{
					"Код", "Неправильный",
				},
			},
			want: &ApplicationError{
				Op:   "",
				User: "",
				Kind: 0,
				Err: &MultiError{
					Errs: map[string]error{
						"Код": errors.New("Неправильный"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EM(tt.args.args...)
			assert.Equal(t, tt.want, got)
		})
	}
}
