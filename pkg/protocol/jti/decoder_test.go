package jti

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/manager"
)

func TestNewJTIDecoder(t *testing.T) {
	decoder := NewJTIDecoder(manager.NewStubDeviceManager(false))
	assert.NotNil(t, decoder)
}

func TestJuniperJTIDecoder_Decode_NilManager(t *testing.T) {
	decoder := JuniperJTIDecoder{}

	data, err := decoder.Decode([]byte{})
	assert.Error(t, err)
	assert.Nil(t, data)
}

func TestJuniperJTIDecoder_Decode_FailedDecode(t *testing.T) {
	decoder := NewJTIDecoder(manager.NewStubDeviceManager(false))

	data, err := decoder.Decode([]byte{})
	assert.Error(t, err)
	assert.Nil(t, data)
}
