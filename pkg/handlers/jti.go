package handlers

import (
	"errors"
	"fmt"

	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/protocol"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
)

// JTIDeviceHandler is the device handler for all devices found by the Juniper
// JTI Plugin using GPB+UDP.
//
// A single device handler may be used because each device has its readings
// pre-built and associated with it elsewhere during the process of receiving
// and parsing the Juniper UDP packets.
//
// All devices are built at runtime, so there should never be a need to reference
// this handler directly in a device configuration. Each device that gets handled
// should have its associated reading data specified in its Data field.
//
// This plugin only supports reads since the data is collected from a unidirectional
// UDP stream.
var JTIDeviceHandler = sdk.DeviceHandler{
	Name: "jti",
	Read: jtiDeviceRead,
}

// jtiDeviceRead implements the read capability for the JTIDeviceHandler.
func jtiDeviceRead(device *sdk.Device) ([]*output.Reading, error) {
	r, exists := device.Data[protocol.ReadingKey]
	if !exists {
		return nil, errors.New("error reading device: expected readings key not found in device data")
	}
	readings, ok := r.([]*output.Reading)
	if !ok {
		return nil, fmt.Errorf("error reading device: unexpected reading data type %T", r)
	}
	return readings, nil
}
