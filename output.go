package operadatatypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

/* Util */
const CSV_FILE_EXTENSION = ".csv"
const BINARY_FILE_EXTENSION = ".raw"

func generateFileName(portentaSerial, dataLabel string, timestamp uint32, isCsv bool) string {
	t := time.Unix(int64(timestamp), 0)
	ret := fmt.Sprintf("OPERA_%s_%s_%04d%02d%02d", portentaSerial, dataLabel, t.Year(), t.Month(), t.Day())
	if isCsv {
		return ret + CSV_FILE_EXTENSION
	} else {
		return ret + BINARY_FILE_EXTENSION
	}
}

func writeStringToBinary(w io.Writer, s string) {
	binary.Write(w, binary.LittleEndian, uint32(len(s)))
	w.Write([]byte(s))
}
func readStringFromBinary(r io.Reader) (string, error) {
	var n uint32
	if err := binary.Read(r, binary.LittleEndian, &n); err != nil {
		return "", err
	}
	buf := make([]byte, n)
	if _, err := r.Read(buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

/* Structs */
type SecondaryData struct {
	UnixSec        uint32
	PortentaSerial string

	Sps30           Sps30Data
	Pressure        float32
	Co2             uint32
	VocIndex        int32
	FlowTemperature float32
	FlowHumidity    float32
	FlowRate        float32

	PortentaImx8Temp    float32
	TeensyMcuTemp       float32
	OpticalTemperatures [3]float32
	OmbTemperatureHtu   float32
	OmbHumidityHtu      float32
	OmbTemperatureScd   float32
	OmbHumidityScd      float32
	Monitor5vMean       float32
	Monitor5vStdDev     float32
}

func (d *SecondaryData) Populate(portentaSerial string, portentaImx8Temp float32, t *NewTeensyData, s *Sps30Data, m *M4SensorMeasurement) {
	d.UnixSec = t.UnixSec
	d.PortentaSerial = portentaSerial

	d.Sps30 = *s
	d.Pressure = m.Pressure
	d.Co2 = m.Co2
	d.VocIndex = m.VocIndex
	d.FlowTemperature = t.FlowTemp
	d.FlowHumidity = t.FlowHum
	d.FlowRate = t.FlowRate

	d.PortentaImx8Temp = portentaImx8Temp
	d.TeensyMcuTemp = t.McuTemp
	d.OpticalTemperatures = [3]float32{m.OpticalTemp0, m.OpticalTemp1, m.OpticalTemp2}
	d.OmbTemperatureHtu = m.TempHtu
	d.OmbHumidityHtu = m.HumHtu
	d.OmbTemperatureScd = m.TempScd
	d.OmbHumidityScd = m.HumScd
	d.Monitor5vMean = m.Monitor5VMean
	d.Monitor5vStdDev = m.Monitor5VStdDev
}

type PrimaryData = NewTeensyData

type OperaData struct {
	UnixSec        uint32
	PortentaSerial string

	Pm2p5       float32
	ClassLabel  string
	ClassLabels []string
	ClassProbs  []float32

	Temp       float32
	RH         float32
	Sps30Pm2p5 float32
	Pressure   float32
	Co2        uint32
	VocIndex   int32
}

func (d *OperaData) Populate(portentaSerial string, m *M4SensorMeasurement, s *Sps30Data, ml *MlPrimaryDataOutput, tr *MlTempHumOutputData) {
	d.UnixSec = ml.UnixSec
	d.PortentaSerial = portentaSerial

	d.Pm2p5 = ml.Pm25.Pm2p5
	d.ClassLabel = ml.Classifcation.GetClass()
	d.ClassLabels = ml.Classifcation.Labels
	d.ClassProbs = ml.Classifcation.Probabilities

	d.Temp = tr.Temp
	d.RH = tr.Hum
	d.Sps30Pm2p5 = s.Pm2p5
	d.Pressure = m.Pressure
	d.Co2 = m.Co2
	d.VocIndex = m.VocIndex
}

type OutputData interface {
	CsvFileWriteJob(string) []CsvFileWriteJob
	Pack(io.Writer)
	Unpack(io.Reader) error
}

/* CSV File Write Job */
func (d *SecondaryData) CsvFileWriteJob(portentaSerial string) []CsvFileWriteJob {
	return []CsvFileWriteJob{{
		Filename: generateFileName(portentaSerial, "SecondaryRaw", d.UnixSec, true),
		Headers:  "unix,portenta,sps30_pm1,sps30_pm2p5,sps30_pm4,sps30_pm10,sps30_pn0p5,sps30_pn1,sps30_pn2p5,sps30_pn4,sps30_pn10,sps30_tps,pressure,co2,voc_index,flow_temp,flow_hum,flow_rate,imx8_temp,teensy_temp,optical_temp0,optical_temp1,optical_temp2,omb_temp_htu,omb_hum_htu,omb_temp_scd,omg_hum_scd,mean_5v_monitor,std_dev_5v_monitor",
		Content: fmt.Sprintf("%d,%s,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%d,%d,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.4f,%.4f",
			d.UnixSec, d.PortentaSerial, d.Sps30.Pm1, d.Sps30.Pm2p5, d.Sps30.Pm4, d.Sps30.Pm10, d.Sps30.Pn0p5, d.Sps30.Pn1, d.Sps30.Pn2p5, d.Sps30.Pn4, d.Sps30.Pn10, d.Sps30.TypicalParticleSize, d.Pressure, d.Co2, d.VocIndex, d.FlowTemperature, d.FlowHumidity, d.FlowRate, d.PortentaImx8Temp, d.TeensyMcuTemp, d.OpticalTemperatures[0], d.OpticalTemperatures[1], d.OpticalTemperatures[2], d.OmbTemperatureHtu, d.OmbHumidityHtu, d.OmbTemperatureScd, d.OmbHumidityScd, d.Monitor5vMean, d.Monitor5vStdDev),
	}}
}

func (d *OperaData) CsvFileWriteJob(portentaSerial string) []CsvFileWriteJob {
	return []CsvFileWriteJob{{
		Filename: generateFileName(portentaSerial, "Output", d.UnixSec, true),
		Headers:  "unix,portenta,pm2p5,class_label,class_labels,class_probs,temp,rh,sps30_pm2p5,pressure,co2,voc_index",
		Content: fmt.Sprintf("%d,%s,%.1f,%s,\"%s\",\"%.1f\",%.1f,%.1f,%.1f,%.1f,%d,%d",
			d.UnixSec, d.PortentaSerial, d.Pm2p5, d.ClassLabel, d.ClassLabels, d.ClassProbs, d.Temp, d.RH, d.Sps30Pm2p5, d.Pressure, d.Co2, d.VocIndex),
	}}
}

func (p NewPulse) String() string {
	indicesStr := fmt.Sprintf("%d", p.Indices[0])
	for _, ind := range p.Indices[1:] {
		indicesStr += fmt.Sprintf(",%d", ind)
	}
	return fmt.Sprintf("(%d,%d,[%s])", p.RawPeak, p.SidePeak, indicesStr)
}

func (d *PrimaryData) CsvFileWriteJob(portentaSerial string) []CsvFileWriteJob {
	ret := []CsvFileWriteJob{}
	filename := generateFileName(portentaSerial, "PrimaryRaw", d.UnixSec, true)
	for _, c := range d.Counts {
		var pulsesStr = ""
		if n := len(c.Pulses); n > 0 {
			pulsesStr += "\"[" + c.Pulses[0].String()
			if n > 1 {
				for _, p := range c.Pulses[1:] {
					pulsesStr += "," + p.String()
				}
			}
			pulsesStr += "]\""
		}

		ret = append(ret, CsvFileWriteJob{
			Filename: filename,
			Headers: "unix,portenta,hv_enabled,hv_set,hv_read,pd0,pd1,laser,raw_scalar0,raw_scalar1,diff_scalar0,diff_scalar1," +
				"baseline0,baseline1,raw_upper_th0,raw_upper_th1,diff_upper_th0,diff_upper_th1,ms_read,num_buff,max_laser_on,num_pulse,pulses_per_second,pulses",
			Content: fmt.Sprintf("%d,%s,%v,%d,%d,%d,%d,%d,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.3f,%.3f,%d,%d,%d,%d,%.2f,%s",
				d.UnixSec, portentaSerial, d.HvEnabled, d.HvSet, d.HvMonitor, c.PinPd0, c.PinPd1, c.PinLaser,
				c.RawScalar0, c.RawScalar1, c.DiffedScalar0, c.DiffedScalar1, c.Baseline0, c.Baseline1, c.RawUpperTh0, c.RawUpperTh1, c.DiffedUpperTh0, c.DiffedUpperTh1,
				c.MsRead, c.BuffersRead, c.MaxLaserOn, c.NumPulses, c.PulsesPerSecond, pulsesStr),
		})
	}
	return ret
}

/* Binary File Write Job */

func (d *SecondaryData) Pack(w io.Writer) {
	binary.Write(w, binary.LittleEndian, d.UnixSec)
	writeStringToBinary(w, d.PortentaSerial)
	d.Sps30.Pack(w)
	binary.Write(w, binary.LittleEndian, d.Pressure)
	binary.Write(w, binary.LittleEndian, d.Co2)
	binary.Write(w, binary.LittleEndian, d.VocIndex)
	binary.Write(w, binary.LittleEndian, d.FlowTemperature)
	binary.Write(w, binary.LittleEndian, d.FlowHumidity)
	binary.Write(w, binary.LittleEndian, d.FlowRate)
	binary.Write(w, binary.LittleEndian, d.PortentaImx8Temp)
	binary.Write(w, binary.LittleEndian, d.TeensyMcuTemp)
	for _, t := range d.OpticalTemperatures {
		binary.Write(w, binary.LittleEndian, t)
	}
	binary.Write(w, binary.LittleEndian, d.OmbTemperatureHtu)
	binary.Write(w, binary.LittleEndian, d.OmbHumidityHtu)
	binary.Write(w, binary.LittleEndian, d.OmbTemperatureScd)
	binary.Write(w, binary.LittleEndian, d.OmbHumidityScd)
	binary.Write(w, binary.LittleEndian, d.Monitor5vMean)
	binary.Write(w, binary.LittleEndian, d.Monitor5vStdDev)
}

func (d *SecondaryData) Unpack(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &d.UnixSec); err != nil {
		return err
	}
	var err error
	d.PortentaSerial, err = readStringFromBinary(r)
	if err != nil {
		return err
	}
	d.Sps30.Unpack(r)
	if err := binary.Read(r, binary.LittleEndian, &d.Pressure); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &d.Co2); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &d.VocIndex); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &d.FlowTemperature); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &d.FlowHumidity); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &d.FlowRate); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &d.PortentaImx8Temp); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &d.TeensyMcuTemp); err != nil {
		return err
	}
	for i := range d.OpticalTemperatures {
		if err := binary.Read(r, binary.LittleEndian, &d.OpticalTemperatures[i]); err != nil {
			return err
		}
	}
	if err := binary.Read(r, binary.LittleEndian, &d.OmbTemperatureHtu); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &d.OmbHumidityHtu); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &d.OmbTemperatureScd); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &d.OmbHumidityScd); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &d.Monitor5vMean); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &d.Monitor5vStdDev); err != nil {
		return err
	}
	return nil
}

// TODO: Can we use generics for this? The only reason I'm even doing a separate function is the filename
func (d *SecondaryData) BinaryFileWriteJob(portentaSerial string) []BinaryFileWriteJob {
	var buf bytes.Buffer
	d.Pack(&buf)
	return []BinaryFileWriteJob{{
		Filename: generateFileName(portentaSerial, "SecondaryRaw", d.UnixSec, false),
		Content:  buf.Bytes(),
	}}
}

func (d *OperaData) Pack(w io.Writer) {
	binary.Write(w, binary.LittleEndian, d.UnixSec)
	writeStringToBinary(w, d.PortentaSerial)
	binary.Write(w, binary.LittleEndian, d.Pm2p5)
	writeStringToBinary(w, d.ClassLabel)
	binary.Write(w, binary.LittleEndian, uint32(len(d.ClassLabels)))
	for _, l := range d.ClassLabels {
		writeStringToBinary(w, l)
	}
	for _, p := range d.ClassProbs {
		binary.Write(w, binary.LittleEndian, p)
	}
	binary.Write(w, binary.LittleEndian, d.Temp)
	binary.Write(w, binary.LittleEndian, d.RH)
	binary.Write(w, binary.LittleEndian, d.Sps30Pm2p5)
	binary.Write(w, binary.LittleEndian, d.Pressure)
	binary.Write(w, binary.LittleEndian, d.Co2)
	binary.Write(w, binary.LittleEndian, d.VocIndex)
}

func (d *OperaData) Unpack(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &d.UnixSec); err != nil {
		return err
	}
	var err error
	d.PortentaSerial, err = readStringFromBinary(r)
	if err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &d.Pm2p5); err != nil {
		return err
	}
	d.ClassLabel, err = readStringFromBinary(r)
	if err != nil {
		return err
	}
	var n uint32
	if err := binary.Read(r, binary.LittleEndian, &n); err != nil {
		return err
	}
	d.ClassLabels = make([]string, n)
	for i := range d.ClassLabels {
		d.ClassLabels[i], err = readStringFromBinary(r)
		if err != nil {
			return err
		}
	}
	d.ClassProbs = make([]float32, n)
	for i := range d.ClassProbs {
		if err := binary.Read(r, binary.LittleEndian, &d.ClassProbs[i]); err != nil {
			return err
		}
	}
	if err := binary.Read(r, binary.LittleEndian, &d.Temp); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &d.RH); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &d.Sps30Pm2p5); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &d.Pressure); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &d.Co2); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &d.VocIndex); err != nil {
		return err
	}
	return nil
}

func (d *OperaData) BinaryFileWriteJob(portentaSerial string) []BinaryFileWriteJob {
	var buf bytes.Buffer
	d.Pack(&buf)
	return []BinaryFileWriteJob{{
		Filename: generateFileName(portentaSerial, "Output", d.UnixSec, false),
		Content:  buf.Bytes(),
	}}
}

func (p NewPulse) Pack(w io.Writer) {
	binary.Write(w, binary.LittleEndian, p.RawPeak)
	binary.Write(w, binary.LittleEndian, p.SidePeak)
	binary.Write(w, binary.LittleEndian, uint32(len(p.Indices)))
	for _, ind := range p.Indices {
		binary.Write(w, binary.LittleEndian, ind)
	}
}

func (p NewPulse) Unpack(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &p.RawPeak); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &p.SidePeak); err != nil {
		return err
	}
	var n uint32
	if err := binary.Read(r, binary.LittleEndian, &n); err != nil {
		return err
	}
	p.Indices = [NUMBER_INDICES_PULSE]uint16{}
	for i := range p.Indices {
		if err := binary.Read(r, binary.LittleEndian, &p.Indices[i]); err != nil {
			return err
		}
	}
	return nil
}

func (d *PrimaryData) Pack(w io.Writer) {
	binary.Write(w, binary.LittleEndian, d.UnixSec)
	binary.Write(w, binary.LittleEndian, uint32(len(d.Counts)))
	for _, c := range d.Counts {
		binary.Write(w, binary.LittleEndian, c.PinPd0)
		binary.Write(w, binary.LittleEndian, c.PinPd1)
		binary.Write(w, binary.LittleEndian, c.PinLaser)

		binary.Write(w, binary.LittleEndian, c.RawScalar0)
		binary.Write(w, binary.LittleEndian, c.RawScalar1)
		binary.Write(w, binary.LittleEndian, c.DiffedScalar0)
		binary.Write(w, binary.LittleEndian, c.DiffedScalar1)

		binary.Write(w, binary.LittleEndian, c.Baseline0)
		binary.Write(w, binary.LittleEndian, c.Baseline1)

		binary.Write(w, binary.LittleEndian, c.RawUpperTh0)
		binary.Write(w, binary.LittleEndian, c.RawUpperTh1)
		binary.Write(w, binary.LittleEndian, c.DiffedUpperTh0)
		binary.Write(w, binary.LittleEndian, c.DiffedUpperTh1)

		binary.Write(w, binary.LittleEndian, c.MsRead)
		binary.Write(w, binary.LittleEndian, c.BuffersRead)
		binary.Write(w, binary.LittleEndian, c.NumPulses)
		binary.Write(w, binary.LittleEndian, c.MaxLaserOn)

		binary.Write(w, binary.LittleEndian, c.PulsesPerSecond)

		binary.Write(w, binary.LittleEndian, uint32(len(c.Pulses)))
		for _, p := range c.Pulses {
			p.Pack(w)
		}
	}
}

func (d *PrimaryData) Unpack(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &d.UnixSec); err != nil {
		return err
	}
	var n uint32
	if err := binary.Read(r, binary.LittleEndian, &n); err != nil {
		return err
	}
	d.Counts = make([]*NewTeensyCounts, n)
	for i := range d.Counts {
		c := &NewTeensyCounts{}
		if err := binary.Read(r, binary.LittleEndian, &c.PinPd0); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &c.PinPd1); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &c.PinLaser); err != nil {
			return err
		}

		if err := binary.Read(r, binary.LittleEndian, &c.RawScalar0); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &c.RawScalar1); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &c.DiffedScalar0); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &c.DiffedScalar1); err != nil {
			return err
		}

		if err := binary.Read(r, binary.LittleEndian, &c.Baseline0); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &c.Baseline1); err != nil {
			return err
		}

		if err := binary.Read(r, binary.LittleEndian, &c.RawUpperTh0); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &c.RawUpperTh1); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &c.DiffedUpperTh0); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &c.DiffedUpperTh1); err != nil {
			return err
		}

		if err := binary.Read(r, binary.LittleEndian, &c.MsRead); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &c.BuffersRead); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &c.NumPulses); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &c.MaxLaserOn); err != nil {
			return err
		}

		if err := binary.Read(r, binary.LittleEndian, &c.PulsesPerSecond); err != nil {
			return err
		}

		var m uint32
		if err := binary.Read(r, binary.LittleEndian, &m); err != nil {
			return err
		}
		c.Pulses = make([]NewPulse, m)
		for j := range c.Pulses {
			if err := c.Pulses[j].Unpack(r); err != nil {
				return err
			}
		}
		d.Counts[i] = c
	}
	return nil
}

func (d *PrimaryData) BinaryFileWriteJob(portentaSerial string) []BinaryFileWriteJob {
	var buf bytes.Buffer
	d.Pack(&buf)
	return []BinaryFileWriteJob{{
		Filename: generateFileName(portentaSerial, "PrimaryRaw", d.UnixSec, false),
		Content:  buf.Bytes(),
	}}
}

// type HousekeepingData struct {
// 	/* Primary Keys */
// 	Unix           uint32
// 	PortentaSerial string

// 	FlowRate            float32    // Flow rate on OPC board
// 	PortentaImx8Temp    float32    // Portenta's main CPU temp
// 	TeensyMcuTemp       float32    // Teensy MCU temp
// 	OpticalTemperatures [3]float32 // Laser temperatures from OPC
// 	OmbTemperatureHtu   float32    // Temp & hum from HTU on OMB
// 	OmbHumidityHtu      float32    //
// 	OmbTemperatureScd   float32    // Temp & hum from SCD41 on OMB
// 	OmbHumidityScd      float32    //
// 	Monitor5vMean       float32
// 	Monitor5vStdDev     float32
// }

// type SecondaryData struct {
// 	/* Primary Keys */
// 	Unix           uint32
// 	PortentaSerial string

// 	Sps30           Sps30Data
// 	Pressure        float32
// 	Co2             uint32
// 	VocIndex        int32
// 	FlowTemperature float32
// 	FlowHumidity    float32
// 	FlowRate        float32
// }

/*

func (d *HousekeepingData) Populate(unix uint32, serial string, teensy *TeensyData, m4 *M4SensorMeasurement,
	portentaImx8Temp float32) {
	d.Unix = unix
	d.PortentaSerial = serial

	d.FlowRate = teensy.FlowRate
	d.PortentaImx8Temp = portentaImx8Temp
	d.TeensyMcuTemp = teensy.McuTemp
	d.OpticalTemperatures[0] = m4.OpticalTemp0
	d.OpticalTemperatures[1] = m4.OpticalTemp1
	d.OpticalTemperatures[2] = m4.OpticalTemp2
	d.OmbHumidityHtu = m4.HumHtu
	d.OmbHumidityScd = m4.HumScd
	d.OmbTemperatureHtu = m4.TempHtu
	d.OmbTemperatureScd = m4.TempScd
	d.Monitor5vMean = m4.Monitor5VMean
	d.Monitor5vStdDev = m4.Monitor5VStdDev
}

func (d *SecondaryData) Populate(unix uint32, serial string, teensy *TeensyData, m4 *M4SensorMeasurement, sps30 *Sps30Data) {
	d.Unix = unix
	d.PortentaSerial = serial

	d.Sps30 = *sps30
	d.Pressure = m4.Pressure
	d.Co2 = m4.Co2
	d.VocIndex = m4.VocIndex
	d.FlowTemperature = teensy.FlowTemp
	d.FlowHumidity = teensy.FlowHum
	d.FlowRate = teensy.FlowRate
}

type LearnedData struct {
	UnixSec uint32
	Temp    float32
	RH      float32
	Pm2p5   float32

	ClassificationLabels       []string
	ClassifcationProbabilities []float32
}

func (d *LearnedData) Populate(t MlTempHumOutputData, p MlPrimaryDataOutput) {
	d.UnixSec = p.UnixSec
	d.Temp = t.Temp
	d.RH = t.Hum

	d.Pm2p5 = p.Pm25.Pm2p5
	d.ClassificationLabels = p.Classifcation.Labels
	d.ClassifcationProbabilities = p.Classifcation.Probabilities
}
func (d *LearnedData) CsvFileWriteJob(portentaSerial string) []CsvFileWriteJob {
	var strProb string = ""
	if len(d.ClassifcationProbabilities) > 0 {
		strProb = fmt.Sprintf("%.1f", d.ClassifcationProbabilities[0])
		for _, p := range d.ClassifcationProbabilities[1:] {
			strProb += fmt.Sprintf(",%.1f", p)
		}
	}

	var strLabels string = ""
	if len(d.ClassificationLabels) > 0 {
		strLabels = "\"" + d.ClassificationLabels[0] + "\""
		for _, l := range d.ClassificationLabels[1:] {
			strLabels += ",\"" + l + "\""
		}
	}

	return []CsvFileWriteJob{{
		Filename: generateFileName(portentaSerial, "Learned", d.UnixSec),
		Headers:  "unix,portenta,temp,hum,pm2p5,classification_labels,classification_probabilities",
		Content:  fmt.Sprintf("%d,%s,%.1f,%.1f,%.3f,\"%s\",\"%s\"", d.UnixSec, portentaSerial, d.Temp, d.RH, d.Pm2p5, strLabels, strProb),
	}}
}

func (d *SecondaryData) CsvFileWriteJob(portentaSerial string) []CsvFileWriteJob {
	return []CsvFileWriteJob{{
		Filename: generateFileName(portentaSerial, "Secondary", d.Unix),
		Headers:  "unix,portenta,sps30,pressure,co2,voc_index,flow_temp,flow_hum,flow_rate",
		Content: fmt.Sprintf("%d,%s,%.1f,%.1f,%d,%d,%.1f,%.1f,%.1f",
			d.Unix, d.PortentaSerial, d.Sps30Pm2p5, d.Pressure, d.Co2, d.VocIndex, d.FlowTemperature, d.FlowHumidity, d.FlowRate),
	}}

}

func (d *HousekeepingData) CsvFileWriteJob(portentaSerial string) []CsvFileWriteJob {
	return []CsvFileWriteJob{{
		Filename: generateFileName(portentaSerial, "Housekeeping", d.Unix),
		Headers:  "unix,portenta,flow_rate,imx8_temp,teensy_temp,optical_temp0,optical_temp1,optical_temp2,omb_temp_htu,omb_hum_htu,omb_temp_scd,omg_hum_scd,mean_5v_monitor,std_dev_5v_monitor",
		Content: fmt.Sprintf("%d,%s,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.4f,%.4f",
			d.Unix, d.PortentaSerial, d.FlowRate, d.PortentaImx8Temp, d.TeensyMcuTemp, d.OpticalTemperatures[0], d.OpticalTemperatures[1], d.OpticalTemperatures[2], d.OmbTemperatureHtu,
			d.OmbHumidityHtu, d.OmbTemperatureScd, d.OmbHumidityScd, d.Monitor5vMean, d.Monitor5vStdDev),
	}}
}

func (p Pulse) String() string {
	return fmt.Sprintf("(%.1f,%.1f,%.1f)", p.Height, p.SidePeak, p.Width)
}

func (d *TeensyData) CsvFileWriteJob(portentaSerial string) []CsvFileWriteJob {
	ret := []CsvFileWriteJob{}
	filename := generateFileName(portentaSerial, "Raw", d.UnixSec)
	for _, c := range d.Counts {
		var pulsesStr = ""
		if n := len(c.Pulses); n > 0 {
			pulsesStr += "\"" + c.Pulses[0].String()
			if n > 1 {
				for _, p := range c.Pulses[1:] {
					pulsesStr += "," + p.String()
				}
			}
			pulsesStr += "\""
		}

		ret = append(ret, CsvFileWriteJob{
			Filename: filename,
			Headers: "unix,portenta,hv_enabled,hv_set,hv_read,pd0,pd1,laser,raw_scalar0,raw_scalar1,diff_scalar0,diff_scalar1," +
				"baseline0,baseline1,raw_upper_th0,raw_upper_th1,diff_upper_th0,diff_upper_th1,ms_read,num_buff,max_laser_on,num_pulse,pulses_per_second,pulses",
			Content: fmt.Sprintf("%d,%s,%v,%d,%d,%d,%d,%d,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.3f,%.3f,%d,%d,%d,%d,%.2f,%s",
				d.UnixSec, portentaSerial, d.HvEnabled, d.HvSet, d.HvMonitor, c.PinPd0, c.PinPd1, c.PinLaser,
				c.RawScalar0, c.RawScalar1, c.DiffedScalar0, c.DiffedScalar1, c.Baseline0, c.Baseline1, c.RawUpperTh0, c.RawUpperTh1, c.DiffedUpperTh0, c.DiffedUpperTh1,
				c.MsRead, c.BuffersRead, c.MaxLaserOn, c.NumPulses, c.PulsesPerSecond, pulsesStr),
		})
	}
	return ret
}

func (p NewPulse) String() string {
	indicesStr := fmt.Sprintf("%d", p.Indices[0])
	for _, ind := range p.Indices[1:] {
		indicesStr += fmt.Sprintf(",%d", ind)
	}
	return fmt.Sprintf("(%d,%d,[%s])", p.RawPeak, p.SidePeak, indicesStr)
}
func (d *NewTeensyData) CsvFileWriteJob(portentaSerial string) []CsvFileWriteJob {
	ret := []CsvFileWriteJob{}
	filename := generateFileName(portentaSerial, "Raw", d.UnixSec)
	for _, c := range d.Counts {
		var pulsesStr = ""
		if n := len(c.Pulses); n > 0 {
			pulsesStr += "\"[" + c.Pulses[0].String()
			if n > 1 {
				for _, p := range c.Pulses[1:] {
					pulsesStr += "," + p.String()
				}
			}
			pulsesStr += "]\""
		}

		ret = append(ret, CsvFileWriteJob{
			Filename: filename,
			Headers: "unix,portenta,hv_enabled,hv_set,hv_read,pd0,pd1,laser,raw_scalar0,raw_scalar1,diff_scalar0,diff_scalar1," +
				"baseline0,baseline1,raw_upper_th0,raw_upper_th1,diff_upper_th0,diff_upper_th1,ms_read,num_buff,max_laser_on,num_pulse,pulses_per_second,pulses",
			Content: fmt.Sprintf("%d,%s,%v,%d,%d,%d,%d,%d,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.3f,%.3f,%d,%d,%d,%d,%.2f,%s",
				d.UnixSec, portentaSerial, d.HvEnabled, d.HvSet, d.HvMonitor, c.PinPd0, c.PinPd1, c.PinLaser,
				c.RawScalar0, c.RawScalar1, c.DiffedScalar0, c.DiffedScalar1, c.Baseline0, c.Baseline1, c.RawUpperTh0, c.RawUpperTh1, c.DiffedUpperTh0, c.DiffedUpperTh1,
				c.MsRead, c.BuffersRead, c.MaxLaserOn, c.NumPulses, c.PulsesPerSecond, pulsesStr),
		})
	}
	return ret
}

func (d *LearnedData) Pack(b io.Writer) []byte {
	binary.Write(b, binary.LittleEndian, d.UnixSec)
	binary.Write(b, binary.LittleEndian, d.Temp)
	binary.Write(b, binary.LittleEndian, d.RH)
	binary.Write(b, binary.LittleEndian, d.Pm2p5)
	binary.Write(b, binary.LittleEndian, uint32(len(d.ClassificationLabels)))
	for _, l := range d.ClassificationLabels {
		binary.Write(b, binary.LittleEndian, uint32(len(l)))
		b.Write([]byte(l))
	}
	for _, p := range d.ClassifcationProbabilities {
		binary.Write(b, binary.LittleEndian, p)
	}
	return nil
}

func (d *LearnedData) Unpack(b io.Reader) error {
	var n uint32
	if err := binary.Read(b, binary.LittleEndian, &d.UnixSec); err != nil {
		return err
	}
	if err := binary.Read(b, binary.LittleEndian, &d.Temp); err != nil {
		return err
	}
	if err := binary.Read(b, binary.LittleEndian, &d.RH); err != nil {
		return err
	}
	if err := binary.Read(b, binary.LittleEndian, &d.Pm2p5); err != nil {
		return err
	}
	if err := binary.Read(b, binary.LittleEndian, &n); err != nil {
		return err
	}
	d.ClassificationLabels = make([]string, n)
	for i := range d.ClassificationLabels {
		var l uint32
		if err := binary.Read(b, binary.LittleEndian, &l); err != nil {
			return err
		}
		buf := make([]byte, l)
		if _, err := b.Read(buf); err != nil {
			return err
		}
		d.ClassificationLabels[i] = string(buf)
	}
	d.ClassifcationProbabilities = make([]float32, n)
	for i := range d.ClassifcationProbabilities {
		if err := binary.Read(b, binary.LittleEndian, &d.ClassifcationProbabilities[i]); err != nil {
			return err
		}
	}
	return nil
}

func (d *SecondaryData) Pack(b io.Writer) []byte {
	binary.Write(b, binary.LittleEndian, d.Unix)
	binary.Write(b, binary.LittleEndian, uint32(len(d.PortentaSerial)))
	b.Write([]byte(d.PortentaSerial))
	binary.Write(b, binary.LittleEndian, d.Sps30Pm2p5)
	binary.Write(b, binary.LittleEndian, d.Pressure)
	binary.Write(b, binary.LittleEndian, d.Co2)
	binary.Write(b, binary.LittleEndian, d.VocIndex)
	binary.Write(b, binary.LittleEndian, d.FlowTemperature)
	binary.Write(b, binary.LittleEndian, d.FlowHumidity)
	binary.Write(b, binary.LittleEndian, d.FlowRate)
	return nil
}

func (d *SecondaryData) Unpack(b io.Reader) error {
	var n uint32
	if err := binary.Read(b, binary.LittleEndian, &d.Unix); err != nil {
		return err
	}
	if err := binary.Read(b, binary.LittleEndian, &n); err != nil {
		return err
	}
	buf := make([]byte, n)
	if _, err := b.Read(buf); err != nil {
		return err
	}
	d.PortentaSerial = string(buf)
	if err := binary.Read(b, binary.LittleEndian, &d.Sps30Pm2p5); err != nil {
		return err
	}
	if err := binary.Read(b, binary.LittleEndian, &d.Pressure); err != nil {
		return err
	}
	if err := binary.Read(b, binary.LittleEndian, &d.Co2); err != nil {
		return err
	}
	if err := binary.Read(b, binary.LittleEndian, &d.VocIndex); err != nil {
		return err
	}
	if err := binary.Read(b, binary.LittleEndian, &d.FlowTemperature); err != nil {
		return err
	}
	if err := binary.Read(b, binary.LittleEndian, &d.FlowHumidity); err != nil {
		return err
	}
	if err := binary.Read(b, binary.LittleEndian, &d.FlowRate); err != nil {
		return err
	}
	return nil
}

func (d *HousekeepingData) Pack(b io.Writer) []byte {
	binary.Write(b, binary.LittleEndian, d.Unix)
	binary.Write(b, binary.LittleEndian, uint32(len(d.PortentaSerial)))
	b.Write([]byte(d.PortentaSerial))
	binary.Write(b, binary.LittleEndian, d.FlowRate)
	binary.Write(b, binary.LittleEndian, d.PortentaImx8Temp)
	binary.Write(b, binary.LittleEndian, d.TeensyMcuTemp)
	for _, t := range d.OpticalTemperatures {
		binary.Write(b, binary.LittleEndian, t)
	}
	binary.Write(b, binary.LittleEndian, d.OmbTemperatureHtu)
	binary.Write(b, binary.LittleEndian, d.OmbHumidityHtu)
	binary.Write(b, binary.LittleEndian, d.OmbTemperatureScd)
	binary.Write(b, binary.LittleEndian, d.OmbHumidityScd)
	binary.Write(b, binary.LittleEndian, d.Monitor5vMean)
	binary.Write(b, binary.LittleEndian, d.Monitor5vStdDev)
	return nil
}

*/
