package main

import (
	"fmt"
	"github.com/krajeswaran/gostartup/config"
	"github.com/krajeswaran/gostartup/internal/adapters"
	"github.com/krajeswaran/gostartup/internal/controllers"
	"github.com/spf13/cobra"
	"os"
	"sync"
)

var once sync.Once
var helloController = controllers.HelloController{}

func main() {
	once.Do(func() {
		config.Init()
		adapters.Init()
	})

	var rootCmd = &cobra.Command{
		Use:   "manage",
		Short: "Manage CLI tool for your service",
		Long: `Home for management commands, scripts,  batch jobs etc
			for this service/backend.`,
	}

	var addUserCmd = &cobra.Command{
		Use:   "adduser",
		Short: "Adds a user and creates JWT token for the user",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			t, err := helloController.CreateUser(args[0])
			if err != nil {
				fmt.Println("Failed to add user, error: ", err)
				os.Exit(-1)
			}

			fmt.Println("User created, JWT token for API: ", t)
		},
	}
	rootCmd.AddCommand(addUserCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
