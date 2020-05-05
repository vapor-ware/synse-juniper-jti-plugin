package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_Ok(t *testing.T) {
	raw := map[string]interface{}{
		"address": "localhost",
	}

	cfg, err := Load(raw)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "localhost", cfg.Address)
}

func TestLoad_Error(t *testing.T) {
	raw := make(map[string]interface{})

	cfg, err := Load(raw)
	assert.Error(t, err)
	assert.Equal(t, ErrNoAddress, err)
	assert.Nil(t, cfg)
}

func TestSet(t *testing.T) {
	defer func() {
		serverConfig = nil
	}()

	assert.Nil(t, serverConfig)

	cfg := &ServerConfig{
		Address: "localhost",
	}

	Set(cfg)
	assert.Equal(t, cfg, serverConfig)
}

func TestGet(t *testing.T) {
	cfg := Get()
	assert.Nil(t, cfg)
}

func TestGet2(t *testing.T) {
	cfg := &ServerConfig{Address: "localhost"}
	serverConfig = cfg

	c := Get()
	assert.Equal(t, cfg, c)
}
