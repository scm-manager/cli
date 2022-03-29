package store

import (
	"encoding/json"
	"fmt"
	"github.com/scm-manager/cli/pkg"
	"github.com/scm-manager/cli/pkg/api"
	"github.com/zalando/go-keyring"
	"io/ioutil"
	"os"
	"path"
)

func readFromFilePath(filePath string) (*pkg.Configuration, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, fmt.Errorf("could not find config file: %w", err)
		}
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}
	configuration := &pkg.Configuration{}
	err = json.Unmarshal(data, configuration)
	if err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)

	}
	key, err := readApiKey(api.KeyName, configuration.Username)
	if err != nil {
		return nil, err
	}
	configuration.ApiKey = key
	return configuration, err
}

func Read() (*pkg.Configuration, error) {
	configFilePath, err := resolveConfigFilePath()
	if err != nil {
		return nil, err
	}
	return readFromFilePath(configFilePath)
}

func writeToFilePath(filePath string, configuration *pkg.Configuration) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("could not create cli config file: %w", err)
	}
	jsonConfig, err := json.Marshal(configuration)
	if err != nil {
		return fmt.Errorf("could not marshal config json: %w", err)
	}
	_, err = file.Write(jsonConfig)
	if err != nil {
		return fmt.Errorf("could not write config to file: %w", err)

	}
	err = storeApiKey(api.KeyName, configuration.Username, configuration.ApiKey)
	if err != nil {
		return fmt.Errorf("could not store api key to keyring: %w", err)
	}
	return nil
}

func Write(configuration *pkg.Configuration) error {
	configFilePath, err := resolveConfigFilePath()
	if err != nil {
		return err
	}
	return writeToFilePath(configFilePath, configuration)
}

func deleteFilePath(filePath string) error {
	config, err := readFromFilePath(filePath)
	if err != nil {
		return err
	}
	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("could not delete config file: %w", err)
	}
	err = deleteApiKey(api.KeyName, config.Username)
	if err != nil {
		return fmt.Errorf("could not delete api key from keystore: %w", err)
	}
	return nil
}

func Delete() error {
	configFilePath, err := resolveConfigFilePath()
	if err != nil {
		return err
	}
	return deleteFilePath(configFilePath)
}

func storeApiKey(apiKeyName string, username string, apiKey string) error {
	err := keyring.Set(apiKeyName, username, apiKey)
	if err != nil {
		return fmt.Errorf("could not store apikey to keyring: %w", err)
	}
	return nil
}

func readApiKey(apiKeyName string, username string) (string, error) {
	apiKey, err := keyring.Get(apiKeyName, username)
	if err != nil {
		return apiKey, fmt.Errorf("could not read apikey from keyring: %w", err)
	}
	return apiKey, nil
}

func deleteApiKey(apiKeyName string, username string) error {
	// Delete stored api key
	err := keyring.Delete(apiKeyName, username)
	if err != nil {
		if err != nil {
			return fmt.Errorf("could not delete stored api key: %w", err)
		}
	}
	return nil
}

func resolveConfigFilePath() (string, error) {
	// Read config (server url, username) => .scm-cli.json
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}
	return path.Join(homeDir, ".scm-cli.json"), nil
}
