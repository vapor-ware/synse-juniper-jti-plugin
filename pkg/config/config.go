package config

import (
	"errors"

	"github.com/mitchellh/mapstructure"
)

// Errors related to loading and parsing data source configurations.
var (
	ErrNoAddress = errors.New("data source configuration does not define required 'address' value")
)

var serverConfig *ServerConfig

// ServerConfig is the configuration for the UDP server. This is loaded in from
// the plugin dynamic registration block.
type ServerConfig struct {

	// Address for the UDP server to listen on. This should be a string specifying
	// the IP/hostname and port. The protocol prefix may be one of: "udp", "udp4",
	// "udp6". If unspecified, "udp" is used.
	Address string
}

// Load the configuration for the plugin's UDP server which will listen for the
// streamed JTI telemetry data.
//
// This also performs basic validation of the data being loaded. It ensures
// that required fields are not empty.
func Load(raw map[string]interface{}) (*ServerConfig, error) {
	var cfg ServerConfig
	if err := mapstructure.Decode(raw, &cfg); err != nil {
		return nil, err
	}

	if cfg.Address == "" {
		return nil, ErrNoAddress
	}
	return &cfg, nil
}

// Set the global server config. This config defines configuration options for the UDP
// server that will be set up to listen for incoming data streams from Juniper equipment.
//
// This is set globally because of restrictions around variable scoping during plugin
// initialization and a fairly locked down interface for interacting with the plugin SDK.
func Set(cfg *ServerConfig) {
	serverConfig = cfg
}

// Get the global server config. This config defines configuration options for the UDP
//// server that will be set up to listen for incoming data streams from Juniper equipment.
//
// If this configuration has not been set yet, nil is returned.
func Get() *ServerConfig {
	return serverConfig
}
