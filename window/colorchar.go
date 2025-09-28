package window

// The character and color info packed into 1 16-bit value
//
// Format:
// RGGB RGGB CCCC CCCC
// FFFF BBBB HHHH HHHH
type ColorChar uint16

// rgb are 5 bit values
func NewColorCharFGColor(r, g, b uint16) ColorChar {
	return ColorChar(((r >> 4) << 15) | ((g >> 3) << 13) | ((b >> 4) << 12))
}

// rgb are 5 bit values
func NewColorCharBGColor(r, g, b uint16) ColorChar {
	return ColorChar(((r >> 4) << 11) | ((g >> 3) << 9) | ((b >> 4) << 8))
}

func (l ColorChar) Char() byte {
	return byte(l & 0xFF)
}

func (l ColorChar) FGColor() (uint8, uint8, uint8) {
	var r uint8 = 0
	if (l & 0x8000) != 0 {
		r = 0xFF
	}
	var g uint8 = 0
	switch l & 0x6000 {
	case 0x6000:
		g = 0xFF
	case 0x4000:
		g = 0xAA
	case 0x2000:
		g = 0x55
	}
	var b uint8 = 0
	if (l & 0x1000) != 0 {
		b = 0xFF
	}
	return r, g, b
}

func (l ColorChar) BGColor() (uint8, uint8, uint8) {
	var r uint8 = 0
	if (l & 0x800) != 0 {
		r = 0xFF
	}
	var g uint8 = 0
	switch l & 0x600 {
	case 0x600:
		g = 0xFF
	case 0x400:
		g = 0xAA
	case 0x200:
		g = 0x55
	}
	var b uint8 = 0
	if (l & 0x100) != 0 {
		b = 0xFF
	}
	return r, g, b
}
