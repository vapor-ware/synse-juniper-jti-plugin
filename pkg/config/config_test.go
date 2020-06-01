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
	assert.Empty(t, cfg.Context)
}

func TestLoad_Ok2(t *testing.T) {
	raw := map[string]interface{}{
		"address": "udp://1.2.3.4:30000",
		"context": map[string]string{
			"foo": "bar",
		},
	}

	cfg, err := Load(raw)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "udp://1.2.3.4:30000", cfg.Address)
	assert.Equal(t, map[string]string{"foo": "bar"}, cfg.Context)
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
