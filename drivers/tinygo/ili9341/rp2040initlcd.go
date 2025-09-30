//go:build pico

package ili9341

import (
	"image/color"
	"machine"

	"tinygo.org/x/drivers/ili9341"
)

// InitDisplay initializes the display of each board.
func InitDisplay() *ili9341.Device {
	machine.SPI0.Configure(machine.SPIConfig{
		SCK:       machine.SPI0_SCK_PIN,
		SDO:       machine.SPI0_SDO_PIN,
		SDI:       machine.SPI0_SDI_PIN,
		Frequency: 40000000,
	})

	// configure backlight
	backlight := machine.GP9
	backlight.Configure(machine.PinConfig{machine.PinOutput})

	display := ili9341.NewSPI(
		machine.SPI0,
		machine.GP10, // LCD_DC,
		machine.GP11, // LCD_SS_PIN,
		machine.GP12, // LCD_RESET,
	)

	display.Configure(ili9341.Config{})
	backlight.High()
	display.SetRotation(ili9341.Rotation270)

	display.FillScreen(color.RGBA{255, 255, 255, 255})

	return display
}
