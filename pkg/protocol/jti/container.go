package jti

import "github.com/vapor-ware/synse-sdk/sdk/output"

// IntermediaryDataContainer is a container for device info and reading data, associating
// the two related pieces prior to SDK Device creation and their subsequent association
// within the SDK Device's Data field.
type IntermediaryDataContainer struct {
	DeviceInfo *DeviceInfo
	Readings   []*output.Reading
}

// DeviceInfo is a light wrapper around basic data needed to create a new
// SDK Device.
type DeviceInfo struct {
	Type         string
	Info         string
	Tags         []string
	Context      map[string]string
	IDComponents map[string]string
}
