package usecases

import (
	"gostartup/src/adapters"
	"gostartup/src/config"
	"gostartup/src/models"
	"reflect"
	"strconv"
	"testing"
)

var accountUsecase *SmsUsecase

func setup() {
	config.Init()
	adapters.Init()
}

func tearDown() {
	c := adapters.Cache()
	c.FlushDB()
}

func TestHelloRepo_CreateUser(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HelloRepo{}
			got, err := h.CreateUser(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHelloRepo_FetchUserName(t *testing.T) {
	type args struct {
		userId string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HelloRepo{}
			got, err := h.FetchUserName(tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchUserName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FetchUserName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHelloRepo_GetApiStats(t *testing.T) {
	tests := []struct {
		name    string
		want    *models.Stat
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HelloRepo{}
			got, err := h.GetApiStats()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetApiStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetApiStats() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHelloRepo_updateApiStats(t *testing.T) {
	type args struct {
		didApiFail bool
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
			h := &HelloRepo{}
			if err := h.updateApiStats(tt.args.didApiFail); (err != nil) != tt.wantErr {
				t.Errorf("updateApiStats() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}