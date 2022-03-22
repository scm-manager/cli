package auth

import (
	"bytes"
	"encoding/json"
	"github.com/zalando/go-keyring"
	"io"
	"log"
	"net/http"
	"os"
)

func Logout() {
	config := ReadConfig()

	// Send logout request to server
	payloadBuf := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuf).Encode("")
	req, _ := http.NewRequest("DELETE", config.ServerUrl+"/api/v2/cli/logout/"+ResolveHostname(), payloadBuf)
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(config.Username, ReadApiKey(config.Username))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Could not revoke api key on server: ", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("Could not resolve response body", err)
		}
	}(res.Body)

	// Delete local config
	filePath := ResolveConfigFilePath()
	err = os.Remove(filePath)
	if err != nil {
		log.Fatal("Could not delete config file", err)
	}

	// Delete stored api key
	err = keyring.Delete("scm-cli", config.Username)
	if err != nil {
		log.Fatal("Could not delete stored api key", err)
	}
}
