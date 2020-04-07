package controllers

import (
	"github.com/labstack/echo/v4"
	"testing"
)

func TestCommonController_DeepStatus(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &CommonController{}
			if err := a.DeepStatus(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("DeepStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommonController_Status(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &CommonController{}
			if err := a.Status(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Status() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}