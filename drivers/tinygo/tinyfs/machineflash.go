//go:build machineflash

package tinyfs

import (
	"machine"

	"tinygo.org/x/tinyfs"
)

func blockDevice() (tinyfs.BlockDevice, error) {
	return machine.Flash, nil
}
