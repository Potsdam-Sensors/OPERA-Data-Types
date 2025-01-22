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
func (d *MlTempHumRawData) Populate(s *SecondaryData) {
	d.Imx8Temp = s.PortentaImx8Temp
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
func TranslateToOldCounts(c *NewTeensyCounts) *TeensyCounts {
	ret := &TeensyCounts{
		PinLaser: c.PinLaser,
		PinPd0:   c.PinPd0,
		PinPd1:   c.PinPd1,

		RawScalar0:    c.RawScalar0,
		RawScalar1:    c.RawScalar1,
		DiffedScalar0: c.DiffedScalar0,
		DiffedScalar1: c.DiffedScalar1,

		Baseline0: c.Baseline0,
		Baseline1: c.Baseline1,

		RawUpperTh0:    c.RawUpperTh0,
		RawUpperTh1:    c.RawUpperTh1,
		DiffedUpperTh0: c.DiffedUpperTh0,
		DiffedUpperTh1: c.DiffedUpperTh1,

		MsRead:      c.MsRead,
		BuffersRead: c.BuffersRead,
		NumPulses:   c.NumPulses,
		MaxLaserOn:  c.MaxLaserOn,

		PulsesPerSecond: c.PulsesPerSecond,

		Pulses: make([]Pulse, len(c.Pulses)),
	}
	usPerPoint := float32(c.MsRead*1000) / float32(c.BuffersRead*3500)
	for idx, p := range c.Pulses {
		ret.Pulses[idx] = Pulse{
			Height:   float32(p.RawPeak) - c.Baseline0,
			Width:    float32(p.Indices[2]+p.Indices[5]) * usPerPoint,
			SidePeak: float32(p.SidePeak) - c.Baseline1,
		}
	}

	return ret
}

func NewTeensyDataToMlPm25(t *NewTeensyData) *MlPm25InputData {
	faux := TeensyData{
		UnixSec: t.UnixSec,
		Counts:  make([]*TeensyCounts, len(t.Counts)),
	}
	for idx, c := range t.Counts {
		faux.Counts[idx] = TranslateToOldCounts(c)
	}
	ret := &MlPm25InputData{}
	ret.Populate(faux)
	return ret
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

type MlClassificationOutputData struct {
	UnixSec       uint32    `json:"unix_sec" binding:"required"`
	Labels        []string  `json:"labels" binding:"required"`
	Probabilities []float32 `json:"probabilities" binding:"required"`
}

func (d *MlClassificationOutputData) GetClass() string {
	maxProbIdx := 0
	for idx, p := range d.Probabilities {
		if p > d.Probabilities[maxProbIdx] {
			maxProbIdx = idx
		}
	}
	if d.Probabilities[maxProbIdx] < CONFIDENCE_INTERVAL_MIN {
		return AEROSOL_NAME_UNKOWN
	}
	return d.Labels[maxProbIdx]
}

type MlPrimaryDataOutput struct {
	UnixSec       uint32
	Classifcation MlClassificationOutputData
	Pm25          MlPm25OutputData
}

/* GOBS */

func (d *MlTempHumOutputData) SendGob(unixSocketPath string) error {
	return sendStructGob(d, DATA_TYPE_ML_TEMP_RH, unixSocketPath)
}

func (d *MlPrimaryDataOutput) SendGob(unixSocketPath string) error {
	return sendStructGob(d, DATA_TYPE_ML_PRIMARY, unixSocketPath)
}

func (d *MlPrimaryDataOutput) DisplayData() *DisplayPrimary {
	return &DisplayPrimary{
		Pm2p5:   d.Pm25.Pm2p5,
		Aerosol: d.Classifcation.GetClass(),
	}
}
