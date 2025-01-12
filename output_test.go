package operadatatypes

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func comparePointers(v1, v2 interface{}) bool {
	if v1 == nil && v2 == nil {
		return true
	}
	if v1 == nil || v2 == nil {
		return false
	}
	return v1 == v2
}
func checkPulseEquality(original, nuevo NewPulse) error {
	if original.RawPeak != nuevo.RawPeak {
		return fmt.Errorf("Pulse structs have differing 'RawPeak', old is %d, got %d", original.RawPeak, nuevo.RawPeak)
	}
	if original.SidePeak != nuevo.SidePeak {
		return fmt.Errorf("Pulse structs have differing 'SidePeak', old is %d, got %d", original.SidePeak, nuevo.SidePeak)
	}
	for idx := 0; idx < 8; idx++ {
		if original.Indices[idx] != nuevo.Indices[idx] {
			return fmt.Errorf("Pulses structs have differing 'Indices': old is %v, got %v", original.Indices, nuevo.Indices)
		}
	}
	return nil
}
func checkCountsEquality(original, nuevo NewTeensyCounts) error {
	for _, test := range []struct {
		origVal interface{}
		newVal  interface{}
		name    string
	}{
		{original.PinPd0, nuevo.PinPd0, "PinPd0"},
		{original.PinPd1, nuevo.PinPd1, "PinPd1"},
		{original.PinLaser, nuevo.PinLaser, "PinLaser"},
		{original.RawScalar0, nuevo.RawScalar0, "RawScalar0"},
		{original.RawScalar1, nuevo.RawScalar1, "RawScalar1"},
		{original.DiffedScalar0, nuevo.DiffedScalar0, "DiffedScalar0"},
		{original.DiffedScalar1, nuevo.DiffedScalar1, "DiffedScalar1"},
		{original.Baseline0, nuevo.Baseline0, "Baseline0"},
		{original.Baseline1, nuevo.Baseline1, "Baseline1"},
		{original.RawUpperTh0, nuevo.RawUpperTh0, "RawUpperTh0"},
		{original.RawUpperTh1, nuevo.RawUpperTh1, "RawUpperTh1"},
		{original.DiffedUpperTh0, nuevo.DiffedUpperTh0, "DiffedUpperTh0"},
		{original.DiffedUpperTh1, nuevo.DiffedUpperTh1, "DiffedUpperTh1"},
		{original.MsRead, nuevo.MsRead, "MsRead"},
		{original.BuffersRead, nuevo.BuffersRead, "BuffersRead"},
		{original.NumPulses, nuevo.NumPulses, "NumPulses"},
		{original.MaxLaserOn, nuevo.MaxLaserOn, "MaxLaserOn"},
		{original.PulsesPerSecond, nuevo.PulsesPerSecond, "PulsesPerSecond"},
		{len(original.Pulses), len(nuevo.Pulses), "Pulses (length)"},
	} {
		if !comparePointers(test.origVal, test.newVal) {
			return fmt.Errorf("Struct's '%s' members do not match, got '%v' originally, new is '%v'.", test.name, test.origVal, test.newVal)
		}
	}
	for idx := 0; idx < len(original.Pulses); idx++ {
		if err := checkPulseEquality(original.Pulses[idx], nuevo.Pulses[idx]); err != nil {
			return fmt.Errorf("Struct has differing pulse #%d, old and new are:\n\t->%v\n\t->%v\r\n%v", idx, original.Pulses[idx], nuevo.Pulses[idx], err)
		}
	}
	return nil
}

func checkPrimaryStructEquality(original, nuevo PrimaryData) error {
	for _, test := range []struct {
		origVal interface{}
		newVal  interface{}
		name    string
	}{
		{original.PortentaSerial, nuevo.PortentaSerial, "PortentaSerial"},
		{original.TeensyData.UnixSec, nuevo.TeensyData.UnixSec, "UnixSec"},
		{original.TeensyData.MilliSec, nuevo.TeensyData.MilliSec, "MilliSec"},
		{original.TeensyData.McuTemp, nuevo.TeensyData.McuTemp, "McuTemp"},
		{original.TeensyData.FlowTemp, nuevo.TeensyData.FlowTemp, "FlowTemp"},
		{original.TeensyData.FlowHum, nuevo.TeensyData.FlowHum, "FlowHum"},
		{original.TeensyData.FlowRate, nuevo.TeensyData.FlowRate, "FlowRate"},
		{original.TeensyData.HvEnabled, nuevo.TeensyData.HvEnabled, "HvEnabled"},
		{len(original.TeensyData.Counts), len(nuevo.TeensyData.Counts), "Counts (length)"},
	} {
		if !comparePointers(test.origVal, test.newVal) {
			return fmt.Errorf("Struct's '%s' members do not match, got '%v' originally, new is '%v'.", test.name, test.origVal, test.newVal)
		}
	}
	for idx := 0; idx < len(original.TeensyData.Counts); idx++ {
		if err := checkCountsEquality(*original.TeensyData.Counts[idx], *nuevo.TeensyData.Counts[idx]); err != nil {
			return fmt.Errorf("Struct's counts #%d are not equal, from old to new, they are: \n\t-> %v\n\t-> %v\r\n%v", idx, original.TeensyData.Counts[idx], nuevo.TeensyData.Counts[idx], err)
		}
	}
	return nil
}

func TestPrimaryPackUnpack(t *testing.T) {
	testData := &PrimaryData{
		PortentaSerial: "abcdefg12345",
		TeensyData: NewTeensyData{
			UnixSec:   uint32(time.Now().Unix()),
			MilliSec:  1002,
			McuTemp:   24.3,
			FlowTemp:  1,
			FlowHum:   2,
			FlowRate:  3,
			HvEnabled: true,
			HvSet:     12,
			HvMonitor: 333,
			Counts: []*NewTeensyCounts{
				{
					PinPd0:          1,
					PinPd1:          2,
					PinLaser:        99,
					RawScalar0:      12,
					RawScalar1:      13,
					DiffedScalar0:   14,
					DiffedScalar1:   100,
					Baseline0:       22.4,
					Baseline1:       -12.1,
					RawUpperTh0:     100.1,
					RawUpperTh1:     12.22,
					DiffedUpperTh0:  -1,
					DiffedUpperTh1:  10,
					MsRead:          255,
					BuffersRead:     254,
					NumPulses:       1,
					MaxLaserOn:      99,
					PulsesPerSecond: 100,
					Pulses: []NewPulse{
						{
							Indices:  [8]uint16{1, 2, 3, 412, 5, 6, 7, 8},
							RawPeak:  25,
							SidePeak: 20,
						},
						{
							Indices:  [8]uint16{1, 2, 3, 4, 5, 6, 7, 8},
							RawPeak:  255,
							SidePeak: 21,
						},
					},
				}, {}, {},
			},
		},
	}

	buffer := new(bytes.Buffer)
	testData.Pack(buffer)
	newStruct := &PrimaryData{}
	if err := newStruct.Unpack(buffer); err != nil {
		t.Errorf("Got an error unpacking newStruct: %v", err)
	}

	if err := checkPrimaryStructEquality(*testData, *newStruct); err != nil {
		t.Error(err)
		return
	}
}

/* Secondary */

func checkSecondaryStructEquality(o, n SecondaryData) error {
	for _, test := range []struct {
		origVal interface{}
		newVal  interface{}
		name    string
	}{
		{o.UnixSec, n.UnixSec, "UnixSec"},
		{o.PortentaSerial, n.PortentaSerial, "PortentaSerial"},
		{o.Sps30, n.Sps30, "Sps30"},
		{o.Pressure, n.Pressure, "Pressure"},
		{o.Co2, n.Co2, "Co2"},
		{o.VocIndex, n.VocIndex, "VocIndex"},
		{o.FlowTemperature, n.FlowTemperature, "FlowTemperature"},
		{o.FlowHumidity, n.FlowHumidity, "FlowHumidity"},
		{o.FlowRate, n.FlowRate, "FlowRate"},
		{o.PortentaImx8Temp, n.PortentaImx8Temp, "PortentaImx8Temp"},
		{o.TeensyMcuTemp, n.TeensyMcuTemp, "TeensyMcuTemp"},
		{o.OpticalTemperatures[0], n.OpticalTemperatures[0], "OpticalTemperatures[0]"},
		{o.OpticalTemperatures[1], n.OpticalTemperatures[1], "OpticalTemperatures[1]"},
		{o.OpticalTemperatures[2], n.OpticalTemperatures[2], "OpticalTemperatures[2]"},
		{o.OmbTemperatureHtu, n.OmbTemperatureHtu, "OmbTemperatureHtu"},
		{o.OmbHumidityHtu, n.OmbHumidityHtu, "OmbHumidityHtu"},
		{o.OmbTemperatureScd, n.OmbTemperatureScd, "OmbTemperatureScd"},
		{o.OmbHumidityScd, n.OmbHumidityScd, "OmbHumidityScd"},
		{o.Monitor5vMean, n.Monitor5vMean, "Monitor5vMean"},
		{o.Monitor5vStdDev, n.Monitor5vStdDev, "Monitor5vStdDev"},
	} {
		if !comparePointers(test.origVal, test.newVal) {
			return fmt.Errorf("Struct '%s' field differs old to new: '%v' vs '%v'", test.name, test.origVal, test.newVal)
		}
	}
	return nil
}

func TestPackUnpackSecondary(t *testing.T) {
	testData := &SecondaryData{
		UnixSec:             12,
		PortentaSerial:      "abcdefg",
		Pressure:            101.2,
		Co2:                 1000,
		VocIndex:            14,
		FlowTemperature:     -10.2,
		FlowHumidity:        1.0,
		FlowRate:            -2,
		PortentaImx8Temp:    100,
		TeensyMcuTemp:       32,
		OpticalTemperatures: [3]float32{-1, 2, 100.1},
		OmbTemperatureHtu:   1,
		OmbHumidityHtu:      2,
		OmbTemperatureScd:   22,
		OmbHumidityScd:      5,
		Monitor5vMean:       10,
		Monitor5vStdDev:     3.1,
	}
	newStruct := &SecondaryData{}

	buffer := new(bytes.Buffer)
	testData.Pack(buffer)

	if err := newStruct.Unpack(buffer); err != nil {
		t.Errorf("Got an error unpacking newStruct: %v", err)
	}

	if err := checkSecondaryStructEquality(*testData, *newStruct); err != nil {
		t.Error(err)
		return
	}
}

func checkOperaDataEquality(o, n OperaData) error {
	for _, test := range []struct {
		origVal interface{}
		newVal  interface{}
		name    string
	}{
		{o.UnixSec, n.UnixSec, "UnixSec"},
		{o.PortentaSerial, n.PortentaSerial, "PortentaSerial"},
		{o.Pm2p5, n.Pm2p5, "Pm2p5"},
		{o.ClassLabel, n.ClassLabel, "ClassLabel"},
		{len(o.ClassLabels), len(n.ClassLabels), "ClassLabels (length)"},
		{o.Temp, n.Temp, "Temp"},
		{o.RH, n.RH, "RH"},
		{o.Sps30Pm2p5, n.Sps30Pm2p5, "Sps30Pm2p5"},
		{o.Pressure, n.Pressure, "Pressure"},
		{o.Co2, n.Co2, "Co2"},
		{o.VocIndex, n.VocIndex, "VocIndex"},
	} {
		if !comparePointers(test.origVal, test.newVal) {
			return fmt.Errorf("Struct '%s' field differs old to new: '%v' vs '%v'", test.name, test.origVal, test.newVal)
		}
	}
	return nil
}

func TestPackUnpackOperaData(t *testing.T) {
	testData := &OperaData{
		UnixSec:        12,
		PortentaSerial: "abcdefg",
		Pm2p5:          -0.5,
		ClassLabel:     "Lemons",
		ClassLabels: []string{
			"crocodiles",
			"alligators",
			"handbags",
		},
		ClassProbs: []float32{
			.3,
			.2,
			.5111,
		},
		Temp:       199.2,
		RH:         -.3,
		Sps30Pm2p5: 12.12,
		Pressure:   101.2,
		Co2:        1000,
		VocIndex:   14,
	}

	newStruct := &OperaData{}

	buffer := new(bytes.Buffer)
	testData.Pack(buffer)

	if err := newStruct.Unpack(buffer); err != nil {
		t.Errorf("Got an error unpacking newStruct: %v", err)
	}

	if err := checkOperaDataEquality(*testData, *newStruct); err != nil {
		t.Error(err)
		return
	}
}
