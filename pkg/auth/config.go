package auth

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func ReadConfig() *Configuration {
	configFilePath := ResolveConfigFilePath()

	_, err := os.Stat(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			log.Fatal("Could not read config file", err)
		}
	}

	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal("Error reading file", err)
	}

	config := &Configuration{}

	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatal("Error parsing config file", err)
	}
	config.ApiKey = ReadApiKey(config.Username)
	return config
}

func ResolveConfigFilePath() string {
	// Read config (server url, username) => .scm-cli.json
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not find home directory", err)
	}

	return path.Join(homeDir, ".scm-cli.json")
}

func ResolveHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Could not resolve hostname: ", err)
	}
	return hostname
}

type Configuration struct {
	ServerUrl string
	Username  string
	ApiKey    string `json:"-"`
}
