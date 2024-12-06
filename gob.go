package operadatatypes

import (
	"encoding/gob"
	"fmt"
	"net"
)

func sendStructGob(d interface{}, dataIdentifier string, unixSocketPath string) error {
	conn, err := net.Dial("unix", unixSocketPath)
	if err != nil {
		return fmt.Errorf("failed to connect to socket, %s: %v", unixSocketPath, err)
	}
	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	if err := encoder.Encode(dataIdentifier); err != nil {
		return fmt.Errorf("failed to send data type: %v", err)
	}
	if err := encoder.Encode(d); err != nil {
		return fmt.Errorf("failed to send data: %v", err)
	}
	return nil
}

func ReceiveStructGob(conn net.Conn) (interface{}, error) {
	decoder := gob.NewDecoder(conn)

	/* Get message type */
	var msgType string
	if err := decoder.Decode(&msgType); err != nil {
		return nil, fmt.Errorf("failed to decode msg type: %v", err)
	}

	/* Interpret Data */
	var data interface{}
	var dataTypeName string

	switch msgType {
	case DATA_TYPE_SPS30:
		data = &Sps30Data{}
		dataTypeName = "sps30"
	case DATA_TYPE_M4_SENSORS:
		data = &M4SensorMeasurement{}
		dataTypeName = "m4 sensor"
	case DATA_TYPE_TEENSY:
		data = &NewTeensyData{}
		dataTypeName = "teensy raw"
	case DATA_TYPE_ML_TEMP_RH:
		data = &MlTempHumOutputData{}
		dataTypeName = "ml temp/rh"
	case DATA_TYPE_ML_PRIMARY:
		data = &MlPm25OutputData{}
		dataTypeName = "ml pm25"
	default:
		return nil, fmt.Errorf("recieved unknown datatype: %v", msgType)
	}

	if err := decoder.Decode(data); err != nil {
		return nil, fmt.Errorf("failed to decode %s data: %v", dataTypeName, err)
	}
	return data, nil
}
