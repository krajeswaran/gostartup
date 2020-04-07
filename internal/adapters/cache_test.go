package adapters

import (
	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)


func TestCacheAdapter_CacheInit(t *testing.T) {
	tests := []struct {
		name string
		want *redis.Cmdable
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CacheAdapter{}
			if got := c.CacheInit(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CacheInit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Setup() {
	DB = mock.Mock{}
}

func TestCacheAdapter_DeepStatus(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CacheAdapter{}
			if err := c.DeepStatus(); (err != nil) != tt.wantErr {
				t.Errorf("DeepStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCacheAdapter_GetApiStats(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CacheAdapter{}
			got, err := c.GetApiStats()
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

func TestCacheAdapter_ResetApiStats(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CacheAdapter{}
			if err := c.ResetApiStats(); (err != nil) != tt.wantErr {
				t.Errorf("ResetApiStats() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCacheAdapter_UpdateApiStats(t *testing.T) {
	type args struct {
		didApiFail bool
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CacheAdapter{}
			got, err := c.UpdateApiStats(tt.args.didApiFail)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateApiStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UpdateApiStats() got = %v, want %v", got, tt.want)
			}
		})
	}
}