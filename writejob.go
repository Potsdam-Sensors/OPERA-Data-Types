package operadatatypes

import (
	"fmt"
)

const (
	USB_MASS_STORAGE_UNIX_SOCKET = "/var/run/usb_mass.sock"
	MAIN_SD_UNIX_SOCKET          = "/var/run/main_sd.sock"

	DATA_TYPE_CSV_FILE = "C"
	DATA_TYPE_BIN_FILE = "B"
)

type FileWriteJob interface {
	String() string
	FileName() string
	SendGob(unixSocketPath string) error
}

type CsvFileWriteJob struct {
	Filename string
	Headers  string
	Content  string
}

func (c CsvFileWriteJob) String() string {
	return fmt.Sprintf("[File: %s, Headers: %s, Content: %s]", c.Filename, c.Headers, c.Content)
}

func (c CsvFileWriteJob) FileName() string {
	return c.Filename
}

func (c CsvFileWriteJob) SendGob(unixSocketPath string) error {
	return sendStructGob(c, DATA_TYPE_CSV_FILE, unixSocketPath)
}

type BinaryFileWriteJob struct {
	Filename string
	Content  []byte
}

func (b BinaryFileWriteJob) String() string {
	return fmt.Sprintf("[File: %s, Content: %d Bytes]", b.Filename, len(b.Content))
}

func (b BinaryFileWriteJob) FileName() string {
	return b.Filename
}

func (b BinaryFileWriteJob) SendGob(unixSocketPath string) error {
	return sendStructGob(b, DATA_TYPE_BIN_FILE, unixSocketPath)
}
