package pkg

import (
	"fmt"
	"net/http"
)

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/User-Agent
// scm-cli/1.0.0 (githash; isodate)

const cliName = "scm-cli"

var (
	version   string = "x.y.z"
	gitHash   string
	buildTime string
)

func createUserAgent(version, gitHash, buildTime string) string {
	if gitHash != "" && buildTime != "" {
		return fmt.Sprintf("%s/%s (%s; %s)", cliName, version, gitHash, buildTime)
	}
	return fmt.Sprintf("%s/%s", cliName, version)
}

func CreateHttpClient() *http.Client {
	client := &http.Client{}
	client.Transport = &userAgentTransport{}
	return client
}

type userAgentTransport struct {
}

func (ua *userAgentTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", createUserAgent(version, gitHash, buildTime))
	return http.DefaultTransport.RoundTrip(r)
}
