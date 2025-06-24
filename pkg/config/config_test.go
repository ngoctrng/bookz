package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	envVars := map[string]string{
		"APP_ENV":       "test",
		"PORT":          "8080",
		"SENTRY_DSN":    "https://test.sentry.io",
		"ALLOW_ORIGINS": "localhost",
		"DB_NAME":       "testdb",
		"DB_HOST":       "localhost",
		"DB_PORT":       "5432",
		"DB_USER":       "testuser",
		"DB_PASS":       "testpass",
		"ENABLE_SSL":    "true",
	}

	for k, v := range envVars {
		os.Setenv(k, v)
	}

	defer func() {
		for k := range envVars {
			os.Unsetenv(k)
		}
	}()

	t.Run("LoadConfig should load all environment variables correctly", func(t *testing.T) {
		cfg, err := LoadConfig()

		assert.NoError(t, err)
		assert.NotNil(t, cfg)

		assert.Equal(t, "test", cfg.AppEnv)
		assert.Equal(t, 8080, cfg.Port)
		assert.Equal(t, "https://test.sentry.io", cfg.SentryDSN)
		assert.Equal(t, "localhost", cfg.AllowOrigins)

		assert.Equal(t, "testdb", cfg.DB.Name)
		assert.Equal(t, "localhost", cfg.DB.Host)
		assert.Equal(t, 5432, cfg.DB.Port)
		assert.Equal(t, "testuser", cfg.DB.User)
		assert.Equal(t, "testpass", cfg.DB.Pass)
		assert.Equal(t, true, cfg.DB.EnableSSL)
	})

	t.Run("Empty config should be initialized", func(t *testing.T) {
		assert.NotNil(t, Empty)
	})
}

func TestLoadConfigWithInvalidValues(t *testing.T) {
	// Test with invalid PORT value
	os.Setenv("PORT", "invalid")
	defer os.Unsetenv("PORT")

	cfg, err := LoadConfig()
	assert.Error(t, err)
	assert.Nil(t, cfg)
}
