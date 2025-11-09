//go:build pico || pico2

package tinyfs

import (
	"tinygo.org/x/drivers/sdcard"
	"tinygo.org/x/tinyfs/littlefs"
)

func (fo *FileOpsDriver) initFS(sd sdcard.Device) {
	fs := littlefs.New(&sd)
	fs.Configure(&littlefs.Config{
		CacheSize:     512,
		LookaheadSize: 512,
		BlockCycles:   100,
	})
	fo.fs = fs
}
