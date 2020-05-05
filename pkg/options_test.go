package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/config"
)

func TestLoadDynamicConfig(t *testing.T) {
	defer func() {
		config.Set(nil)
		configCount = 0
	}()
	assert.Nil(t, config.Get())

	devices, err := LoadDynamicConfig(map[string]interface{}{
		"address": "localhost",
	})
	assert.NoError(t, err)
	assert.Empty(t, devices) // For this plugin, no devices are created at this time

	assert.NotNil(t, config.Get())
	assert.Equal(t, "localhost", config.Get().Address)
}

func TestLoadDynamicConfig_ErrorTooManyServers(t *testing.T) {
	defer func() {
		config.Set(nil)
		configCount = 0
	}()

	configCount = 1

	devices, err := LoadDynamicConfig(map[string]interface{}{
		"address": "localhost",
	})
	assert.Error(t, err)
	assert.Equal(t, ErrTooManyServers, err)
	assert.Nil(t, devices)

	assert.Nil(t, config.Get())
}

func TestLoadDynamicConfig_ErrorLoadConfig(t *testing.T) {
	defer func() {
		config.Set(nil)
		configCount = 0
	}()
	assert.Nil(t, config.Get())

	devices, err := LoadDynamicConfig(map[string]interface{}{})
	assert.Error(t, err)
	assert.Nil(t, devices)

	assert.Nil(t, config.Get())
}

func TestDeviceIdentifier(t *testing.T) {
	id := DeviceIdentifier(map[string]interface{}{
		"id": map[string]string{
			"test":  "",
			"foo":   "bar",
			"index": "1",
		},
	})

	assert.Equal(t, "bar1", id)
}

func TestDeviceIdentifier_NoID(t *testing.T) {
	assert.Panics(t, func() {
		DeviceIdentifier(map[string]interface{}{
			"foo":   "bar",
			"index": "1",
		})
	})
}

func TestDeviceIdentifier_BadType(t *testing.T) {
	assert.Panics(t, func() {
		DeviceIdentifier(map[string]interface{}{
			"id": map[string]int{
				"test":  1,
				"foo":   2,
				"index": 3,
			},
		})
	})
}
