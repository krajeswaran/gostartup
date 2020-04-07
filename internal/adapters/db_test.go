package adapters

import (
	"github.com/go-pg/pg/v9"
	"github.com/krajeswaran/gostartup/internal/models"
	"reflect"
	"testing"
)

func TestDBAdapter_CreateUser(t *testing.T) {
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
			d := &DBAdapter{}
			got, err := d.CreateUser(tt.args.name)
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

func TestDBAdapter_DBInit(t *testing.T) {
	tests := []struct {
		name string
		want *pg.DB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DBAdapter{}
			if got := d.DBInit(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DBInit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBAdapter_DeepStatus(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DBAdapter{}
			if err := d.DeepStatus(); (err != nil) != tt.wantErr {
				t.Errorf("DeepStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBAdapter_FetchUser(t *testing.T) {
	type args struct {
		id string
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
			d := &DBAdapter{}
			got, err := d.FetchUser(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}