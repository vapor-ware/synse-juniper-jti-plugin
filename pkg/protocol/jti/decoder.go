package jti

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/manager"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/protocol/jti/protos/optics"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/protocol/jti/protos/port"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/protocol/jti/protos/telemetry_top"
)

// JuniperJTIDecoder is used to decode bytes from the UDP data stream
// into its appropriate JTI protobuf type and from there into a format
// which is usable for Synse devices and readings.
//
// For details on Juniper sensors, see:
// https://www.juniper.net/documentation/en_US/junos/topics/reference/general/junos-telemetry-interface-grpc-sensors.html
type JuniperJTIDecoder struct {
	deviceManager manager.DeviceManager
}

// NewJTIDecoder creates a new JuniperJTIDecoder.
func NewJTIDecoder(deviceManager manager.DeviceManager) *JuniperJTIDecoder {
	return &JuniperJTIDecoder{
		deviceManager: deviceManager,
	}
}

// Decode the given bytes into a format consumable by the Synse platform.
func (decoder *JuniperJTIDecoder) Decode(buffer []byte) ([]*IntermediaryDataContainer, error) {
	if decoder.deviceManager == nil {
		return nil, errors.New("JTI decoder does not have a device manager defined")
	}

	ts := &telemetry_top.TelemetryStream{}
	if err := proto.Unmarshal(buffer, ts); err != nil {
		return nil, err
	}

	var decoded []*IntermediaryDataContainer

	if proto.HasExtension(ts.Enterprise, telemetry_top.E_JuniperNetworks) {
		jnsIface, err := proto.GetExtension(ts.Enterprise, telemetry_top.E_JuniperNetworks)
		if err != nil {
			log.WithError(err).Error("[jti] failed to get extension")
			return nil, err
		}

		switch jns := jnsIface.(type) {
		case *telemetry_top.JuniperNetworksSensors:

			if proto.HasExtension(jns, optics.E_JnprOpticsExt) {
				/*
					OPTICS
				*/
				opticsIface, err := proto.GetExtension(jns, optics.E_JnprOpticsExt)
				if err != nil {
					log.WithError(err).Error("[jti] failed to get extension")
					return nil, err
				}

				switch opt := opticsIface.(type) {
				case *optics.Optics:
					res, err := NewOpticsContextFromStream(ts).Decode(opt)
					if err != nil {
						return nil, err
					}
					decoded = append(decoded, res...)
				default:
					log.Error("[jti] found no matching optics iface")
					return nil, fmt.Errorf("found no matching optics interface")
				}

			} else if proto.HasExtension(jns, port.E_JnprInterfaceExt) {
				/*
					PORT INTERFACE
				*/
				portIface, err := proto.GetExtension(jns, port.E_JnprInterfaceExt)
				if err != nil {
					log.WithError(err).Error("[jti] failed to get extension")
					return nil, err
				}

				switch p := portIface.(type) {
				case *port.Port:
					res, err := NewPortContextFromStream(ts).Decode(p)
					if err != nil {
						return nil, err
					}
					decoded = append(decoded, res...)
				default:
					log.Error("[jti] found no matching port iface")
					return nil, fmt.Errorf("found no matching port interface")
				}
			} else {
				/*
					UNKNOWN
				*/
				// TODO (etd): Figure out how to get an identifier for the extension to log.
				// 	 It is not immediately intuitive how to do this.
				log.Info("[jti] received message with extension not currently supported by the plugin")
			}

		default:
			log.WithFields(log.Fields{
				"ext": jns,
			}).Warning("[jti] unsupported JTI protobuf extension")
		}

	} else {
		log.Warning("[jti] message does not provide juniper network extension")
	}

	return decoded, nil
}
