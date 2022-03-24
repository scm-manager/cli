package command

import (
	"bytes"
	"encoding/json"
	"github.com/Xuanwo/go-locale"
	"github.com/scm-manager/cli/pkg/auth"
	"io"
	"io/ioutil"
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
	language, err := locale.Detect()
	if err != nil {
		log.Fatal("Could not detect client locale")
	}
	baseLang, _ := language.Base()
	req.Header.Add("Accept-Language", baseLang.String())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Could not send login request: ", err)
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Could not read response", err)
	}
	//TODO Remove after implementation
	//fmt.Println(string(data))
	decoder := json.NewDecoder(bytes.NewBuffer(data))
	_, err = decoder.Token()
	if err != nil {
		log.Fatal("Could not read token from json", err)
	}

	for decoder.More() {
		response := &ExecuteResponse{}
		err := decoder.Decode(response)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal("Error reading response body", err)
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
		if response.Exit != 0 {
			os.Exit(response.Exit)
		}
	}
}

type ExecuteResponse struct {
	Out  string
	Err  string
	Exit int
}
