package pkg

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateUserAgent(t *testing.T) {
	userAgent := createUserAgent("1.0.0", "bb4d50d", "2022-04-04T09:41:03Z")
	assert.Equal(t, "scm-cli/1.0.0 (bb4d50d; 2022-04-04T09:41:03Z)", userAgent)
}

func TestCreateUserAgent_WithoutGitHashAndBuildTime(t *testing.T) {
	fmt.Println(time.Now().Format(time.RFC3339))
	userAgent := createUserAgent("x.y.z", "", "")
	assert.Equal(t, "scm-cli/x.y.z", userAgent)
}

func TestCreateHttpClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "scm-cli/x.y.z", r.Header.Get("User-Agent"))
	}))
	defer server.Close()

	client := CreateHttpClient()
	_, err := client.Get(server.URL)
	assert.NoError(t, err)
}
