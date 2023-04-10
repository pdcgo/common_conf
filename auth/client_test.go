package auth_test

import (
	_ "embed"
	"testing"

	"github.com/pdcgo/common_conf/auth"
	"github.com/stretchr/testify/assert"
)

//go:embed ..\endpoint.url
var endpoint []byte

func getEndpoint() string {
	return string(endpoint)
}

func TestClient(t *testing.T) {
	client := auth.NewAuthClient(getEndpoint())

	err := client.Login("testtiktok@gmail.com", "password", 8, "testing")
	assert.Nil(t, err)

	t.Run("testing gagal", func(t *testing.T) {
		err := client.Login("testtiktok@gmail.com", "passwords", 8, "testing")
		assert.NotNil(t, err)
	})
}
