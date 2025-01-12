package operadatatypes

import (
	"os"
	"testing"
)

const TEST_CONFIG_LOCATION = "./test_config.json"

func TestReadConfigFile(t *testing.T) {
	f, err := os.OpenFile(TEST_CONFIG_LOCATION, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0677)
	if err != nil {
		t.Errorf("failed to generate test file: %v", err)
		return
	}
	defer f.Close()

	testDataString := "{\"output_to_raw\": true, \"output_to_csv\": false}"
	if n, err := f.WriteString(testDataString); err != nil {
		t.Errorf("failed to write to test file: %v", err)
		return
	} else if n != len(testDataString) {
		t.Errorf("only wrote %d bytes out of expected %d bytes", n, len(testDataString))
		return
	}

	c, err := readConfigFile("./test_config.json")
	if err != nil {
		t.Errorf("readConfigFile(): %v", err)
		return
	}
	if (c.OutputToCsv != false) || (c.OutputToRaw != true) {
		t.Errorf("ConfigStruct was expected like %v, got %v", ConfigStruct{false, true}, c)
	}
}
