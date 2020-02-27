package main

import (
	"github.com/krajeswaran/gostartup/config"
	"github.com/krajeswaran/gostartup/internal/adapters"
	"github.com/krajeswaran/gostartup/internal/servers/api"
	"sync"
)

var once sync.Once

func main() {
	once.Do(func() {
		config.Init()
		adapters.Init()
		api.Init()
	})
}
