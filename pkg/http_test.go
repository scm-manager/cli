package pkg

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUserAgent(t *testing.T) {
	userAgent := createUserAgent("1.0.0", "linux", "arm64", "bb4d50d", "2022-04-04T09:41:03Z")
	assert.Equal(t, "scm-cli/1.0.0 (linux arm64; bb4d50d; 2022-04-04T09:41:03Z)", userAgent)
}

func TestCreateUserAgent_WithoutGitHashAndBuildTime(t *testing.T) {
	userAgent := createUserAgent("x.y.z", "darwin", "arm64", "", "")
	assert.Equal(t, "scm-cli/x.y.z (darwin arm64)", userAgent)
}

func TestCreateHttpClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.Header.Get("User-Agent"), "scm-cli/x.y.z")
	}))
	defer server.Close()

	client := CreateHttpClient()
	_, err := client.Get(server.URL)
	assert.NoError(t, err)
}
