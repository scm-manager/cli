package main

import (
	"fmt"
	"github.com/scm-manager/cli/pkg/auth"
	"github.com/scm-manager/cli/pkg/command"
	"log"
	"os"
)

func main() {
	config, err := auth.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	if config == nil {
		if len(os.Args) > 2 && os.Args[1] == "login" {
			err := auth.Login()
			if err != nil {
				log.Fatal(err)
			}
			return
		} else {
			fmt.Println("Please login first calling \"scm login {server-url}\"")
			os.Exit(1)
		}
	}

	if len(os.Args) > 1 && os.Args[1] == "logout" {
		err := auth.Logout()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Successfully logged out")
	} else {
		executor, err := command.CreateDefaultExecutor(config)
		if err != nil {
			log.Fatal("Failed to create default executor", err)
		}
		exitCode, err := executor.Execute(os.Args[1:]...)
		if err != nil {
			log.Fatal("Failed to execute command", err)
		}
		os.Exit(exitCode)
	}
}
