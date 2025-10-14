//go:build pico || pico2

package ili948x

import (
	"machine"
)

// InitDisplay initializes the display of each board.
func InitDisplay() *Ili948x {
	machine.SPI1.Configure(machine.SPIConfig{
		SCK:       machine.GP10,
		SDO:       machine.GP11,
		SDI:       machine.GP12,
		Frequency: 40000000,
	})

	display := NewIli9488(
		NewSPITransport(*machine.SPI1),
		machine.GP13, // chip select
		machine.GP14, // data / command
		machine.GP15, // reset
	)

	display.FillScreen(0)

	return display
}
