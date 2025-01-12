package operadatatypes

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const CONFIG_FILE_LOCATION = "/etc/telosair/opera.conf"

type ConfigStruct struct {
	OutputToCsv bool `json:"output_to_csv"`
	OutputToRaw bool `json:"output_to_raw"`
}

func GetDefaultConfig() ConfigStruct {
	return ConfigStruct{
		true, true,
	}
}

func readConfigFile(filepath string) (ConfigStruct, error) {
	ret := ConfigStruct{}

	f, err := os.Open(filepath)
	if err != nil {
		return ret, fmt.Errorf("failed to open file, '%s': %v", filepath, err)
	}
	defer f.Close()
	buff, err := io.ReadAll(f)
	if err != nil {
		return ret, fmt.Errorf("failed to read file, '%s': %v", filepath, err)
	}
	if err := json.Unmarshal(buff, &ret); err != nil {
		return ret, fmt.Errorf("failed to read json contents of file, '%s': %v", filepath, err)
	}
	return ret, nil
}

func ReadConfig() (ConfigStruct, error) {
	return readConfigFile(CONFIG_FILE_LOCATION)
}
