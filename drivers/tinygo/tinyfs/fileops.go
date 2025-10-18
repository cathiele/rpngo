//go:build pico || pico2

package tinyfs

import (
	"errors"
	"machine"
	"os"

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
	// present working directory.  Must start wth
	// '/' and not end with '/'
	pwd string
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
	fo.pwd = "/"
	return nil
}

func (fo *FileOpsDriver) FileSize(path string) (int, error) {
	if fo.initErr != nil {
		return 0, fo.initErr
	}
	s, err := fo.fs.Stat(fo.absPath(path))
	if err != nil {
		return 0, err
	}
	return int(s.Size()), nil
}

func (fo *FileOpsDriver) ReadFile(path string) ([]byte, error) {
	if fo.initErr != nil {
		return nil, fo.initErr
	}
	f, err := fo.fs.Open(fo.absPath(path))
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
	return fo.writeOrAppend(path, data, os.O_CREATE|os.O_WRONLY)
}

func (fo *FileOpsDriver) AppendToFile(path string, data []byte) error {
	return fo.writeOrAppend(path, data, os.O_CREATE|os.O_APPEND)
}

func (fo *FileOpsDriver) writeOrAppend(path string, data []byte, flags int) error {
	if fo.initErr != nil {
		return fo.initErr
	}
	f, err := fo.fs.OpenFile(fo.absPath(path), flags)
	if err != nil {
		return err
	}
	defer f.Close()
	totalWritten := 0
	for totalWritten < len(data) {
		written, err := f.Write(data[totalWritten:])
		if err != nil {
			return err
		}
		totalWritten += written
	}
	return nil
}

func (fo *FileOpsDriver) Chdir(path string) error {
	newPath := fo.absPath(path)
	if len(newPath) > 1 {
		for newPath[len(newPath)-1] == '/' {
			newPath = newPath[:len(newPath)-1]
		}
	}
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

// This doesn't handle more complex cases like /foo/bar/../baz
// Maybe it can be reimplemented later.
func (fo *FileOpsDriver) absPath(path string) string {
	if (len(path) == 0) || (path == ".") {
		return fo.pwd
	}
	if path == ".." {
		return fo.parentOfPwd()
	}
	if path[0] == '/' {
		return path
	}
	if fo.pwd == "/" {
		return fo.pwd + path
	}
	return fo.pwd + "/" + path
}

func (fo *FileOpsDriver) parentOfPwd() string {
	lastSlashIdx := 0
	for i, c := range fo.pwd {
		if c == '/' {
			lastSlashIdx = i
		}
	}
	if lastSlashIdx == 0 {
		return "/"
	}
	return fo.pwd[:lastSlashIdx]
}
