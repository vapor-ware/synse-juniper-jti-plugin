package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/protocol"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
)

func Test_jtiDeviceRead(t *testing.T) {
	d := &sdk.Device{
		Data: map[string]interface{}{
			protocol.ReadingKey: []*output.Reading{
				{
					Type:  "test",
					Value: 100,
				},
			},
		},
	}

	readings, err := jtiDeviceRead(d)
	assert.NoError(t, err)
	assert.Len(t, readings, 1)
	assert.Equal(t, "test", readings[0].Type)
	assert.Equal(t, 100, readings[0].Value)
}

func Test_jtiDeviceRead_ErrNoReadings(t *testing.T) {
	d := &sdk.Device{
		Data: map[string]interface{}{},
	}

	readings, err := jtiDeviceRead(d)
	assert.Error(t, err)
	assert.Nil(t, readings)
}

func Test_jtiDeviceRead_ErrBadType(t *testing.T) {
	d := &sdk.Device{
		Data: map[string]interface{}{
			protocol.ReadingKey: []string{
				"foo", "bar",
			},
		},
	}

	readings, err := jtiDeviceRead(d)
	assert.Error(t, err)
	assert.Nil(t, readings)
}
