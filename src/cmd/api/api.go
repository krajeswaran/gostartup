package main

import (
	"gostartup/src/adapters"
	"gostartup/src/config"
	"gostartup/src/server/api"
)

func main() {
	config.Init()
	adapters.Init()
	api.Init()
}
