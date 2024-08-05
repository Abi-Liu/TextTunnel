package server_test

import (
	"os"
	"testing"

	"github.com/Abi-Liu/TextTunnel/internal/server"
	"github.com/stretchr/testify/assert"
)

const PORT = "8080"
const DB_STRING = "randomdbstring"

func TestLoadEnv(t *testing.T) {
	t.Run("missing port env", func(t *testing.T) {
		// mocking an unset port env
		os.Setenv("PORT", "")

		_, err := server.LoadEnv()
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "port env missing")
	})

	t.Run("missing db string", func(t *testing.T) {
		os.Setenv("PORT", PORT)
		_, err := server.LoadEnv()
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "db string missing")
	})

	t.Run("valid env vars", func(t *testing.T) {
		os.Setenv("PORT", PORT)
		os.Setenv("DATABASE_URL", DB_STRING)
		env, err := server.LoadEnv()

		assert.NoError(t, err)
		assert.Equal(t, PORT, env.PORT)
		assert.Equal(t, DB_STRING, env.DB_URL)
	})

}
