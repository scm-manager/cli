package main

import (
	"fmt"
	"github.com/scm-manager/cli/pkg/auth"
	"github.com/scm-manager/cli/pkg/command"
	"os"
)

func main() {
	config := auth.ReadConfig()
	if config == nil {
		if len(os.Args) > 2 && os.Args[1] == "login" {
			auth.Login()
		} else {
			fmt.Println("Please login first calling \"scm login {server-url}\"")
			os.Exit(1)
		}
	}

	if os.Args[1] == "logout" {
		auth.Logout()
		fmt.Println("Successfully logged out")
	} else {
		command.ExecuteCommand()
	}
}
