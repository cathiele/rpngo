//go:build pico

package ili948x

import (
	"image/color"
	"machine"
)

// InitDisplay initializes the display of each board.
func InitDisplay() *Ili948x {
	machine.SPI0.Configure(machine.SPIConfig{
		SCK:       machine.GP10,
		SDO:       machine.GP11,
		SDI:       machine.GP12,
		Frequency: 40000000,
	})

	display := NewIli9488(
		NewSPITransport(*machine.SPI1),
		machine.GP13,  // chip select
		machine.GP14,  // data / command
		machine.NoPin, // backlight
		machine.GP15,  // reset
		TFT_DEFAULT_WIDTH,
		TFT_DEFAULT_HEIGHT)

	display.FillScreen(color.RGBA{255, 255, 255, 255})

	return display
}
