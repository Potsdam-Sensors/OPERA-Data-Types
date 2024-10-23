package operadatatypes

import (
	"fmt"
)

const DATA_TYPE_TEENSY = "T"

type TeensyData struct {
	UnixSec   uint32  `json:"unix_sec" binding:"required"`
	McuTemp   float32 `json:"mcu_temp" binding:"required"`
	FlowTemp  float32 `json:"flow_temp" binding:"required"`
	FlowHum   float32 `json:"flow_hum" binding:"required"`
	FlowRate  float32 `json:"flow_rate" binding:"required"`
	HvEnabled bool    `json:"hv_enabled" binding:"required"`
	HvSet     uint8   `json:"hv_set" binding:"required"`
	HvMonitor uint16  `json:"hv_monitor" binding:"required"`

	Counts []*TeensyCounts `json:"counts" binding:"required"`
}

type Pulse struct { // 12
	Height   float32 `json:"height" binding:"required"`
	Width    float32 `json:"width" binding:"required"`
	SidePeak float32 `json:"side_peak" binding:"required"`
}

type TeensyCounts struct { // 55
	PinPd0   uint8 `json:"pin_pd0" binding:"required"`
	PinPd1   uint8 `json:"pin_pd1" binding:"required"`
	PinLaser uint8 `json:"pin_laser" binding:"required"`
	//3

	RawScalar0    float32 `json:"raw_scalar0" binding:"required"`
	RawScalar1    float32 `json:"raw_scalar1" binding:"required"`
	DiffedScalar0 float32 `json:"diffed_scalar0" binding:"required"`
	DiffedScalar1 float32 `json:"diffed_scalar1" binding:"required"`
	//16

	Baseline0 float32 `json:"baseline0" binding:"required"`
	Baseline1 float32 `json:"baseline1" binding:"required"`
	// 8

	RawUpperTh0    float32 `json:"raw_upper_th0" binding:"required"`
	RawUpperTh1    float32 `json:"raw_upper_th1" binding:"required"`
	DiffedUpperTh0 float32 `json:"diffed_upper_th0" binding:"required"`
	DiffedUpperTh1 float32 `json:"diffed_upper_th1" binding:"required"`
	//16

	MsRead      uint32 `json:"ms_read" binding:"required"`
	BuffersRead uint32 `json:"buffers_read" binding:"required"`
	NumPulses   uint32 `json:"num_pulses" binding:"required"`
	// 12
	PulsesPerSecond float32 `json:"pulses_per_second" binding:"required"`
	// x

	Pulses []Pulse `json:"pulses" binding:"required"`
	// 12 * NumPulses
}

func (c *TeensyCounts) String() string {
	return fmt.Sprintf("[Counts %d,%d:%d | %d ms, %d Buffers, %d Pulses [%.3f pulses/s] | Baselines: %.2f & %.2f]", c.PinPd0, c.PinPd1, c.PinLaser, c.MsRead, c.BuffersRead, c.NumPulses, c.PulsesPerSecond, c.Baseline0, c.Baseline1)
}
func (d *TeensyData) String() string {
	return fmt.Sprintf("[Teensy Data | Unix %d | MCU Temp %.1f degC | Flow %.1f degC, %.1f perc., %.4f m/s | Hv Enabled: %v, Set: %d, Val: %d]", d.UnixSec, d.McuTemp, d.FlowTemp, d.FlowHum, d.FlowRate,
		d.HvEnabled, d.HvSet, d.HvMonitor)
}
func (d *TeensyData) SendGob(unixSocketPath string) error {
	return sendStructGob(d, DATA_TYPE_TEENSY, unixSocketPath)
}
