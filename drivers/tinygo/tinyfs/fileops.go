//go:build pico || pico2

package tinyfs

import (
	"errors"
	"machine"
	"os"
	"strconv"

	"tinygo.org/x/drivers/sdcard"
	"tinygo.org/x/tinyfs/fatfs"
)

var errNotADirectory = errors.New("not a directory")

// Setp for a PicoCalc or hardware with the same SD Card Connection
var spi = machine.SPI0

const sckPin = machine.GP18
const sdoPin = machine.GP19
const sdiPin = machine.GP16
const csPin = machine.GP17

type FileOpsDriver struct {
	initErr error
	fs      *fatfs.FATFS
	// present working directory.  '/' should only be used
	// between directories
	pwd string
}

func (fo *FileOpsDriver) Init() error {
	sd := sdcard.New(spi, sckPin, sdoPin, sdiPin, csPin) // object allocated on the heap: escapes at line 39
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
	s, err := fo.fs.Stat(absPath(fo.pwd, path, false, false))
	if err != nil {
		return 0, err
	}
	return int(s.Size()), nil
}

func (fo *FileOpsDriver) ReadFile(path string) ([]byte, error) {
	if fo.initErr != nil {
		return nil, fo.initErr
	}
	f, err := fo.fs.Open(absPath(fo.pwd, path, false, false))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	ff := f.(*fatfs.File)
	sz, err := ff.Size()
	if err != nil {
		return nil, err
	}
	data := make([]byte, sz) // object allocated on the heap: size is not constant
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
	// Note, the tiyfs fat driver is very picky that all three of these are
	// set or it falls back to read mode.
	return fo.writeOrAppend(path, data, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
}

func (fo *FileOpsDriver) AppendToFile(path string, data []byte) error {
	// Note, the tiyfs fat driver is very picky that all three of these are
	// set or it falls back to read mode.
	return fo.writeOrAppend(path, data, os.O_WRONLY|os.O_CREATE|os.O_APPEND)
}

func (fo *FileOpsDriver) writeOrAppend(path string, data []byte, flags int) error { // object allocated on the heap: escapes at line 102
	if fo.initErr != nil {
		return fo.initErr
	}
	for len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	print("Opening file ")
	println(path)
	f, err := fo.fs.OpenFile(absPath(fo.pwd, path, false, false), flags)
	if err != nil {
		return err
	}
	defer f.Close()
	totalWritten := 0
	for totalWritten < len(data) {
		print("writing '")
		print(data[totalWritten:])
		println("'")
		written, err := f.Write(data[totalWritten:])
		if err != nil {
			return err
		}
		totalWritten += written
	}
	print("Wrote bytes:")
	println(strconv.Itoa(totalWritten))
	return nil
}

func (fo *FileOpsDriver) Chdir(path string) error {
	newPath := absPath(fo.pwd, path, false, false)
	s, err := fo.fs.Stat(newPath)
	if err != nil {
		return err
	}
	if !s.IsDir() {
		return errNotADirectory
	}
	fo.pwd = newPath
	return nil
}
