package terminal

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"os"
	"strings"
)

func ReadCredentials() (string, string, error) {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, err := r.ReadString('\n')
	if err != nil {
		return "", "", fmt.Errorf("could not read username: %w", err)
	}
	fmt.Print("Password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", "", fmt.Errorf("could not read password: %w", err)
	}
	return strings.TrimSpace(username), string(password), nil
}
