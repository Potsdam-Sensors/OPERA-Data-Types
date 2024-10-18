package operadatatypes

import "fmt"

type CsvFileWriteJob struct {
	Filename string
	Headers  string
	Content  string
}

func (c CsvFileWriteJob) String() string {
	return fmt.Sprintf("[File: %s, Headers: %s, Content: %s]", c.Filename, c.Headers, c.Content)
}

const (
	USB_MASS_STORAGE_UNIX_SOCKET = "/var/run/usb_mass.sock"
	MAIN_SD_UNIX_SOCKET          = "/var/run/main_sd.sock"
)
