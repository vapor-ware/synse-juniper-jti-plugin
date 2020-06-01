package pkg

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/config"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/manager"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/protocol"
	"github.com/vapor-ware/synse-sdk/sdk"
)

// RunBackgroundListener is a plugin pre-run action which starts the UDP server, listening
// for incoming streamed data from Juniper equipment.
var RunBackgroundListener = sdk.PluginAction{
	Name: "run background JTI listener",
	Action: func(p *sdk.Plugin) error {

		// First, get the pre-loaded server configuration parsed from the plugin
		// configuration's dynamicRegistration block.
		serverConfig := config.Get()
		if serverConfig == nil {
			return errors.New("failed to load cached UDP server configuration")
		}

		// Create the device manager used by the server to get, create, and register devices.
		deviceManager := manager.NewPluginDeviceManager(p)

		// Create the UDP server from the configuration.
		svr := protocol.NewJtiUDPServer(serverConfig, deviceManager)

		log.Info("[jti] starting UDP server listen")
		go func() {
			if err := svr.Listen(); err != nil {
				log.WithError(err).Error("[jti] failed UDP server listen")
				panic(err)
			}
			log.Info("[jti] finished UDP server listen")
		}()

		return nil
	},
}
