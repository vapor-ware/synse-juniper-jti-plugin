package protocol

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/config"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/manager"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/protocol/jti"
)

func TestNewJtiUDPServer(t *testing.T) {
	svr := NewJtiUDPServer(
		&config.ServerConfig{
			Address: "localhost",
			Context: map[string]string{
				"site": "test",
			},
		},
		manager.NewStubDeviceManager(false),
	)

	assert.Equal(t, svr.Address, "localhost")
	assert.Equal(t, map[string]string{"site": "test"}, svr.GlobalContext)
	assert.Equal(t, svr.BufferSize, uint64(64*1024))
	assert.False(t, svr.stopped)
	assert.Nil(t, svr.conn)
	assert.NotNil(t, svr.decoder)
	assert.NotNil(t, svr.deviceManager)
}

func TestJtiUDPServer_Stop_nilConn(t *testing.T) {
	svr := JtiUDPServer{}
	assert.False(t, svr.stopped)
	assert.Nil(t, svr.conn)

	svr.Stop()
	assert.True(t, svr.stopped)
}

func TestJtiUDPServer_Stop(t *testing.T) {
	svr := JtiUDPServer{
		conn: &net.UDPConn{},
	}
	assert.False(t, svr.stopped)
	assert.NotNil(t, svr.conn)

	svr.Stop()
	assert.True(t, svr.stopped)
}

func TestNewDeviceFromInfo(t *testing.T) {
	// Context is only set via device info.

	info := jti.DeviceInfo{
		Type: "device-type",
		Info: "device-info",
		Tags: []string{
			"a/b:c",
		},
		Context: map[string]string{
			"device-ctx": "abc",
		},
		IDComponents: map[string]string{
			"foo": "bar",
		},
	}

	svr := JtiUDPServer{
		GlobalContext: map[string]string{},
		deviceManager: manager.NewStubDeviceManager(false),
	}

	dev, err := svr.newDeviceFromInfo(&info)
	assert.NoError(t, err)
	assert.Equal(t, "device-type", dev.Type)
	assert.Equal(t, "device-info", dev.Info)
	assert.Equal(t, "jti", dev.Handler)
	assert.Equal(t, map[string]string{
		"device-ctx": "abc",
	}, dev.Context)
	assert.Equal(t, map[string]interface{}{
		"id": map[string]string{
			"foo": "bar",
		},
	}, dev.Data)
	assert.Len(t, dev.Tags, 1)
	tag := dev.Tags[0]
	assert.Equal(t, "a", tag.Namespace)
	assert.Equal(t, "b", tag.Annotation)
	assert.Equal(t, "c", tag.Label)
}

func TestNewDeviceFromInfo2(t *testing.T) {
	// Context is only set via device info and global, no conflicts.

	info := jti.DeviceInfo{
		Type: "device-type",
		Info: "device-info",
		Tags: []string{
			"a/b:c",
		},
		Context: map[string]string{
			"device-ctx": "abc",
		},
		IDComponents: map[string]string{
			"foo": "bar",
		},
	}

	svr := JtiUDPServer{
		GlobalContext: map[string]string{
			"global-ctx": "123",
		},
		deviceManager: manager.NewStubDeviceManager(false),
	}

	dev, err := svr.newDeviceFromInfo(&info)
	assert.NoError(t, err)
	assert.Equal(t, "device-type", dev.Type)
	assert.Equal(t, "device-info", dev.Info)
	assert.Equal(t, "jti", dev.Handler)
	assert.Equal(t, map[string]string{
		"device-ctx": "abc",
		"global-ctx": "123",
	}, dev.Context)
	assert.Equal(t, map[string]interface{}{
		"id": map[string]string{
			"foo": "bar",
		},
	}, dev.Data)
	assert.Len(t, dev.Tags, 1)
	tag := dev.Tags[0]
	assert.Equal(t, "a", tag.Namespace)
	assert.Equal(t, "b", tag.Annotation)
	assert.Equal(t, "c", tag.Label)
}

func TestNewDeviceFromInfo3(t *testing.T) {
	// Context is only set via device info and global, with conflicts.

	info := jti.DeviceInfo{
		Type: "device-type",
		Info: "device-info",
		Tags: []string{
			"a/b:c",
		},
		Context: map[string]string{
			"device-ctx": "abc",
			"common":     "device-value",
		},
		IDComponents: map[string]string{
			"foo": "bar",
		},
	}

	svr := JtiUDPServer{
		GlobalContext: map[string]string{
			"global-ctx": "123",
			"common":     "global-value",
		},
		deviceManager: manager.NewStubDeviceManager(false),
	}

	dev, err := svr.newDeviceFromInfo(&info)
	assert.NoError(t, err)
	assert.Equal(t, "device-type", dev.Type)
	assert.Equal(t, "device-info", dev.Info)
	assert.Equal(t, "jti", dev.Handler)
	assert.Equal(t, map[string]string{
		"device-ctx": "abc",
		"global-ctx": "123",
		"common":     "device-value",
	}, dev.Context)
	assert.Equal(t, map[string]interface{}{
		"id": map[string]string{
			"foo": "bar",
		},
	}, dev.Data)
	assert.Len(t, dev.Tags, 1)
	tag := dev.Tags[0]
	assert.Equal(t, "a", tag.Namespace)
	assert.Equal(t, "b", tag.Annotation)
	assert.Equal(t, "c", tag.Label)
}

func TestNewDeviceFromInfo_Error(t *testing.T) {
	info := jti.DeviceInfo{
		Type: "device-type",
		Info: "device-info",
		Tags: []string{
			"a/b:c",
		},
		Context: map[string]string{
			"device-ctx": "abc",
		},
		IDComponents: map[string]string{
			"foo": "bar",
		},
	}

	svr := JtiUDPServer{
		GlobalContext: map[string]string{},
		deviceManager: manager.NewStubDeviceManager(true),
	}

	dev, err := svr.newDeviceFromInfo(&info)
	assert.Error(t, err)
	assert.Nil(t, dev)
}
