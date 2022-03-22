package auth

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/term"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func Login() {
	// Collect login parameters
	serverUrl := os.Args[2]
	username, password := ReadCredentials()

	// Send login request
	loginRequest := LoginRequest{Hostname: ResolveHostname()}
	payloadBuf := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuf).Encode(loginRequest)
	if err != nil {
		log.Fatal("Could not encode hostname: ", err)
	}
	req, _ := http.NewRequest("POST", serverUrl+"/api/v2/cli/login", payloadBuf)
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(username, password)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Could not send login request: ", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Could not read login response", err)
	}

	// Create .scm-cli.json config file
	configFilePath := ResolveConfigFilePath()
	file, err := os.Create(configFilePath)
	if err != nil {
		log.Fatal("Could not create cli config file", err)
	}
	configuration := Configuration{ServerUrl: serverUrl, ApiKey: string(body), Username: username}
	jsonConfig, err := json.Marshal(configuration)
	if err != nil {
		log.Fatal("Could not marshal config json", err)
	}
	_, err = file.Write(jsonConfig)
	if err != nil {
		log.Fatal("Could not write config to file", err)
	}

	// Store api key to keyring
	StoreApiKey(username, configuration.ApiKey)
}

func ReadCredentials() (string, string) {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, err := r.ReadString('\n')
	if err != nil {
		log.Fatal("Could not read username", err)
	}
	fmt.Print("Password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal("Could not read password", err)
	}
	return strings.TrimSpace(username), string(password)
}

type LoginRequest struct {
	Hostname string `json:"hostname"`
}
