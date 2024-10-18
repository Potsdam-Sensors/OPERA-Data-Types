package operadatatypes

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
