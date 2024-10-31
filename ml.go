package operadatatypes

import (
	"bytes"
	"encoding/binary"
)

/* ~~ Temperature & Humidity ~~ */

// Data going into ML for calculation of flow temp/hum from raw data
type MlTempHumRawData struct {
	Imx8Temp float32 `json:"imx8_temp" binding:"required"`
	FlowTemp float32 `json:"flow_temp" binding:"required"`
	FlowHum  float32 `json:"flow_hum" binding:"required"`
}

const DATA_TYPE_ML_TEMP_RH = "R"

// Data output by ML for sample temp/hum
type MlTempHumOutputData struct {
	Temp float32 `json:"temp" binding:"required"`
	Hum  float32 `json:"hum" binding:"required"`
}

func (d *MlTempHumOutputData) DisplayData() *DisplayFlowConditions {
	return &DisplayFlowConditions{
		FlowTemp: d.Temp,
		FlowHum:  d.Hum,
	}
}
func (d *MlTempHumRawData) Populate(h *HousekeepingData, s *SecondaryData) {
	d.Imx8Temp = h.PortentaImx8Temp
	d.FlowTemp = s.FlowTemperature
	d.FlowHum = s.FlowHumidity
}

/* ~~ PM2.5 ~~ */

// For raw data we will just take all of the teensy data

// For output, output just a number for PM2.5
const DATA_TYPE_ML_PRIMARY = "P"

type MlPm25OutputData struct {
	UnixSec uint32
	Pm2p5   float32
}

type mlPm25InputDataPulses struct {
	Laser, Pd0, Pd1 uint8
	Pulses          []Pulse
}

func (d *mlPm25InputDataPulses) Serialize() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, d.Laser)
	binary.Write(buf, binary.LittleEndian, d.Pd0)
	binary.Write(buf, binary.LittleEndian, d.Pd1)
	binary.Write(buf, binary.LittleEndian, uint32(len(d.Pulses)))
	for _, p := range d.Pulses {
		binary.Write(buf, binary.LittleEndian, p.Height)
		binary.Write(buf, binary.LittleEndian, p.Width)
		binary.Write(buf, binary.LittleEndian, p.SidePeak)
	}
	return buf.Bytes()
}

type MlPm25InputData struct {
	UnixSec   uint32
	PulseData []mlPm25InputDataPulses
}

func (d *MlPm25InputData) Serialize() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, d.UnixSec)
	binary.Write(buf, binary.LittleEndian, uint32(len(d.PulseData)))
	for _, p := range d.PulseData {
		buf.Write(p.Serialize())
	}
	return buf.Bytes()
}

func (d *MlPm25InputData) Populate(t TeensyData) {
	d.UnixSec = t.UnixSec
	d.PulseData = make([]mlPm25InputDataPulses, len(t.Counts))
	for idx, c := range t.Counts {
		d.PulseData[idx].Laser = c.PinLaser
		d.PulseData[idx].Pd0 = c.PinPd0
		d.PulseData[idx].Pd1 = c.PinPd1
		d.PulseData[idx].Pulses = make([]Pulse, len(c.Pulses))
		copy(d.PulseData[idx].Pulses, c.Pulses)
	}
}

func TeensyDataToMlPm25(t *TeensyData) *MlPm25InputData {
	ret := &MlPm25InputData{}
	ret.Populate(*t)
	return ret
}

func (d *MlPm25OutputData) DisplayData() *DisplayPrimary {
	return &DisplayPrimary{
		Pm2p5:   d.Pm2p5,
		Aerosol: "nil",
	}
}

/* GOBS */

func (d *MlTempHumOutputData) SendGob(unixSocketPath string) error {
	return sendStructGob(d, DATA_TYPE_ML_TEMP_RH, unixSocketPath)
}

func (d *MlPm25OutputData) SendGob(unixSocketPath string) error {
	return sendStructGob(d, DATA_TYPE_ML_PRIMARY, unixSocketPath)
}
