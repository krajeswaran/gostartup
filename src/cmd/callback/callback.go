package main

import (
	"gostartup/src/adapters"
	"gostartup/src/config"
	"gostartup/src/server/callback"
)

func main() {
	config.Init()
	adapters.Init()
	callback.CallbackInit()
}
