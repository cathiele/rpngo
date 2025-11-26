//go:build littlefs

package tinyfs

import (
	"tinygo.org/x/tinyfs"
	"tinygo.org/x/tinyfs/littlefs"
)

func (fo *FileOpsDriver) initFS(sd tinyfs.BlockDevice) {
	fs := littlefs.New(sd)
	fs.Configure(&littlefs.Config{
		CacheSize:     512,
		LookaheadSize: 512,
		BlockCycles:   100,
	})
	fo.fs = fs
}
