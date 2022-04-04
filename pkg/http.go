package pkg

import (
	"fmt"
	"net/http"
	"runtime"
)

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/User-Agent
// name/version (os arch; commithash; buildtime)

const cliName = "scm-cli"

var (
	version    string = "x.y.z"
	commitHash string
	buildTime  string
)

func createDevUserAgent(version, osName, arch string) string {
	return fmt.Sprintf("%s/%s (%s %s)", cliName, version, osName, arch)
}

func createProdUserAgent(version, osName, arch, gitHash, buildTime string) string {
	return fmt.Sprintf("%s/%s (%s %s; %s; %s)", cliName, version, osName, arch, gitHash, buildTime)
}

func createUserAgent(version, osName, arch, gitHash, buildTime string) string {
	if gitHash != "" && buildTime != "" {
		return createProdUserAgent(version, osName, arch, gitHash, buildTime)
	}
	return createDevUserAgent(version, osName, arch)
}

func CreateHttpClient() *http.Client {
	client := &http.Client{}
	client.Transport = &userAgentTransport{}
	return client
}

type userAgentTransport struct {
}

func (t *userAgentTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	ua := createUserAgent(version, runtime.GOOS, runtime.GOARCH, commitHash, buildTime)
	r.Header.Set("User-Agent", ua)
	return http.DefaultTransport.RoundTrip(r)
}
