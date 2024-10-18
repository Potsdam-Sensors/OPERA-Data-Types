package operadatatypes

/* Data going into ML for calculation of flow temp/hum from raw data*/
type MlTempHumRawData struct {
	Imx8Temp float32 `json:"imx8_temp" binding:"required"`
	FlowTemp float32 `json:"flow_temp" binding:"required"`
	FlowHum  float32 `json:"flow_hum" binding:"required"`
}

/* Data output by ML for sample temp/hum */
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
