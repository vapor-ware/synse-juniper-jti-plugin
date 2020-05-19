package pkg

import (
	"errors"
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/config"
	"github.com/vapor-ware/synse-sdk/sdk"
)

// Errors relating to loading and executing upon the plugin dynamic configs.
var (
	ErrTooManyServers = errors.New("too many UDP servers configured")
)

var configCount = 0

// LoadDynamicConfig loads the dynamic configuration provided in the plugin
// config, creates the necessary data, and starts receiving data for the
// configured data sources.
func LoadDynamicConfig(data map[string]interface{}) ([]*sdk.Device, error) {
	log.Debug("[jti] loading dynamic registration config")
	if configCount > 0 {
		log.Error("[jti] too many UDP servers configured; only one should be defined in the plugin dynamicRegistration configuration")
		return nil, ErrTooManyServers
	}
	configCount++

	serverConfig, err := config.Load(data)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"data": data,
		}).Error("[jti] failed to load configuration data")
		return nil, err
	}

	// The server config is used to generate a UDP server to listen for the
	// streamed telemetry data. We need a reference to the Plugin in order
	// to create this server, as it needs the reference in order to execute
	// callbacks to register any new devices that it finds in the data stream.
	//
	// Since this plugin handler is passed as an initialization option to the
	// plugin constructor, we cannot get a reference to the plugin within the
	// scope of this functions. By caching the loaded config, we can defer its
	// use until a plugin PreRunAction, where we will have the necessary
	// reference to the plugin.
	log.WithFields(log.Fields{
		"config": serverConfig,
	}).Debug("[jti] caching server config")
	config.Set(serverConfig)

	return []*sdk.Device{}, nil
}

// DeviceIdentifier is the custom device identifier function for the JTI plugin.
//
// Since all devices are created dynamically at runtime, we have some guarantees of
// what fields exist in the data field and how they are structured. The runtime device
// loader will always put all fields pertaining to the ID of the plugin into an map
// under the "id" key.
func DeviceIdentifier(data map[string]interface{}) string {
	var identifier string

	idComponents, exists := data["id"]
	if !exists {
		panic("device does not contain the expected ID info in its data")
	}

	components, ok := idComponents.(map[string]string)
	if !ok {
		panic("device ID info is not a map of string:string")
	}

	// To ensure that we get the same identifier reliably, we want to make sure
	// we append the components reliably, so we will sort the keys.
	var keys []string
	for k := range components {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		// Instead of implementing our own type checking and casting, just
		// use Sprint. Note that this may be meaningless for complex types.
		identifier += fmt.Sprint(components[key])
	}

	log.WithField("identifier", identifier).Debug("[jti] constructed device identifier")
	return identifier
}
