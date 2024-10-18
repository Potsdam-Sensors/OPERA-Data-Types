package operadatatypes

import (
	"encoding/json"
	"fmt"
	"net"
)

var DISPLAY_DATA_KEYS = struct {
	PM25           string
	Aerosol        string
	SPS30_PM25     string
	CO2            string
	PRESSURE       string
	VOC            string
	TEMP           string
	RH             string
	IMX8_TEMP_OK   string
	TEENSY_TEMP_OK string
	LASER_TEMP_OK  string
	FLOW_RATE_OK   string
}{
	PM25:           "Primary PM2.5",
	Aerosol:        "Aerosol",
	SPS30_PM25:     "PM2.5",
	CO2:            "CO2",
	PRESSURE:       "pressure",
	VOC:            "VOC Index",
	TEMP:           "Flow Temp",
	RH:             "Flow RH",
	IMX8_TEMP_OK:   "I.MX8 T Ok",
	TEENSY_TEMP_OK: "Teensy T Ok",
	LASER_TEMP_OK:  "Laser T Ok",
	FLOW_RATE_OK:   "Flow Rate Ok",
}

type DisplayPrimary struct {
	Pm2p5   float32 `json:"Primary PM2.5" binding:"required"`
	Aerosol string  `json:"Aerosol" binding:"required"`
}

type DisplayDeviceStatus struct {
	Imx8Temp        bool `json:"I.MX8 T Ok" binding:"required"`
	TeensyMcuTemp   bool `json:"Teensy T Ok" binding:"required"`
	LaserTemps      bool `json:"Laser T Ok" binding:"required"`
	OpticalFlowRate bool `json:"Flow Rate Ok"`
}

type DisplayM4Sensors struct {
	Co2      uint32  `json:"CO2" binding:"required"`
	VocIndex int32   `json:"VOC Index" binding:"required"`
	Pressure float32 `json:"pressure" binding:"required"`
}

type DisplaySps30 struct {
	Pm2p5 float32 `json:"PM2.5" binding:"required"`
}

type DisplayFlowConditions struct {
	FlowTemp float32 `json:"Flow Temp" binding:"required"`
	FlowHum  float32 `json:"Flow RH" binding:"required"`
}

/* Per Type */
func (m *M4SensorMeasurement) DisplayData() any {
	return &DisplayM4Sensors{
		Co2:      m.Co2,
		VocIndex: m.VocIndex,
		Pressure: m.Pressure,
	}
}
func (s *Sps30Data) DisplayData() any {
	return &DisplaySps30{
		Pm2p5: s.Pm2p5,
	}
}

func (t *TeensyData) DisplayData() any {
	return &DisplayFlowConditions{
		FlowTemp: t.FlowTemp,
		FlowHum:  t.FlowHum,
	}
}

func GetDisplayData(m4 *M4SensorMeasurement, teensy *TeensyData, imx8Temp float32) any {
	laserTempOk := true
	for _, tmp := range []float32{m4.OpticalTemp0, m4.OpticalTemp1, m4.OpticalTemp2} {
		if tmp > OK_LASER_TEMP_MAX_C {
			laserTempOk = false
			break
		}
	}

	return &DisplayDeviceStatus{
		Imx8Temp:        imx8Temp <= OK_IMX8_TEMP_MAX_C,
		TeensyMcuTemp:   teensy.McuTemp <= OK_TEENSY_TEMP_MAX_C,
		LaserTemps:      laserTempOk,
		OpticalFlowRate: teensy.FlowRate >= OK_FLOW_RATE_MIN,
	}
}

func SendDisplayData(data any, unixSocketPath string) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to convert to json: %v", err)
	}

	conn, err := net.Dial("unix", unixSocketPath)
	if err != nil {
		return fmt.Errorf("failed to connect to unix socket: %v", err)
	}
	defer conn.Close()

	if _, err := conn.Write(jsonBytes); err != nil {
		return fmt.Errorf("failed to write bytes to socket conn: %v", err)
	}
	return nil
}
