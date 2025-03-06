package operadatatypes

import (
	"os"
	"testing"
)

func TestMlInputCounts(t *testing.T) {
	test_data_output := "./test_data.raw"
	test_data := mlPm25InputDataPulses{
		Laser:     12,
		Pd0:       13,
		Pd1:       14,
		MsRead:    102,
		Baseline0: 12.3,
		Baseline1: 54.3,
		Pulses: []NewPulse{
			NewPulse{[8]uint16{1, 2, 3, 4, 5, 6, 7, 8}, 100, 200},
		},
	}

	serial_bytes := test_data.Serialize()
	f, err := os.OpenFile(test_data_output, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	if err != nil {
		t.Errorf("failed to create test file: %v", err)
		return
	}
	defer f.Close()

	_, err = f.Write(serial_bytes)
	if err != nil {
		t.Errorf("failed to write bytes to test file: %v", err)
		return
	}
}
