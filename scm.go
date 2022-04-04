package main

import (
	"fmt"
	"github.com/scm-manager/cli/pkg"
	"github.com/scm-manager/cli/pkg/api"
	"github.com/scm-manager/cli/pkg/command"
	"github.com/scm-manager/cli/pkg/store"
	"github.com/scm-manager/cli/pkg/terminal"
	"log"
	"os"
	"strings"
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

func createApiKeyName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "scm-cli"
	}
	return "scm-cli-" + hostname
}

func readConfig() *pkg.Configuration {
	configuration, err := store.Read()
	if err != nil {
		log.Fatalf("Could not read configuration: %v", err)
	}
	return configuration
}

func login() {
	// Collect login parameters
	serverUrl := os.Args[2]
	if strings.HasSuffix(serverUrl, "/") {
		serverUrl = serverUrl[0 : len(serverUrl)-1]
	}

	username, password, err := terminal.ReadCredentials()
	// line break after Password:
	fmt.Println()
	if err != nil {
		printLoginError("Could not read credentials: %v", err)
	}

	// Create api key
	apiKey, err := api.Create(serverUrl, username, password, createApiKeyName())
	if err != nil {
		printLoginError("Could not create api key: %v", err)
	}
	// Write Config
	err = store.Write(&pkg.Configuration{ServerUrl: serverUrl, Username: username, ApiKey: apiKey})
	if err != nil {
		printLoginError("Could not write config to store: %v", err)
	}
}

func printLoginError(format string, err error) {
	fmt.Println()
	log.Fatalf(format, err)
}

func logout(configuration *pkg.Configuration) {
	err := api.Remove(configuration.ServerUrl, configuration.ApiKey, createApiKeyName())
	if err != nil {
		fmt.Printf("Failed to remove api key from server: %v", err)
		fmt.Println("We suggest you remove the api key manually on your SCM-Manager server.")
	}
	err = store.Remove()
	if err != nil {
		log.Fatalf("Could not remove local config: %v", err)
	}
	fmt.Println("Successfully logged out")
}

func executeCommand(configuration *pkg.Configuration) {
	executor, err := command.CreateDefaultExecutor(configuration)
	if err != nil {
		log.Fatalf("Failed to create default executor: %v", err)
	}
	exitCode, err := executor.Execute(os.Args[1:]...)
	if err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}
	os.Exit(exitCode)
}
