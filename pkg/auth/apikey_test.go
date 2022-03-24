package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStoreApiKey(t *testing.T) {

	username := "scmadmin"
	apiKey := "secret_key"
	err := StoreApiKey(username, apiKey)

	assert.NoError(t, err)

	storedKey, err := ReadApiKey(username)
	assert.NoError(t, err)

	assert.Equal(t, apiKey, storedKey)
}
