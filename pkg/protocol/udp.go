package protocol

import (
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/manager"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/protocol/jti"
	"github.com/vapor-ware/synse-sdk/sdk/config"
)

const (
	// ReadingKey is the key into a device's Data field which stores the
	// device reading data.
	ReadingKey = "_device_readings"
)

// JtiUDPServer is the UDP server for collecting streamed JTI data over UDP.
type JtiUDPServer struct {
	Address    string
	BufferSize uint64

	stopped       bool
	conn          *net.UDPConn
	decoder       *jti.JuniperJTIDecoder
	deviceManager manager.DeviceManager
}

// NewJtiUDPServer creates a new instance of a JtiUDPServer.
func NewJtiUDPServer(address string, deviceManager manager.DeviceManager) *JtiUDPServer {
	return &JtiUDPServer{
		Address:       address,
		BufferSize:    64 * 1024, // 64kb, max size of UDP datagram.
		decoder:       jti.NewJTIDecoder(deviceManager),
		deviceManager: deviceManager,
	}
}

// Connect creates the UDP server connection.
func (server *JtiUDPServer) Connect() error {
	if server.conn != nil {
		log.WithFields(log.Fields{
			"conn": server.conn,
		}).Debug("[jti] UDP server already connected")
		return nil
	}

	var (
		network string
		address string
	)

	splt := strings.SplitN(server.Address, "://", 2)
	if len(splt) != 2 {
		network = "udp"
		address = splt[0]
	} else {
		network = splt[0]
		address = splt[1]
	}

	addr, err := net.ResolveUDPAddr(network, address)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP(network, addr)
	if err != nil {
		return err
	}
	server.conn = conn

	return nil
}

// Stop the UDP server from running and close the server connection.
func (server *JtiUDPServer) Stop() {
	server.stopped = true

	if server.conn != nil {
		_ = server.conn.Close()
		server.conn = nil
	}
}

// Listen is the entry point for the server run. It will listen for incoming packets
// and attempt to decode them into device readings.
//
// If new devices are found, it will add them to the device manager. All readings are
// associated with a device via the device's Data field.
func (server *JtiUDPServer) Listen() error {
	buf := make([]byte, server.BufferSize)

	if err := server.Connect(); err != nil {
		log.WithError(err).Error("[jti] error creating UDP connection")
		return err
	}

	log.WithFields(log.Fields{
		"address": server.Address,
		"buffer":  server.BufferSize,
	}).Info("[jti] listening...")

	for !server.stopped {
		n, _, err := server.conn.ReadFrom(buf)
		if err != nil {
			log.WithError(err).Error("[jti] error reading from UDP connection")
			return err
		}

		data, err := server.decoder.Decode(buf[:n])
		if err != nil {
			log.WithError(err).Warning("[jti] failed to decode payload into readings - discarding")
			continue
		}

		for _, d := range data {
			dev, err := server.deviceManager.NewDevice(
				&config.DeviceProto{
					Type:    d.DeviceInfo.Type,
					Context: d.DeviceInfo.Context,
					Tags:    d.DeviceInfo.Tags,
					Data: map[string]interface{}{
						"id": d.DeviceInfo.IDComponents,
					},
					Handler: "jti",
				},
				&config.DeviceInstance{
					Info: d.DeviceInfo.Info,
				},
			)
			if err != nil {
				log.WithFields(log.Fields{
					"err":  err,
					"type": d.DeviceInfo.Type,
					"info": d.DeviceInfo.Info,
					"id":   d.DeviceInfo.IDComponents,
					"ctx":  d.DeviceInfo.Context,
					"tags": d.DeviceInfo.Tags,
				}).Error("[jti] failed to create a new device")
				return err
			}

			// Attempt to get the device. If the device does not yet exist, register it
			// with the plugin.
			deviceID := server.deviceManager.GenerateDeviceID(dev)
			device := server.deviceManager.GetDevice(deviceID)
			if device == nil {
				log.WithFields(log.Fields{
					"id": deviceID,
				}).Info("[jti] device with ID does not exist - creating new device")
				if err := server.deviceManager.RegisterDevice(dev); err != nil {
					log.WithFields(log.Fields{
						"err":  err,
						"id":   deviceID,
						"info": dev.Info,
						"ctx":  dev.Context,
						"type": dev.Type,
					}).Error("[jti] failed to register new device")
					return err
				}

				// Since the device is now registered, we can use the Device reference
				// to add the readings to.
				device = dev
			}

			// Add the readings to the device data.
			device.Data[ReadingKey] = d.Readings
		}
	}

	return nil
}
