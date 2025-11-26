//go:build fatfs

package tinyfs

import (
	"tinygo.org/x/tinyfs"
	"tinygo.org/x/tinyfs/fatfs"
)

func (fo *FileOpsDriver) initFS(sd tinyfs.BlockDevice) {
	fs := fatfs.New(sd)
	fs.Configure(&fatfs.Config{
		SectorSize: 512,
	})
	fo.fs = fs
}
