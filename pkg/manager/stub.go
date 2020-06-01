package manager

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/config"
)

// StubDeviceManager implements the DeviceManager interface. It provides a stub
// device manager useful for testing.
type StubDeviceManager struct {
	withError bool
	cache     map[string]*sdk.Device
}

// NewStubDeviceManager creates a new DeviceManager for testing.
func NewStubDeviceManager(withError bool) DeviceManager {
	return &StubDeviceManager{
		withError: withError,
		cache:     make(map[string]*sdk.Device),
	}
}

// GetDevice gets an SDK Device.
func (dm *StubDeviceManager) GetDevice(id string) *sdk.Device {
	return dm.cache[id]
}

// NewDevice creates a new SDK Device.
func (dm *StubDeviceManager) NewDevice(proto *config.DeviceProto, inst *config.DeviceInstance) (*sdk.Device, error) {
	if dm.withError {
		return nil, fmt.Errorf("error creating stub device")
	}
	return sdk.NewDeviceFromConfig(proto, inst, map[string]*sdk.DeviceHandler{"jti": {}})
}

// RegisterDevice registers an SDK Device.
func (dm *StubDeviceManager) RegisterDevice(device *sdk.Device) error {
	if dm.withError {
		return fmt.Errorf("error registering stub device")
	}
	dm.cache[device.GetID()] = device
	return nil
}

// GenerateDeviceID generates a fake device ID for the given device.
func (dm *StubDeviceManager) GenerateDeviceID(device *sdk.Device) string {
	return "test-device-id"
}
