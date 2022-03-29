package command

import (
	"bytes"
	"fmt"
	"github.com/scm-manager/cli/pkg"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExecutor_Execute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "[{\"out\":\"Hello World\"}]")
		assert.NoError(t, err)
	}))
	defer server.Close()
	var stdout bytes.Buffer
	configuration := &pkg.Configuration{ServerUrl: server.URL, Username: "scmadmin", ApiKey: "secret"}

	executor := CreateExecutor(&stdout, nil, nil, configuration)
	exitCode, err := executor.Execute("some", "command")

	assert.NoError(t, err)
	assert.Equal(t, 0, exitCode)
	assert.Equal(t, "Hello World", stdout.String())
}

func TestExecutor_ExecuteCheckForArgs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, []string{"some", "command"}, r.Form["args"])
		_, err = fmt.Fprintf(w, "[{\"out\":\"Hello World\"}]")
		assert.NoError(t, err)
	}))
	defer server.Close()
	var stdout bytes.Buffer
	configuration := &pkg.Configuration{ServerUrl: server.URL, Username: "scmadmin", ApiKey: "secret"}

	executor := CreateExecutor(&stdout, nil, nil, configuration)
	_, err := executor.Execute("some", "command")

	assert.NoError(t, err)
}

func TestExecutor_ExecuteCheckWithApiKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer secret", r.Header.Get("Authorization"))
		_, err := fmt.Fprintf(w, "[{\"out\":\"Hello World\"}]")
		assert.NoError(t, err)
	}))
	defer server.Close()
	var stdout bytes.Buffer
	configuration := &pkg.Configuration{ServerUrl: server.URL, Username: "scmadmin", ApiKey: "secret"}

	executor := CreateExecutor(&stdout, nil, nil, configuration)
	_, err := executor.Execute("some", "command")

	assert.NoError(t, err)
}

func TestExecutor_ExecuteCheckStderr(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "[{\"err\":\"Missing entity\"}]")
		assert.NoError(t, err)
	}))
	defer server.Close()
	var stderr bytes.Buffer
	configuration := &pkg.Configuration{ServerUrl: server.URL, Username: "scmadmin", ApiKey: "secret"}

	executor := CreateExecutor(nil, &stderr, nil, configuration)
	_, err := executor.Execute("some", "command")

	assert.NoError(t, err)
	assert.Equal(t, "Missing entity", stderr.String())
}

func TestExecutor_ExecuteCheckExitCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "[{\"exit\":42}]")
		assert.NoError(t, err)
	}))
	defer server.Close()
	configuration := &pkg.Configuration{ServerUrl: server.URL, Username: "scmadmin", ApiKey: "secret"}

	executor := CreateExecutor(nil, nil, nil, configuration)
	exitCode, err := executor.Execute("some", "command")

	assert.NoError(t, err)
	assert.Equal(t, 42, exitCode)
}

func TestExecutor_ExecuteCheckLocale(t *testing.T) {
	t.Setenv("LANGUAGE", "es_MX")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "es", r.Header.Get("Accept-Language"))
		_, err := fmt.Fprintf(w, "[{\"exit\":0}]")
		assert.NoError(t, err)
	}))
	defer server.Close()
	configuration := &pkg.Configuration{ServerUrl: server.URL, Username: "scmadmin", ApiKey: "secret"}

	executor := CreateExecutor(nil, nil, nil, configuration)
	_, err := executor.Execute("some", "command")

	assert.NoError(t, err)
}

func TestExecutor_ExecuteCheckHttpError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not found", http.StatusNotFound)
	}))
	defer server.Close()
	configuration := &pkg.Configuration{ServerUrl: server.URL, Username: "scmadmin", ApiKey: "secret"}

	executor := CreateExecutor(nil, nil, nil, configuration)
	_, err := executor.Execute("some", "command")

	assert.ErrorContains(t, err, "HTTP Error 404")
}

func TestCreateDefaultExecutor(t *testing.T) {
	configuration := &pkg.Configuration{ServerUrl: "myServer", Username: "scmadmin", ApiKey: "secret"}

	executor, err := CreateDefaultExecutor(configuration)

	assert.NoError(t, err)
	assert.NotNil(t, executor)
}
