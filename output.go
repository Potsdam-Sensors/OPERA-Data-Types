package operadatatypes

import (
	"fmt"
	"time"
)

/* Util */
func generateFileName(portentaSerial, dataLabel string, timestamp uint32) string {
	t := time.Unix(int64(timestamp), 0)
	return fmt.Sprintf("OPERA_%s_%s_%04d%02d%02d.csv", portentaSerial, dataLabel, t.Year(), t.Month(), t.Day())
}

/* Structs */
type HousekeepingData struct {
	/* Primary Keys */
	Unix           uint32
	PortentaSerial string

	FlowRate            float32    // Flow rate on OPC board
	PortentaImx8Temp    float32    // Portenta's main CPU temp
	TeensyMcuTemp       float32    // Teensy MCU temp
	OpticalTemperatures [3]float32 // Laser temperatures from OPC
	OmbTemperatureHtu   float32    // Temp & hum from HTU on OMB
	OmbHumidityHtu      float32    //
	OmbTemperatureScd   float32    // Temp & hum from SCD41 on OMB
	OmbHumidityScd      float32    //
	Monitor5vMean       float32
	Monitor5vStdDev     float32
}

type SecondaryData struct {
	/* Primary Keys */
	Unix           uint32
	PortentaSerial string

	Sps30Pm2p5      float32
	Pressure        float32
	Co2             uint32
	VocIndex        int32
	FlowTemperature float32
	FlowHumidity    float32
	FlowRate        float32
}

type OutputData interface {
	CsvFileWriteJob(string) []CsvFileWriteJob
}

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

	d.Sps30Pm2p5 = sps30.Pm2p5
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

/* To Csv Functions */

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
		Headers:  "unix,portenta,sps30_pm2p5,pressure,co2,voc_index,flow_temp,flow_hum,flow_rate",
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
