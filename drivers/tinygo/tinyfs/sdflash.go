//go:build sdflash

package tinyfs

import (
	"mattwach/rpngo/elog"

	"tinygo.org/x/drivers/sdcard"
	"tinygo.org/x/tinyfs"
)

func blockDevice() (tinyfs.BlockDevice, error) {
	elog.Heap("alloc: /drivers/tinygo/tinyfs/fileops.go:34: sd := sdcard.New(spi, sckPin, sdoPin, sdiPin, csPin)")
	sd := sdcard.New(spi, sckPin, sdoPin, sdiPin, csPin) // object allocated on the heap: escapes at line 39
	err := sd.Configure()
	if err != nil {
		return nil, err
	}
	return &sd, err
}
