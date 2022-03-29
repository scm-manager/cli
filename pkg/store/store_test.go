package store

import (
	"github.com/scm-manager/cli/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/zalando/go-keyring"
	"os"
	"path"
	"testing"
)

func init() {
	keyring.MockInit()
}

func TestRead(t *testing.T) {
	err := keyring.Set("scm-cli", "trillian", "secret")

	config, err := readFromFilePath(path.Join("testdata", "cli-config.json"))

	assert.NoError(t, err)
	assert.Equal(t, "localhost", config.ServerUrl)
	assert.Equal(t, "trillian", config.Username)
	assert.Equal(t, "secret", config.ApiKey)
}

func TestStore(t *testing.T) {
	filePath := path.Join(t.TempDir(), "testpath")
	err := writeToFilePath(filePath, &pkg.Configuration{ServerUrl: "server", Username: "user"})
	assert.NoError(t, err)

	config, err := readFromFilePath(filePath)
	assert.Equal(t, "server", config.ServerUrl)
	assert.Equal(t, "user", config.Username)
}

func TestDeleteIfConfigNotExist(t *testing.T) {
	filePath := path.Join(t.TempDir(), "testpath")

	err := deleteFilePath(filePath)
	assert.Error(t, err)

	_, err = os.Stat(filePath)
	assert.True(t, os.IsNotExist(err))
}

func TestDelete(t *testing.T) {
	filePath := path.Join(t.TempDir(), "testpath")

	err := writeToFilePath(filePath, &pkg.Configuration{ServerUrl: "server", Username: "user"})
	assert.NoError(t, err)

	err = deleteFilePath(filePath)
	assert.NoError(t, err)

	_, err = os.Stat(filePath)
	assert.True(t, os.IsNotExist(err))
}

func TestStoreApiKey(t *testing.T) {
	keyname := "scm-cli"
	username := "scmadmin"
	apiKey := "secret_key"
	err := storeApiKey(keyname, username, apiKey)

	assert.NoError(t, err)

	storedKey, err := readApiKey(keyname, username)
	assert.NoError(t, err)

	assert.Equal(t, apiKey, storedKey)
}

func TestDeleteApiKey(t *testing.T) {
	keyname := "scm-cli"
	username := "scmadmin"
	apiKey := "secret_key"
	err := keyring.Set(keyname, username, apiKey)
	assert.NoError(t, err)

	err = deleteApiKey(keyname, username)
	assert.NoError(t, err)

	key, err := keyring.Get(keyname, username)
	assert.Empty(t, key)
}
