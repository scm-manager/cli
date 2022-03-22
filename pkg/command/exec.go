package command

import (
	"bytes"
	"fmt"
	"github.com/scm-manager/cli/pkg/auth"
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
	req.SetBasicAuth(config.Username, config.ApiKey)
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

	fmt.Println(string(body))

}
