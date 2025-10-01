package main

import (
	"image/color"
	"mattwach/rpngo/drivers/tinygo/ili9341"
	"runtime"
	"strconv"

	"mattwach/rpngo/drivers/tinygo/fonts"
)

func freeMemOverlay(screen ili9341.Ili9341Screen) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	s := strconv.Itoa(int(ms.HeapIdle/1024)) + "k free"
	var x int16 = 280
	var y int16 = 20
	g := color.RGBA{G: 255}
	for _, r := range s {
		fonts.NimbusMono12p.GetGlyph(rune(r&0xFF)).Draw(screen.Device, x, y, g)
		x += ili9341.FontCharWidth
	}
}
