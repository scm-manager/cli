package api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zalando/go-keyring"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	keyring.MockInit()
}

func TestCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "api-secret")
		assert.NoError(t, err)
		assert.Equal(t, "Basic YXJ0aHVyOnNlY3JldA==", r.Header.Get("Authorization"))
	}))
	defer server.Close()

	apiKey, err := Create(server.URL, "arthur", "secret", "test-key")
	assert.NoError(t, err)
	assert.Equal(t, "api-secret", apiKey)
}

func TestRemoveIfKeyNotExist(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer api-token", r.Header.Get("Authorization"))
	}))
	defer server.Close()

	err := Remove(server.URL, "api-token", "scm-test")
	assert.NoError(t, err)
}

func TestRemove(t *testing.T) {
	serviceName := "scm-test"
	apiKey := "api-token"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer "+apiKey, r.Header.Get("Authorization"))
	}))
	defer server.Close()

	err := Remove(server.URL, apiKey, serviceName)
	assert.NoError(t, err)
}
