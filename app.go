package main

import (
	"fmt"
	"github.com/scm-manager/cli/pkg/api"
	"github.com/scm-manager/cli/pkg/auth"
	"github.com/scm-manager/cli/pkg/command"
	"github.com/scm-manager/cli/pkg/config"
	"log"
	"os"
)

func main() {
	configuration := readConfig()
	if configuration == nil {
		if len(os.Args) > 2 && os.Args[1] == "login" {
			login()
			return
		} else {
			fmt.Println("Please login first calling \"scm login {server-url}\"")
			os.Exit(1)
		}
	}

	if len(os.Args) > 1 && os.Args[1] == "logout" {
		logout(configuration)
	} else {
		executeCommand(configuration)
	}
}

func readConfig() *config.Configuration {
	configuration, err := config.Read()
	if err != nil {
		log.Fatalf("Could not read configuration: %v", err)
	}
	return configuration
}

func login() {
	// Collect login parameters
	serverUrl := os.Args[2]
	username, password, err := auth.ReadCredentials()
	if err != nil {
		log.Fatalf("Could not read credentials: %v", err)
	}
	// Create api key
	apiKey, err := api.Create(serverUrl, username, password)
	if err != nil {
		log.Fatal(err)
	}
	// Write Config
	err = config.Store(&config.Configuration{ServerUrl: serverUrl, Username: username, ApiKey: apiKey})
	if err != nil {
		log.Fatal(err)
	}
}

func logout(configuration *config.Configuration) {
	err := api.Remove(api.KeyName, configuration.ServerUrl, configuration.Username, configuration.ApiKey)
	if err != nil {
		log.Fatal(err)
	}
	err = config.Delete()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully logged out")
}

func executeCommand(configuration *config.Configuration) {
	executor, err := command.CreateDefaultExecutor(configuration)
	if err != nil {
		log.Fatal("Failed to create default executor", err)
	}
	exitCode, err := executor.Execute(os.Args[1:]...)
	if err != nil {
		log.Fatal("Failed to execute command", err)
	}
	os.Exit(exitCode)
}
