package command

import (
	"encoding/json"
	"fmt"
	"github.com/Xuanwo/go-locale"
	"github.com/scm-manager/cli/pkg"
	"io"
	"net/http"
	"os"
)

type ExecuteResponse struct {
	Out  string
	Err  string
	Exit int
}

func CreateExecutor(stdout io.Writer, stderr io.Writer, stdin io.Reader, config *pkg.Configuration) *Executor {
	return &Executor{stdout: stdout, stderr: stderr, stdin: stdin, config: config}
}

func CreateDefaultExecutor(config *pkg.Configuration) (*Executor, error) {
	stdin, err := createStdin()
	if err != nil {
		return nil, err
	}
	return CreateExecutor(os.Stdout, os.Stderr, stdin, config), nil
}

func createStdin() (io.Reader, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to read stat from stdin: %w", err)
	}
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return os.Stdin, nil
	}
	return nil, nil
}

type Executor struct {
	stdout io.Writer
	stderr io.Writer
	stdin  io.Reader
	config *pkg.Configuration
}

func (e *Executor) Execute(args ...string) (int, error) {
	req, err := e.createRequest(args)
	if err != nil {
		return -1, err
	}
	res, err := e.sendRequest(req)
	if err != nil {
		return -1, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	return e.processResponse(res)
}

func (e *Executor) sendRequest(req *http.Request) (*http.Response, error) {
	client := pkg.CreateHttpClient()
	res, err := client.Do(req)
	if err != nil {
		return res, fmt.Errorf("could not send execution request: %w", err)
	}
	if res.StatusCode >= 400 {
		return res, fmt.Errorf("HTTP Error %d", res.StatusCode)
	}
	return res, nil
}

func (e *Executor) createRequest(args []string) (*http.Request, error) {
	req, _ := http.NewRequest("POST", e.config.ServerUrl+"/api/v2/cli/exec", e.stdin)
	req.Header.Add("Authorization", "Bearer "+e.config.ApiKey)
	e.setCommandVarArgs(req, args)
	err := e.setRequestLocale(req)
	if err != nil {
		return req, err
	}
	return req, nil
}

func (e *Executor) processResponse(res *http.Response) (int, error) {
	decoder := json.NewDecoder(res.Body)
	_, err := decoder.Token()
	if err != nil {
		return -1, fmt.Errorf("could not read token from json: %w", err)
	}

	for decoder.More() {
		response := &ExecuteResponse{}
		err := decoder.Decode(response)
		if err != nil {
			if err == io.EOF {
				break
			}
			return -1, fmt.Errorf("could not decode response body: %w", err)
		}
		if response.Out != "" {
			_, err := e.stdout.Write([]byte(response.Out))
			if err != nil {
				return -1, fmt.Errorf("could not write to stdout: %w", err)
			}
		}
		if response.Err != "" {
			_, err := e.stderr.Write([]byte(response.Err))
			if err != nil {
				return -1, fmt.Errorf("could not write to stderr: %w", err)
			}
		}
		if response.Exit != 0 {
			return response.Exit, nil
		}
	}
	return 0, nil
}

func (e *Executor) setCommandVarArgs(req *http.Request, args []string) {
	queryString := req.URL.Query()
	for _, arg := range args {
		queryString.Add("args", arg)
	}
	req.URL.RawQuery = queryString.Encode()
}

func (e *Executor) setRequestLocale(req *http.Request) error {
	language, err := locale.Detect()
	if err != nil {
		return fmt.Errorf("could not detect client locale: %w", err)
	}
	baseLang, _ := language.Base()
	req.Header.Add("Accept-Language", baseLang.String())
	return nil
}
