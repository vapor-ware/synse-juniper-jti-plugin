package jti

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/protocol/jti/protos/optics"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg/protocol/jti/protos/telemetry_top"
)

func TestNewOpticsContextFromStream(t *testing.T) {
	ctx := NewOpticsContextFromStream(&telemetry_top.TelemetryStream{
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

func TestOpticsContext_Decode(t *testing.T) {
	ctx := OpticsContext{
		SensorName:     "sensor",
		SystemID:       "test",
		ComponentID:    2,
		SubComponentID: 0,
	}
	o := &optics.Optics{
		OpticsDiag: []*optics.OpticsInfos{{
			IfName: &stringVal,
			OpticsDiagStats: &optics.OpticsDiagStats{
				OpticsType: &uint32Val,
				ModuleTemp: &float64Val,
			},
		}},
	}

	data, err := ctx.Decode(o)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Len(t, data, 1)
}

func TestOpticsContext_Decode_NilOptics(t *testing.T) {
	ctx := OpticsContext{
		SensorName:     "sensor",
		SystemID:       "test",
		ComponentID:    2,
		SubComponentID: 0,
	}

	data, err := ctx.Decode(nil)
	assert.NoError(t, err)
	assert.Len(t, data, 0)
}

func TestOpticsContext_Decode_ErrMakeDeviceInfo(t *testing.T) {
	ctx := OpticsContext{
		SensorName:     "sensor",
		SystemID:       "test",
		ComponentID:    2,
		SubComponentID: 0,
	}

	data, err := ctx.Decode(&optics.Optics{
		OpticsDiag: []*optics.OpticsInfos{{
			//  No interface name specified
			OpticsDiagStats: &optics.OpticsDiagStats{
				OpticsType: &uint32Val,
				ModuleTemp: &float64Val,
			},
		}},
	})
	assert.Error(t, err)
	assert.Nil(t, data)
}

func TestOpticsContext_MakeDeviceInfo(t *testing.T) {
	ctx := OpticsContext{
		SensorName:     "sensor",
		SystemID:       "test",
		ComponentID:    2,
		SubComponentID: 0,
	}
	infos := &optics.OpticsInfos{
		IfName: &stringVal,
		OpticsDiagStats: &optics.OpticsDiagStats{
			OpticsType: &uint32Val,
			ModuleTemp: &float64Val,
		},
	}

	info, err := ctx.MakeDeviceInfo(infos)
	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, "interface", info.Type)
	assert.Equal(t, "test interface string", info.Info)
	assert.Equal(t, []string{"vapor/networking:interface"}, info.Tags)
	assert.Equal(t, map[string]string{
		"component_id": "2",
		"name":         "string",
		"system_id":    "test",
		"sensor_name":  "sensor",
	}, info.Context)
	assert.Equal(t, map[string]string{
		"sys":  "test",
		"if":   "string",
		"cid":  "2",
		"scid": "0",
	}, info.IDComponents)
}

func TestOpticsContext_MakeDeviceInfo_ErrNilInfo(t *testing.T) {
	ctx := OpticsContext{
		SensorName:     "sensor",
		SystemID:       "test",
		ComponentID:    2,
		SubComponentID: 0,
	}

	info, err := ctx.MakeDeviceInfo(nil)
	assert.Error(t, err)
	assert.Nil(t, info)
}

func TestOpticsContext_MakeDeviceInfo_ErrNoSystemID(t *testing.T) {
	ctx := OpticsContext{
		SensorName:     "sensor",
		SystemID:       "",
		ComponentID:    2,
		SubComponentID: 0,
	}
	infos := &optics.OpticsInfos{
		IfName: &stringVal,
		OpticsDiagStats: &optics.OpticsDiagStats{
			OpticsType: &uint32Val,
			ModuleTemp: &float64Val,
		},
	}

	info, err := ctx.MakeDeviceInfo(infos)
	assert.Error(t, err)
	assert.Nil(t, info)
}

func TestOpticsContext_MakeDeviceInfo_ErrNoIfaceName(t *testing.T) {
	ctx := OpticsContext{
		SensorName:     "sensor",
		SystemID:       "test",
		ComponentID:    2,
		SubComponentID: 0,
	}
	infos := &optics.OpticsInfos{
		OpticsDiagStats: &optics.OpticsDiagStats{
			OpticsType: &uint32Val,
			ModuleTemp: &float64Val,
		},
	}

	info, err := ctx.MakeDeviceInfo(infos)
	assert.Error(t, err)
	assert.Nil(t, info)
}

func TestOpticsContext_MakeReadings(t *testing.T) {
	ctx := OpticsContext{
		SensorName:     "sensor",
		SystemID:       "test",
		ComponentID:    2,
		SubComponentID: 0,
	}
	infos := &optics.OpticsInfos{
		IfName: &stringVal,
		OpticsDiagStats: &optics.OpticsDiagStats{
			OpticsType: &uint32Val,
			ModuleTemp: &float64Val,
		},
	}

	readings, err := ctx.MakeReadings(infos)
	assert.NoError(t, err)
	assert.Len(t, readings, 22)
}

func TestOpticsContext_MakeReadings2(t *testing.T) {
	ctx := OpticsContext{
		SensorName:     "sensor",
		SystemID:       "test",
		ComponentID:    2,
		SubComponentID: 0,
	}
	infos := &optics.OpticsInfos{
		IfName: &stringVal,
		OpticsDiagStats: &optics.OpticsDiagStats{
			OpticsType: &uint32Val,
			ModuleTemp: &float64Val,
			OpticsLaneDiagStats: []*optics.OpticsDiagLaneStats{{
				LaneNumber: &uint32Val,
			}},
		},
	}

	readings, err := ctx.MakeReadings(infos)
	assert.NoError(t, err)
	assert.Len(t, readings, 41)
}
