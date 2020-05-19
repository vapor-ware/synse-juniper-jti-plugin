package pkg

import (
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/handlers"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/outputs"
	"github.com/vapor-ware/synse-sdk/sdk"
)

// MakePlugin creates a new instance of an SDK plugin and registers
// all of the Juniper JTI-specific capabilities with the plugin.
func MakePlugin() (*sdk.Plugin, error) {
	log.Debug("[jti] making plugin")

	// Create a new plugin instance
	plugin, err := sdk.NewPlugin(
		sdk.CustomDynamicDeviceRegistration(LoadDynamicConfig),
		sdk.CustomDeviceIdentifier(DeviceIdentifier),
		sdk.PluginConfigRequired(),
		sdk.DeviceConfigOptional(),
		sdk.DynamicConfigRequired(),
	)
	if err != nil {
		return nil, err
	}

	// Register pre-run action(s) with the plugin.
	plugin.RegisterPreRunActions(
		&RunBackgroundListener,
	)

	// Register custom output types
	err = plugin.RegisterOutputs(
		&outputs.Boolean,
		&outputs.BytesCounter,
		&outputs.BytesPerSecond,
		&outputs.DecibelMilliwatts,
		&outputs.MegabitPerSecond,
		&outputs.PacketsCounter,
		&outputs.PacketsPerSecond,
		&outputs.TimeTicks,
	)
	if err != nil {
		return nil, err
	}

	// Register device handler(s)
	err = plugin.RegisterDeviceHandlers(
		&handlers.JTIDeviceHandler,
	)
	if err != nil {
		return nil, err
	}

	return plugin, nil
}
