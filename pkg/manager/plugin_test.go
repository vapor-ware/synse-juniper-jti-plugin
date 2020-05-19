package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-sdk/sdk"
)

func TestNewPluginDeviceManager(t *testing.T) {
	dm := NewPluginDeviceManager(&sdk.Plugin{})
	assert.NotNil(t, dm)
}
