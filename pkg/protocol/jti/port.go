package jti

import (
	"errors"
	"fmt"

	"github.com/prometheus/common/log"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/outputs"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/protocol/jti/protos/port"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/protocol/jti/protos/telemetry_top"
	"github.com/vapor-ware/synse-sdk/sdk/output"
)

// PortContext provides contextual information used to generate devices and
// readings from a JTI GPB port message.
type PortContext struct {
	SensorName     string
	SystemID       string
	ComponentID    uint32
	SubComponentID uint32
}

// NewPortContextFromStream creates a new PortContext populated with values from
// the higher-level TelemetryStream GPB message associated with the Port message.
func NewPortContextFromStream(ts *telemetry_top.TelemetryStream) *PortContext {
	return &PortContext{
		SensorName:     ts.GetSensorName(),
		SystemID:       ts.GetSystemId(),
		ComponentID:    ts.GetComponentId(),
		SubComponentID: ts.GetSubComponentId(),
	}
}

// Decode the Port GPB message into a data container which can be translated into
// Synse devices and readings.
func (ctx *PortContext) Decode(prt *port.Port) ([]*IntermediaryDataContainer, error) {
	var decoded []*IntermediaryDataContainer

	if prt == nil {
		log.Info("[jti] port decode: port is nil, no data to collect")
		return decoded, nil
	}

	for _, info := range prt.GetInterfaceStats() {
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

// MakeDeviceInfo creates a DeviceInfo corresponding to an InterfaceInfos. The DeviceInfo
// is used to generate SDK devices.
func (ctx *PortContext) MakeDeviceInfo(iface *port.InterfaceInfos) (*DeviceInfo, error) {
	if iface == nil {
		return nil, errors.New("unable to load device info from port context: nil interface info")
	}

	if ctx.SystemID == "" {
		return nil, errors.New("unable to load device info from port context: context has no system ID")
	}

	ifaceName := iface.GetIfName()
	if ifaceName == "" {
		return nil, errors.New("unable to load device info from port context: interface has no name")
	}

	return &DeviceInfo{
		Type: "interface",
		Info: fmt.Sprintf("%s interface %s", ctx.SystemID, ifaceName),
		Tags: []string{
			"vapor/networking:interface",
		},
		Context: map[string]string{
			"component_id":   fmt.Sprint(ctx.ComponentID),
			"name":           ifaceName,
			"system_id":      ctx.SystemID,
			"sensor_name":    ctx.SensorName,
			"parent_ae_name": iface.GetParentAeName(),
		},
		IDComponents: map[string]string{
			"sys":  ctx.SystemID,
			"if":   ifaceName,
			"cid":  fmt.Sprint(ctx.ComponentID),
			"scid": fmt.Sprint(ctx.SubComponentID),
		},
	}, nil
}

// MakeReadings creates device readings for an InterfaceInfos message. The message contains many data
// points, all of which are translated into Synse readings.
func (ctx *PortContext) MakeReadings(iface *port.InterfaceInfos) ([]*output.Reading, error) {
	// FIXME (etd): This is a first pass at what interface reading data could look like. This
	//   is still a work in progress and should be reviewed + refined.

	// TODO (etd): maybe there is a more programmatic way of being able to construct these? e.g.
	// 	 via reflection?

	var readings = []*output.Reading{
		// -*- Bytes Counter Outputs -*-
		outputs.BytesCounter.MakeReading(iface.IngressStats.GetIfOctets()).WithContext(map[string]string{
			"direction": "ingress",
			"metric":    "if_octets",
		}),

		outputs.BytesCounter.MakeReading(iface.EgressStats.GetIfOctets()).WithContext(map[string]string{
			"direction": "egress",
			"metric":    "if_octets",
		}),

		// -*- Bytes per Second Outputs -*-
		outputs.PacketsPerSecond.MakeReading(iface.IngressStats.GetIf_1SecOctets()).WithContext(map[string]string{
			"direction": "ingress",
			"metric":    "if_1sec_octets",
		}),

		outputs.PacketsPerSecond.MakeReading(iface.EgressStats.GetIf_1SecOctets()).WithContext(map[string]string{
			"direction": "egress",
			"metric":    "if_1sec_octets",
		}),

		// -*- Megabits per second Outputs -*-
		outputs.MegabitPerSecond.MakeReading(iface.GetIfHighSpeed()).WithContext(map[string]string{
			"metric": "if_high_speed",
		}),

		// -*- Number Outputs -*-
		output.Number.MakeReading(iface.IngressErrors.GetIfInFifoErrors()).WithContext(map[string]string{
			"direction": "ingress",
			"metric":    "if_in_fifo_errors",
		}),
		output.Number.MakeReading(iface.IngressErrors.GetIfInResourceErrors()).WithContext(map[string]string{
			"direction": "ingress",
			"metric":    "if_in_resource_errors",
		}),

		// -*- Packets Counter Outputs -*-
		outputs.PacketsCounter.MakeReading(iface.IngressStats.GetIfPkts()).WithContext(map[string]string{
			"direction": "ingress",
			"metric":    "if_pkts",
		}),
		outputs.PacketsCounter.MakeReading(iface.IngressStats.GetIfUcPkts()).WithContext(map[string]string{
			"direction": "ingress",
			"metric":    "if_uc_pkts",
		}),
		outputs.PacketsCounter.MakeReading(iface.IngressStats.GetIfMcPkts()).WithContext(map[string]string{
			"direction": "ingress",
			"metric":    "if_mc_pkts",
		}),
		outputs.PacketsCounter.MakeReading(iface.IngressStats.GetIfBcPkts()).WithContext(map[string]string{
			"direction": "ingress",
			"metric":    "if_bc_pkts",
		}),
		outputs.PacketsCounter.MakeReading(iface.IngressStats.GetIfError()).WithContext(map[string]string{
			"direction": "ingress",
			"metric":    "if_error",
		}),
		outputs.PacketsCounter.MakeReading(iface.IngressStats.GetIfPausePkts()).WithContext(map[string]string{
			"direction": "ingress",
			"metric":    "if_pause_pkts",
		}),
		outputs.PacketsCounter.MakeReading(iface.IngressStats.GetIfUnknownProtoPkts()).WithContext(map[string]string{
			"direction": "ingress",
			"metric":    "if_unknown_proto_pkts",
		}),

		outputs.PacketsCounter.MakeReading(iface.EgressStats.GetIfPkts()).WithContext(map[string]string{
			"direction": "egress",
			"metric":    "if_pkts",
		}),
		outputs.PacketsCounter.MakeReading(iface.EgressStats.GetIfUcPkts()).WithContext(map[string]string{
			"direction": "egress",
			"metric":    "if_uc_pkts",
		}),
		outputs.PacketsCounter.MakeReading(iface.EgressStats.GetIfMcPkts()).WithContext(map[string]string{
			"direction": "egress",
			"metric":    "if_mc_pkts",
		}),
		outputs.PacketsCounter.MakeReading(iface.EgressStats.GetIfBcPkts()).WithContext(map[string]string{
			"direction": "egress",
			"metric":    "if_bc_pkts",
		}),
		outputs.PacketsCounter.MakeReading(iface.EgressStats.GetIfError()).WithContext(map[string]string{
			"direction": "egress",
			"metric":    "if_error",
		}),
		outputs.PacketsCounter.MakeReading(iface.EgressStats.GetIfPausePkts()).WithContext(map[string]string{
			"direction": "egress",
			"metric":    "if_pause_pkts",
		}),
		outputs.PacketsCounter.MakeReading(iface.EgressStats.GetIfUnknownProtoPkts()).WithContext(map[string]string{
			"direction": "egress",
			"metric":    "if_unknown_proto_pkts",
		}),

		outputs.PacketsCounter.MakeReading(iface.EgressErrors.GetIfErrors()).WithContext(map[string]string{
			"direction": "egress",
			"metric":    "if_errors",
		}),
		outputs.PacketsCounter.MakeReading(iface.EgressErrors.GetIfDiscards()).WithContext(map[string]string{
			"direction": "egress",
			"metric":    "if_discards",
		}),

		outputs.PacketsCounter.MakeReading(iface.IngressErrors.GetIfErrors()).WithContext(map[string]string{
			"direction": "ingress",
			"metric":    "if_errors",
		}),
		outputs.PacketsCounter.MakeReading(iface.IngressErrors.GetIfInQdrops()).WithContext(map[string]string{
			"direction": "ingress",
			"metric":    "if_in_qdrops",
		}),
		outputs.PacketsCounter.MakeReading(iface.IngressErrors.GetIfInFrameErrors()).WithContext(map[string]string{
			"direction": "ingress",
			"metric":    "if_in_frame_errors",
		}),
		outputs.PacketsCounter.MakeReading(iface.IngressErrors.GetIfDiscards()).WithContext(map[string]string{
			"direction": "ingress",
			"metric":    "if_discards",
		}),
		outputs.PacketsCounter.MakeReading(iface.IngressErrors.GetIfInRunts()).WithContext(map[string]string{
			"direction": "ingress",
			"metric":    "if_in_runts",
		}),
		outputs.PacketsCounter.MakeReading(iface.IngressErrors.GetIfInL3Incompletes()).WithContext(map[string]string{
			"direction": "ingress",
			"metric":    "if_in_l3_incompletes",
		}),
		outputs.PacketsCounter.MakeReading(iface.IngressErrors.GetIfInL2ChanErrors()).WithContext(map[string]string{
			"direction": "ingress",
			"metric":    "if_in_l2chan_errors",
		}),
		outputs.PacketsCounter.MakeReading(iface.IngressErrors.GetIfInL2MismatchTimeouts()).WithContext(map[string]string{
			"direction": "ingress",
			"metric":    "if_in_l2_mismatch_timeouts",
		}),

		// -*- Packets per Second Outputs -*-
		outputs.PacketsPerSecond.MakeReading(iface.IngressStats.GetIf_1SecPkts()).WithContext(map[string]string{
			"direction": "ingress",
			"metric":    "if_1sec_pkts",
		}),

		outputs.PacketsPerSecond.MakeReading(iface.EgressStats.GetIf_1SecPkts()).WithContext(map[string]string{
			"direction": "egress",
			"metric":    "if_1sec_pkts",
		}),

		// -*- Timestamp Outputs -*-
		output.Timestamp.MakeReading(iface.GetInitTime()).WithContext(map[string]string{
			"metric": "init_time",
		}),

		// -*- Time Tick Outputs -*-
		outputs.TimeTicks.MakeReading(iface.GetIfLastChange()).WithContext(map[string]string{
			"metric": "if_last_change",
		}),

		// -*- Status Outputs -*-
		output.Status.MakeReading(iface.GetIfAdministrationStatus()).WithContext(map[string]string{
			"metric": "if_administration_status",
		}),
		output.Status.MakeReading(iface.GetIfOperationalStatus()).WithContext(map[string]string{
			"metric": "if_operational_status",
		}),

		// -*- String Outputs -*-
		output.String.MakeReading(iface.GetIfDescription()).WithContext(map[string]string{
			"metric": "if_description",
		}),
		output.String.MakeReading(iface.GetParentAeName()).WithContext(map[string]string{
			"metric": "parent_ae_name",
		}),
	}

	var ingresQueueStats []*output.Reading
	for _, qstat := range iface.GetIngressQueueInfo() {
		queueNumber := fmt.Sprint(qstat.GetQueueNumber())

		var q = []*output.Reading{
			outputs.PacketsCounter.MakeReading(qstat.GetPackets()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "ingress",
				"metric":       "packets",
			}),
			outputs.PacketsCounter.MakeReading(qstat.GetTailDropPackets()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "ingress",
				"metric":       "tail_drop_packets",
			}),
			outputs.PacketsCounter.MakeReading(qstat.GetRlDropPackets()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "ingress",
				"metric":       "rl_drop_packets",
			}),
			outputs.PacketsCounter.MakeReading(qstat.GetRedDropPackets()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "ingress",
				"metric":       "red_drop_packets",
			}),
			outputs.PacketsCounter.MakeReading(qstat.GetAvgBufferOccupancy()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "ingress",
				"metric":       "avg_buffer_occupancy",
			}),
			outputs.PacketsCounter.MakeReading(qstat.GetCurBufferOccupancy()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "ingress",
				"metric":       "cur_buffer_occupancy",
			}),
			outputs.PacketsCounter.MakeReading(qstat.GetPeakBufferOccupancy()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "ingress",
				"metric":       "peak_buffer_occupancy",
			}),
			outputs.BytesCounter.MakeReading(qstat.GetBytes()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "ingress",
				"metric":       "bytes",
			}),
			outputs.BytesCounter.MakeReading(qstat.GetRlDropBytes()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "ingress",
				"metric":       "rl_drop_bytes",
			}),
			outputs.BytesCounter.MakeReading(qstat.GetRedDropBytes()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "ingress",
				"metric":       "red_drop_bytes",
			}),
			output.Number.MakeReading(qstat.GetAllocatedBufferSize()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "ingress",
				"metric":       "allocated_buffer_size",
			}),
		}
		ingresQueueStats = append(ingresQueueStats, q...)
	}

	var egressQueueStats []*output.Reading
	for _, qstat := range iface.GetEgressQueueInfo() {
		queueNumber := fmt.Sprint(qstat.GetQueueNumber())

		var q = []*output.Reading{
			outputs.PacketsCounter.MakeReading(qstat.GetPackets()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "egress",
				"metric":       "packets",
			}),
			outputs.PacketsCounter.MakeReading(qstat.GetTailDropPackets()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "egress",
				"metric":       "tail_drop_packets",
			}),
			outputs.PacketsCounter.MakeReading(qstat.GetRlDropPackets()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "egress",
				"metric":       "rl_drop_packets",
			}),
			outputs.PacketsCounter.MakeReading(qstat.GetRedDropPackets()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "egress",
				"metric":       "red_drop_packets",
			}),
			outputs.PacketsCounter.MakeReading(qstat.GetAvgBufferOccupancy()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "egress",
				"metric":       "avg_buffer_occupancy",
			}),
			outputs.PacketsCounter.MakeReading(qstat.GetCurBufferOccupancy()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "egress",
				"metric":       "cur_buffer_occupancy",
			}),
			outputs.PacketsCounter.MakeReading(qstat.GetPeakBufferOccupancy()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "egress",
				"metric":       "peak_buffer_occupancy",
			}),
			outputs.BytesCounter.MakeReading(qstat.GetBytes()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "egress",
				"metric":       "bytes",
			}),
			outputs.BytesCounter.MakeReading(qstat.GetRlDropBytes()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "egress",
				"metric":       "rl_drop_bytes",
			}),
			outputs.BytesCounter.MakeReading(qstat.GetRedDropBytes()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "egress",
				"metric":       "red_drop_bytes",
			}),
			output.Number.MakeReading(qstat.GetAllocatedBufferSize()).WithContext(map[string]string{
				"queue_number": queueNumber,
				"direction":    "egress",
				"metric":       "allocated_buffer_size",
			}),
		}
		egressQueueStats = append(egressQueueStats, q...)
	}

	readings = append(readings, ingresQueueStats...)
	readings = append(readings, egressQueueStats...)

	return readings, nil
}
