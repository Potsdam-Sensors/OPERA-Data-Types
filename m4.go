package operadatatypes

import (
	"fmt"
)

const DATA_TYPE_M4_SENSORS = "M"

type M4SensorMeasurement struct {
	UnixSec         uint32  `json:"unix_sec" binding:"required"`
	Pressure        float32 `json:"pressure" binding:"required"`
	TempHtu         float32 `json:"temp_htu" binding:"required"`
	TempScd         float32 `json:"temp_scd" binding:"required"`
	HumHtu          float32 `json:"hum_htu" binding:"required"`
	HumScd          float32 `json:"hum_scnd" binding:"required"`
	Co2             uint32  `json:"co2" binding:"required"`
	VocIndex        int32   `json:"voc" binding:"required"`
	OpticalTemp0    float32 `json:"optical_temp0" binding:"required"`
	OpticalTemp1    float32 `json:"optical_temp1" binding:"required"`
	OpticalTemp2    float32 `json:"optical_temp2" binding:"required"`
	Monitor5VMean   float32 `json:"monitor5vMean" binding:"required"`
	Monitor5VStdDev float32 `json:"monitor5vStdDev" binding:"required"`
}

func (m *M4SensorMeasurement) String() string {
	return fmt.Sprintf("[M4| Temp %.2f, %.2f | Hum %.2f, %.2f | CO2 %d | Pressure %.1f | VOC Index %d | Optical Temps %.2f, %.2f & %.2f, 5v: Mean %.4f, Std. %.4f]",
		m.TempHtu, m.TempScd, m.HumHtu, m.HumScd, m.Co2, m.Pressure, m.VocIndex, m.OpticalTemp0, m.OpticalTemp1, m.OpticalTemp2, m.Monitor5VMean, m.Monitor5VStdDev)
}

func (d *M4SensorMeasurement) SendGob(unixSocketPath string) error {
	return sendStructGob(d, DATA_TYPE_M4_SENSORS, unixSocketPath)
}
