package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/config"
	"github.com/vapor-ware/synse-sdk/sdk"
)

func TestRunBackgroundListenerAction_ErrNoConfig(t *testing.T) {
	defer config.Set(nil)

	err := RunBackgroundListener.Action(&sdk.Plugin{})
	assert.Error(t, err)
}
