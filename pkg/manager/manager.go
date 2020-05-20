package manager

import (
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/config"
)

// DeviceManager provides an interface for managing SDK devices for the plugin.
//
// This interface is used as a sort of side-step around the fact that the SDK has
// a tightly controlled API and doesn't provide direct access to the internal
// device manager.
type DeviceManager interface {
	GetDevice(string) *sdk.Device
	NewDevice(*config.DeviceProto, *config.DeviceInstance) (*sdk.Device, error)
	RegisterDevice(*sdk.Device) error
	GenerateDeviceID(*sdk.Device) string
}
