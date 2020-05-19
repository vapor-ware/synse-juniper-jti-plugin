package jti

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/protocol/jti/protos/port"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/protocol/jti/protos/telemetry_top"
)

var (
	stringVal  = "string"
	uint32Val  = uint32(1)
	uint64Val  = uint64(1)
	float64Val = float64(1)
)

func TestNewPortContextFromStream(t *testing.T) {
	ctx := NewPortContextFromStream(&telemetry_top.TelemetryStream{
		SensorName:     &stringVal,
		SystemId:       &stringVal,
		ComponentId:    &uint32Val,
		SubComponentId: &uint32Val,
	})
	assert.NotNil(t, ctx)
	assert.Equal(t, stringVal, ctx.SensorName)
	assert.Equal(t, stringVal, ctx.SystemID)
	assert.Equal(t, uint32Val, ctx.ComponentID)
	assert.Equal(t, uint32Val, ctx.SubComponentID)
}

func TestPortContext_Decode(t *testing.T) {
	ctx := PortContext{
		SensorName:     "sensor",
		SystemID:       "test",
		ComponentID:    2,
		SubComponentID: 0,
	}
	p := &port.Port{
		InterfaceStats: []*port.InterfaceInfos{{
			IfName:                 &stringVal,
			InitTime:               &uint64Val,
			SnmpIfIndex:            &uint32Val,
			ParentAeName:           &stringVal,
			IfAdministrationStatus: &stringVal,
			IfOperationalStatus:    &stringVal,
			IfDescription:          &stringVal,
			IfTransitions:          &uint64Val,
			IfLastChange:           &uint32Val,
			IfHighSpeed:            &uint32Val,
		}},
	}

	data, err := ctx.Decode(p)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Len(t, data, 1)
}

func TestPortContext_Decode_NilPort(t *testing.T) {
	ctx := PortContext{
		SensorName:     "sensor",
		SystemID:       "test",
		ComponentID:    2,
		SubComponentID: 0,
	}

	data, err := ctx.Decode(nil)
	assert.NoError(t, err)
	assert.Len(t, data, 0)
}

func TestPortContext_Decode_ErrMakeDeviceInfo(t *testing.T) {
	ctx := PortContext{
		SensorName:     "sensor",
		SystemID:       "test",
		ComponentID:    2,
		SubComponentID: 0,
	}

	data, err := ctx.Decode(&port.Port{
		InterfaceStats: []*port.InterfaceInfos{{
			// No interface name
			InitTime:               &uint64Val,
			SnmpIfIndex:            &uint32Val,
			ParentAeName:           &stringVal,
			IfAdministrationStatus: &stringVal,
			IfOperationalStatus:    &stringVal,
			IfDescription:          &stringVal,
			IfTransitions:          &uint64Val,
			IfLastChange:           &uint32Val,
			IfHighSpeed:            &uint32Val,
		}},
	})
	assert.Error(t, err)
	assert.Nil(t, data)
}

func TestPortContext_MakeDeviceInfo(t *testing.T) {
	ctx := PortContext{
		SensorName:     "sensor",
		SystemID:       "test",
		ComponentID:    2,
		SubComponentID: 0,
	}
	infos := &port.InterfaceInfos{
		IfName:                 &stringVal,
		InitTime:               &uint64Val,
		SnmpIfIndex:            &uint32Val,
		ParentAeName:           &stringVal,
		IfAdministrationStatus: &stringVal,
		IfOperationalStatus:    &stringVal,
		IfDescription:          &stringVal,
		IfTransitions:          &uint64Val,
		IfLastChange:           &uint32Val,
		IfHighSpeed:            &uint32Val,
	}

	info, err := ctx.MakeDeviceInfo(infos)
	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, "interface", info.Type)
	assert.Equal(t, "test interface string", info.Info)
	assert.Equal(t, []string{"vapor/networking:interface"}, info.Tags)
	assert.Equal(t, map[string]string{
		"component_id":   "2",
		"name":           "string",
		"system_id":      "test",
		"sensor_name":    "sensor",
		"parent_ae_name": "string",
	}, info.Context)
	assert.Equal(t, map[string]string{
		"sys":  "test",
		"if":   "string",
		"cid":  "2",
		"scid": "0",
	}, info.IDComponents)
}

func TestPortContext_MakeDeviceInfo_ErrNilInfos(t *testing.T) {
	ctx := PortContext{
		SensorName:     "sensor",
		SystemID:       "test",
		ComponentID:    2,
		SubComponentID: 0,
	}

	info, err := ctx.MakeDeviceInfo(nil)
	assert.Error(t, err)
	assert.Nil(t, info)
}

func TestPortContext_MakeDeviceInfo_ErrNoSystemID(t *testing.T) {
	ctx := PortContext{
		SensorName:     "sensor",
		SystemID:       "",
		ComponentID:    2,
		SubComponentID: 0,
	}
	infos := &port.InterfaceInfos{
		IfName:                 &stringVal,
		InitTime:               &uint64Val,
		SnmpIfIndex:            &uint32Val,
		ParentAeName:           &stringVal,
		IfAdministrationStatus: &stringVal,
		IfOperationalStatus:    &stringVal,
		IfDescription:          &stringVal,
		IfTransitions:          &uint64Val,
		IfLastChange:           &uint32Val,
		IfHighSpeed:            &uint32Val,
	}

	info, err := ctx.MakeDeviceInfo(infos)
	assert.Error(t, err)
	assert.Nil(t, info)
}

func TestPortContext_MakeDeviceInfo_ErrNoIfaceName(t *testing.T) {
	ctx := PortContext{
		SensorName:     "sensor",
		SystemID:       "test",
		ComponentID:    2,
		SubComponentID: 0,
	}
	infos := &port.InterfaceInfos{
		InitTime:               &uint64Val,
		SnmpIfIndex:            &uint32Val,
		ParentAeName:           &stringVal,
		IfAdministrationStatus: &stringVal,
		IfOperationalStatus:    &stringVal,
		IfDescription:          &stringVal,
		IfTransitions:          &uint64Val,
		IfLastChange:           &uint32Val,
		IfHighSpeed:            &uint32Val,
	}

	info, err := ctx.MakeDeviceInfo(infos)
	assert.Error(t, err)
	assert.Nil(t, info)
}

func TestPortContext_MakeReadings(t *testing.T) {
	ctx := PortContext{
		SensorName:     "sensor",
		SystemID:       "test",
		ComponentID:    2,
		SubComponentID: 0,
	}
	infos := &port.InterfaceInfos{
		IfName:                 &stringVal,
		InitTime:               &uint64Val,
		SnmpIfIndex:            &uint32Val,
		ParentAeName:           &stringVal,
		IfAdministrationStatus: &stringVal,
		IfOperationalStatus:    &stringVal,
		IfDescription:          &stringVal,
		IfTransitions:          &uint64Val,
		IfLastChange:           &uint32Val,
		IfHighSpeed:            &uint32Val,
	}

	readings, err := ctx.MakeReadings(infos)
	assert.NoError(t, err)

	assert.Len(t, readings, 39)
}

func TestPortContext_MakeReadings2(t *testing.T) {
	ctx := PortContext{
		SensorName:     "sensor",
		SystemID:       "test",
		ComponentID:    2,
		SubComponentID: 0,
	}
	infos := &port.InterfaceInfos{
		IfName:                 &stringVal,
		InitTime:               &uint64Val,
		SnmpIfIndex:            &uint32Val,
		ParentAeName:           &stringVal,
		IfAdministrationStatus: &stringVal,
		IfOperationalStatus:    &stringVal,
		IfDescription:          &stringVal,
		IfTransitions:          &uint64Val,
		IfLastChange:           &uint32Val,
		IfHighSpeed:            &uint32Val,
		EgressQueueInfo: []*port.QueueStats{
			{},
		},
	}

	readings, err := ctx.MakeReadings(infos)
	assert.NoError(t, err)

	assert.Len(t, readings, 50)
}

func TestPortContext_MakeReadings3(t *testing.T) {
	ctx := PortContext{
		SensorName:     "sensor",
		SystemID:       "test",
		ComponentID:    2,
		SubComponentID: 0,
	}
	infos := &port.InterfaceInfos{
		IfName:                 &stringVal,
		InitTime:               &uint64Val,
		SnmpIfIndex:            &uint32Val,
		ParentAeName:           &stringVal,
		IfAdministrationStatus: &stringVal,
		IfOperationalStatus:    &stringVal,
		IfDescription:          &stringVal,
		IfTransitions:          &uint64Val,
		IfLastChange:           &uint32Val,
		IfHighSpeed:            &uint32Val,
		IngressQueueInfo: []*port.QueueStats{
			{},
		},
	}

	readings, err := ctx.MakeReadings(infos)
	assert.NoError(t, err)

	assert.Len(t, readings, 50)
}

func TestPortContext_MakeReadings4(t *testing.T) {
	ctx := PortContext{
		SensorName:     "sensor",
		SystemID:       "test",
		ComponentID:    2,
		SubComponentID: 0,
	}
	infos := &port.InterfaceInfos{
		IfName:                 &stringVal,
		InitTime:               &uint64Val,
		SnmpIfIndex:            &uint32Val,
		ParentAeName:           &stringVal,
		IfAdministrationStatus: &stringVal,
		IfOperationalStatus:    &stringVal,
		IfDescription:          &stringVal,
		IfTransitions:          &uint64Val,
		IfLastChange:           &uint32Val,
		IfHighSpeed:            &uint32Val,
		EgressQueueInfo: []*port.QueueStats{
			{},
		},
		IngressQueueInfo: []*port.QueueStats{
			{},
		},
	}

	readings, err := ctx.MakeReadings(infos)
	assert.NoError(t, err)

	assert.Len(t, readings, 61)
}
