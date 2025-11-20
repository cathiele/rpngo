//go:build fatfs

package tinyfs

import (
	"tinygo.org/x/drivers/sdcard"
	"tinygo.org/x/tinyfs/fatfs"
)

func (fo *FileOpsDriver) initFS(sd sdcard.Device) {
	fs := fatfs.New(&sd)
	fs.Configure(&fatfs.Config{
		SectorSize: 512,
	})
	fo.fs = fs
}
