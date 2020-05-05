package manager

import (
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/config"
)

// pluginDeviceManager implements the DeviceManager interface. It provides access to the
// SDK's built-in device manager through the exposed methods on a Plugin instance.
type pluginDeviceManager struct {
	plugin *sdk.Plugin
}

// NewPluginDeviceManager creates a new DeviceManager for SDK Plugin instances.
func NewPluginDeviceManager(plugin *sdk.Plugin) DeviceManager {
	return &pluginDeviceManager{
		plugin: plugin,
	}
}

// GetDevice gets an SDK Device.
func (dm *pluginDeviceManager) GetDevice(id string) *sdk.Device {
	return dm.plugin.GetDevice(id)
}

// NewDevice creates a new SDK Device.
func (dm *pluginDeviceManager) NewDevice(proto *config.DeviceProto, inst *config.DeviceInstance) (*sdk.Device, error) {
	return dm.plugin.NewDevice(proto, inst)
}

// RegisterDevice registers an SDK Device with the backing Plugin device manager.
func (dm *pluginDeviceManager) RegisterDevice(device *sdk.Device) error {
	return dm.plugin.AddDevice(device)
}
