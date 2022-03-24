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
	req := createRequest()
	res := sendRequest(req)
	defer res.Body.Close()
	processResponse(res)
}

func sendRequest(req *http.Request) *http.Response {
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Could not send login request: ", err)
	}
	return res
}

func createRequest() *http.Request {
	config := auth.ReadConfig()

	payloadBuf := new(bytes.Buffer)
	req, _ := http.NewRequest("POST", config.ServerUrl+"/api/v2/cli/exec", payloadBuf)
	req.Header.Add("Authorization", "Bearer "+config.ApiKey)
	setCommandVarArgs(req)
	setRequestLocale(req)
	return req
}

func processResponse(res *http.Response) {
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

func setCommandVarArgs(req *http.Request) {
	queryString := req.URL.Query()
	for i, arg := range os.Args {
		if i > 0 {
			queryString.Add("args", arg)
		}
	}
	req.URL.RawQuery = queryString.Encode()
}

func setRequestLocale(req *http.Request) {
	language, err := locale.Detect()
	if err != nil {
		log.Fatal("Could not detect client locale")
	}
	baseLang, _ := language.Base()
	req.Header.Add("Accept-Language", baseLang.String())
}

type ExecuteResponse struct {
	Out  string
	Err  string
	Exit int
}
