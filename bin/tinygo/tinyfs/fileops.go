//go:build pico || pico2

package tinyfs

import (
	"machine"
	"mattwach/rpngo/rpn"

	"tinygo.org/x/drivers/sdcard"
	"tinygo.org/x/tinyfs/fatfs"
)

// Setp for a PicoCalc or hardware with the same SD Card Connection
var spi = machine.SPI0

const sckPin = machine.GP18
const sdoPin = machine.GP19
const sdiPin = machine.GP16
const csPin = machine.GP17

type FileOpsDriver struct {
	initErr error
	fs      *fatfs.FATFS
}

func (fo *FileOpsDriver) Init() error {
	sd := sdcard.New(spi, sckPin, sdoPin, sdiPin, csPin)
	fo.initErr = sd.Configure()
	if fo.initErr != nil {
		return fo.initErr
	}
	fo.fs = fatfs.New(&sd)
	fo.fs.Configure(&fatfs.Config{
		SectorSize: 512,
	})
	return nil
}

func (fo *FileOpsDriver) FileSize(path string) (int, error) {
	if fo.initErr != nil {
		return 0, fo.initErr
	}
	s, err := fo.fs.Stat(path)
	if err != nil {
		return 0, err
	}
	return int(s.Size()), nil
}

func (fo *FileOpsDriver) ReadFile(path string) ([]byte, error) {
	if fo.initErr != nil {
		return nil, fo.initErr
	}
	f, err := fo.fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	ff := f.(*fatfs.File)
	sz, err := ff.Size()
	if err != nil {
		return nil, err
	}
	data := make([]byte, sz)
	totalRead := 0
	for totalRead < int(sz) {
		read, err := f.Read(data[totalRead:])
		if err != nil {
			return nil, err
		}
		totalRead += read
	}
	return data, nil
}

func (fo *FileOpsDriver) WriteFile(path string, data []byte) error {
	// TODO: Implement
	return rpn.ErrNotSupported
}

func (fo *FileOpsDriver) AppendToFile(path string, data []byte) error {
	// TODO: Implement
	return rpn.ErrNotSupported
}

func (fo *FileOpsDriver) Chdir(path string) error {
	// TODO: Implement
	return rpn.ErrNotSupported
}
