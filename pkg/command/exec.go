package command

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/scm-manager/cli/pkg/auth"
	"io"
	"log"
	"net/http"
	"os"
)

func ExecuteCommand() {
	config := auth.ReadConfig()

	payloadBuf := new(bytes.Buffer)
	req, _ := http.NewRequest("POST", config.ServerUrl+"/api/v2/cli/exec", payloadBuf)
	queryString := req.URL.Query()
	for i, arg := range os.Args {
		if i > 0 {
			queryString.Add("args", arg)
		}
	}
	req.URL.RawQuery = queryString.Encode()
	req.Header.Add("Authorization", "Bearer "+config.ApiKey)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Could not send login request: ", err)
	}
	defer res.Body.Close()

	reader := bufio.NewReader(res.Body)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal("Error reading response body", err)
		}
		response := &ExecuteResponse{}
		err = json.Unmarshal(line, response)
		if err != nil {
			log.Fatal("Could not decode json response", err)
		}
		if response.Out != "" {
			_, err := os.Stdout.WriteString(response.Out)
			if err != nil {
				log.Fatal("Could not write to stdout", err)
			}
		}
		if response.Err != "" {
			_, err := os.Stderr.WriteString(response.Err)
			if err != nil {
				log.Fatal("Could not write to stderr", err)
			}
		}
	}
}

type ExecuteResponse struct {
	Out string
	Err string
}
