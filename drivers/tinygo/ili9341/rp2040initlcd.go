//go:build pico || pico2

package ili9341

import (
	"image/color"
	"machine"

	"tinygo.org/x/drivers/ili9341"
)

// InitDisplay initializes the display of each board.
func InitDisplay() *ili9341.Device {
	machine.SPI1.Configure(machine.SPIConfig{
		SCK:       machine.GP10,
		SDO:       machine.GP11,
		SDI:       machine.GP12,
		Frequency: 40000000,
	})

	// configure backlight
	backlight := machine.GP9
	backlight.Configure(machine.PinConfig{machine.PinOutput})

	display := ili9341.NewSPI(
		machine.SPI1,
		machine.GP14, // LCD_DC,
		machine.GP13, // LCD_SS_PIN,
		machine.GP15, // LCD_RESET,
	)

	display.Configure(ili9341.Config{})
	backlight.High()
	display.SetRotation(ili9341.Rotation270)

	display.FillScreen(color.RGBA{255, 255, 255, 255})

	return display
}
