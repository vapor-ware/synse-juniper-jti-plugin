package outputs

import (
	"github.com/vapor-ware/synse-sdk/sdk/output"
)

// Boolean is a reading output which describes a true/false value.
var Boolean = output.Output{
	Name: "boolean",
	Type: "bool",
}

// BytesCounter is a reading output which describes a count of bytes. This is
// not associated with a time, so it is not a rate -- merely just a total count
// of bytes.
var BytesCounter = output.Output{
	Name: "bytes",
	Type: "counter",
	Unit: &output.Unit{
		Name:   "bytes",
		Symbol: "b",
	},
}

// BytesPerSecond is a reading output which describes a rate of bytes over a second.
var BytesPerSecond = output.Output{
	Name: "bytes-per-second",
	Type: "throughput",
	Unit: &output.Unit{
		Name:   "bytes per second",
		Symbol: "bytes/s",
	},
}

// DecibelMilliwatts is a reading output which describes a measure of absolute power expressed
// as a ratio between decibels to one milliwatt.
var DecibelMilliwatts = output.Output{
	Name: "decibel-milliwatt",
	Type: "power",
	Unit: &output.Unit{
		Name:   "decibel-milliwatt",
		Symbol: "dBm",
	},
}

// MegabitPerSecond is a reading output which describes a rate of 1,000,000 bits over a second.
var MegabitPerSecond = output.Output{
	Name: "megabit-per-second",
	Type: "throughput",
	Unit: &output.Unit{
		Name:   "Megabits per second",
		Symbol: "Mbit/s",
	},
}

// Milliamperes is a reading output which describes an electrical current as measured in thousandths
// of an Ampere (milli-amperes).
var Milliamperes = output.Output{
	Name: "milliampere",
	Type: "current",
	Unit: &output.Unit{
		Name:   "milliamperes",
		Symbol: "mA",
	},
}

// PacketsCounter is a reading output which describes a count of packets.
// This is not associated with a time, so is not a rate -- merely just a
// total count of packets.
var PacketsCounter = output.Output{
	Name: "packets",
	Type: "counter",
	Unit: &output.Unit{
		Name:   "packets",
		Symbol: "pkts",
	},
}

// PacketsPerSecond is a reading output which describes a rate of packets over a second.
var PacketsPerSecond = output.Output{
	Name: "packets-per-second",
	Type: "throughput",
	Unit: &output.Unit{
		Name:   "packets per second",
		Symbol: "pkts/s",
	},
}

// TimeTicks is a reading output which describes the passage of time, as measured
// in "time ticks".
var TimeTicks = output.Output{
	Name: "time-ticks",
	Type: "time",
	Unit: &output.Unit{
		Name:   "ticks",
		Symbol: "ticks",
	},
}
