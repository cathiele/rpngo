//go:build fatfs

package tinyfs

import (
	"tinygo.org/x/drivers/sdcard"
	"tinygo.org/x/tinyfs/fatfs"
)

func (fo *FileOpsDriver) initFS(sd sdcard.Device) {
	fo.fs = fatfs.New(&sd)
	fo.fs.Configure(&fatfs.Config{
		SectorSize: 512,
	})
}
