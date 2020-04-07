package controllers

import (
	"github.com/labstack/echo/v4"
	"reflect"
	"testing"
)

func TestNewApiResponse(t *testing.T) {
	type args struct {
		status int
		msg    string
		data   interface{}
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 *echo.Map
	}{
		{
			name:  "test null response",
			args:  args{
				status: 0,
				msg:    "",
				data:   nil,
			},
			want:  0,
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := NewApiResponse(tt.args.status, tt.args.msg, tt.args.data)
			if got != tt.want {
				t.Errorf("NewApiResponse() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("NewApiResponse() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}