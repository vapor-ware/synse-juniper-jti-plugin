package jti

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/outputs"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/protocol/jti/protos/optics"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/protocol/jti/protos/telemetry_top"
	"github.com/vapor-ware/synse-sdk/sdk/output"
)

// OpticsContext provides contextual information used to generate devices and
// readings from a JTI GPB optics message.
type OpticsContext struct {
	SensorName     string
	SystemID       string
	ComponentID    uint32
	SubComponentID uint32
}

// NewOpticsContextFromStream creates a new OpticsContext populated with values from
// the higher-level TelemetryStream GPB message associated with the Optics message.
func NewOpticsContextFromStream(ts *telemetry_top.TelemetryStream) *OpticsContext {
	return &OpticsContext{
		SensorName:     ts.GetSensorName(),
		SystemID:       ts.GetSystemId(),
		ComponentID:    ts.GetComponentId(),
		SubComponentID: ts.GetSubComponentId(),
	}
}

// Decode the Optics GPB message into a data container which can be translated into
// Synse devices and readings.
func (ctx *OpticsContext) Decode(opt *optics.Optics) ([]*IntermediaryDataContainer, error) {
	var decoded []*IntermediaryDataContainer

	if opt == nil {
		log.Info("[jti] optics decode: optics is nil, no data to collect")
		return decoded, nil
	}

	for _, info := range opt.GetOpticsDiag() {
		deviceInfo, err := ctx.MakeDeviceInfo(info)
		if err != nil {
			return nil, err
		}
		readings, err := ctx.MakeReadings(info)
		if err != nil {
			return nil, err
		}

		decoded = append(decoded, &IntermediaryDataContainer{
			DeviceInfo: deviceInfo,
			Readings:   readings,
		})
	}
	return decoded, nil
}

// MakeDeviceInfo creates a DeviceInfo corresponding to an OpticsInfos. The DeviceInfo
// is used to generate SDK devices.
func (ctx *OpticsContext) MakeDeviceInfo(info *optics.OpticsInfos) (*DeviceInfo, error) {
	if info == nil {
		return nil, errors.New("unable to load device info from optics context: nil info")
	}

	if ctx.SystemID == "" {
		return nil, errors.New("unable to load device info from optics context: context has no system ID")
	}

	ifaceName := info.GetIfName()
	if ifaceName == "" {
		return nil, errors.New("unable to load device info from optics context: info has no name")
	}

	// FIXME (etd): is this an interface? is this an "optic"? what do we call this device?
	return &DeviceInfo{
		Type: "interface",
		Info: fmt.Sprintf("%s interface %s", ctx.SystemID, ifaceName),
		Tags: []string{
			"vapor/networking:interface",
		},
		Context: map[string]string{
			"component_id": fmt.Sprint(ctx.ComponentID),
			"name":         ifaceName,
			"system_id":    ctx.SystemID,
			"sensor_name":  ctx.SensorName,
		},
		IDComponents: map[string]string{
			"sys":  ctx.SystemID,
			"if":   ifaceName,
			"cid":  fmt.Sprint(ctx.ComponentID),
			"scid": fmt.Sprint(ctx.SubComponentID),
		},
	}, nil
}

// MakeReadings creates device readings for an OpticsInfos message. The message contains many data
// points, all of which are translated into Synse readings.
func (ctx *OpticsContext) MakeReadings(info *optics.OpticsInfos) ([]*output.Reading, error) {
	stats := info.GetOpticsDiagStats()

	var readings = []*output.Reading{
		output.Number.MakeReading(stats.GetOpticsType()).WithContext(map[string]string{
			"metric": "optics_type",
		}),
		output.Temperature.MakeReading(stats.GetModuleTemp()).WithContext(map[string]string{
			"metric": "module_temp",
		}),
		output.Temperature.MakeReading(stats.GetModuleTempHighAlarmThreshold()).WithContext(map[string]string{
			"metric": "module_temp_high_alarm_threshold",
		}),
		output.Temperature.MakeReading(stats.GetModuleTempLowAlarmThreshold()).WithContext(map[string]string{
			"metric": "module_temp_low_alarm_threshold",
		}),
		output.Temperature.MakeReading(stats.GetModuleTempHighWarningThreshold()).WithContext(map[string]string{
			"metric": "module_temp_high_warning_threshold",
		}),
		output.Temperature.MakeReading(stats.GetModuleTempLowWarningThreshold()).WithContext(map[string]string{
			"metric": "module_temp_low_warning_threshold",
		}),
		outputs.DecibelMilliwatts.MakeReading(stats.GetLaserOutputPowerHighAlarmThresholdDbm()).WithContext(map[string]string{
			"metric": "laser_output_power_high_alarm_threshold_dbm",
		}),
		outputs.DecibelMilliwatts.MakeReading(stats.GetLaserOutputPowerLowAlarmThresholdDbm()).WithContext(map[string]string{
			"metric": "laser_output_power_low_alarm_threshold_dbm",
		}),
		outputs.DecibelMilliwatts.MakeReading(stats.GetLaserOutputPowerHighWarningThresholdDbm()).WithContext(map[string]string{
			"metric": "laser_output_power_high_warning_threshold_dbm",
		}),
		outputs.DecibelMilliwatts.MakeReading(stats.GetLaserOutputPowerLowWarningThresholdDbm()).WithContext(map[string]string{
			"metric": "laser_output_power_low_warning_threshold_dbm",
		}),
		outputs.DecibelMilliwatts.MakeReading(stats.GetLaserRxPowerHighAlarmThresholdDbm()).WithContext(map[string]string{
			"metric": "laser_rx_power_high_alarm_threshold_dbm",
		}),
		outputs.DecibelMilliwatts.MakeReading(stats.GetLaserRxPowerLowAlarmThresholdDbm()).WithContext(map[string]string{
			"metric": "laser_rx_power_low_alarm_threshold_dbm",
		}),
		outputs.DecibelMilliwatts.MakeReading(stats.GetLaserRxPowerHighWarningThresholdDbm()).WithContext(map[string]string{
			"metric": "laser_rx_power_high_warning_threshold_dbm",
		}),
		outputs.DecibelMilliwatts.MakeReading(stats.GetLaserRxPowerLowWarningThresholdDbm()).WithContext(map[string]string{
			"metric": "laser_rx_power_low_warning_threshold_dbm",
		}),
		// FIXME (etd): I don't know what these units should be, so just sticking with "number" for now.
		output.Number.MakeReading(stats.GetLaserBiasCurrentHighAlarmThreshold()).WithContext(map[string]string{
			"metric": "laser_bias_current_high_alarm_threshold",
		}),
		output.Number.MakeReading(stats.GetLaserBiasCurrentLowAlarmThreshold()).WithContext(map[string]string{
			"metric": "laser_bias_current_low_alarm_threshold",
		}),
		output.Number.MakeReading(stats.GetLaserBiasCurrentHighWarningThreshold()).WithContext(map[string]string{
			"metric": "laser_bias_current_high_warning_threshold",
		}),
		output.Number.MakeReading(stats.GetLaserBiasCurrentLowWarningThreshold()).WithContext(map[string]string{
			"metric": "laser_bias_current_low_warning_threshold",
		}),
		outputs.Boolean.MakeReading(stats.GetModuleTempHighAlarm()).WithContext(map[string]string{
			"metric": "module_temp_high_alarm",
		}),
		outputs.Boolean.MakeReading(stats.GetModuleTempLowAlarm()).WithContext(map[string]string{
			"metric": "module_temp_low_alarm",
		}),
		outputs.Boolean.MakeReading(stats.GetModuleTempHighWarning()).WithContext(map[string]string{
			"metric": "module_temp_high_warning",
		}),
		outputs.Boolean.MakeReading(stats.GetModuleTempLowWarning()).WithContext(map[string]string{
			"metric": "module_temp_low_warning",
		}),
	}

	var laneStats []*output.Reading
	for _, stat := range info.OpticsDiagStats.GetOpticsLaneDiagStats() {
		laneNumber := fmt.Sprint(stat.GetLaneNumber())

		var q = []*output.Reading{
			output.Temperature.MakeReading(stat.GetLaneLaserTemperature()).WithContext(map[string]string{
				"lane_number": laneNumber,
				"metric":      "lane_laser_temperature",
			}),
			outputs.DecibelMilliwatts.MakeReading(stat.GetLaneLaserOutputPowerDbm()).WithContext(map[string]string{
				"lane_number": laneNumber,
				"metric":      "lane_laser_output_power_dbm",
			}),
			outputs.DecibelMilliwatts.MakeReading(stat.GetLaneLaserReceiverPowerDbm()).WithContext(map[string]string{
				"lane_number": laneNumber,
				"metric":      "lane_laser_receiver_power_dbm",
			}),
			// TODO (etd): don't know the unit, so just using Number here
			output.Number.MakeReading(stat.GetLaneLaserBiasCurrent()).WithContext(map[string]string{
				"lane_number": laneNumber,
				"metric":      "lane_laser_bias_current",
			}),
			outputs.Boolean.MakeReading(stat.GetLaneLaserOutputPowerHighAlarm()).WithContext(map[string]string{
				"lane_number": laneNumber,
				"metric":      "lane_laser_output_power_high_alarm",
			}),
			outputs.Boolean.MakeReading(stat.GetLaneLaserOutputPowerLowAlarm()).WithContext(map[string]string{
				"lane_number": laneNumber,
				"metric":      "lane_laser_output_power_low_alarm",
			}),
			outputs.Boolean.MakeReading(stat.GetLaneLaserOutputPowerHighWarning()).WithContext(map[string]string{
				"lane_number": laneNumber,
				"metric":      "lane_laser_output_power_high_warning",
			}),
			outputs.Boolean.MakeReading(stat.GetLaneLaserOutputPowerLowWarning()).WithContext(map[string]string{
				"lane_number": laneNumber,
				"metric":      "lane_laser_output_power_low_warning",
			}),
			outputs.Boolean.MakeReading(stat.GetLaneLaserReceiverPowerHighAlarm()).WithContext(map[string]string{
				"lane_number": laneNumber,
				"metric":      "lane_laser_receiver_power_high_alarm",
			}),
			outputs.Boolean.MakeReading(stat.GetLaneLaserReceiverPowerLowAlarm()).WithContext(map[string]string{
				"lane_number": laneNumber,
				"metric":      "lane_laser_receiver_power_low_alarm",
			}),
			outputs.Boolean.MakeReading(stat.GetLaneLaserReceiverPowerHighWarning()).WithContext(map[string]string{
				"lane_number": laneNumber,
				"metric":      "lane_laser_receiver_power_high_warning",
			}),
			outputs.Boolean.MakeReading(stat.GetLaneLaserReceiverPowerLowWarning()).WithContext(map[string]string{
				"lane_number": laneNumber,
				"metric":      "lane_laser_receiver_power_low_warning",
			}),
			outputs.Boolean.MakeReading(stat.GetLaneLaserBiasCurrentHighAlarm()).WithContext(map[string]string{
				"lane_number": laneNumber,
				"metric":      "lane_laser_bias_current_high_alarm",
			}),
			outputs.Boolean.MakeReading(stat.GetLaneLaserBiasCurrentLowAlarm()).WithContext(map[string]string{
				"lane_number": laneNumber,
				"metric":      "lane_laser_bias_current_low_alarm",
			}),
			outputs.Boolean.MakeReading(stat.GetLaneLaserBiasCurrentHighWarning()).WithContext(map[string]string{
				"lane_number": laneNumber,
				"metric":      "lane_laser_bias_current_high_warning",
			}),
			outputs.Boolean.MakeReading(stat.GetLaneLaserBiasCurrentLowWarning()).WithContext(map[string]string{
				"lane_number": laneNumber,
				"metric":      "lane_laser_bias_current_low_warning",
			}),
			outputs.Boolean.MakeReading(stat.GetLaneTxLossOfSignalAlarm()).WithContext(map[string]string{
				"lane_number": laneNumber,
				"metric":      "lane_tx_loss_of_signal_alarm",
			}),
			outputs.Boolean.MakeReading(stat.GetLaneRxLossOfSignalAlarm()).WithContext(map[string]string{
				"lane_number": laneNumber,
				"metric":      "lane_rx_loss_of_signal_alarm",
			}),
			outputs.Boolean.MakeReading(stat.GetLaneTxLaserDisabledAlarm()).WithContext(map[string]string{
				"lane_number": laneNumber,
				"metric":      "lane_tx_laser_disabled_alarm",
			}),
		}
		laneStats = append(laneStats, q...)
	}

	readings = append(readings, laneStats...)

	return readings, nil
}
