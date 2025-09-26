// Makes the pixel package compatible with the displayer interface
//
// It seem like this has to be implemented somewhere, but I can't yet find it.
package pixel565

import (
	"image/color"

	"tinygo.org/x/drivers/pixel"
)

type Pixel565 struct {
	Image pixel.Image[pixel.RGB565BE]
	// these are for quicker SetPixel chexks (to avoid a panic)
	w int16
	h int16
}

func (p *Pixel565) Init(w, h int16) {
	p.Image = pixel.NewImage[pixel.RGB565BE](int(w), int(h))
	p.w = w
	p.h = h
}

func (p *Pixel565) Size() (int16, int16) {
	return p.w, p.h
}

func (p *Pixel565) SetPixel(x, y int16, c color.RGBA) {
	if (x < 0) || (x >= p.w) || (y < 0) || (y >= p.h) {
		return
	}
	p.Image.Set(int(x), int(y), pixel.NewRGB565BE(c.R, c.G, c.B))
}

func (p *Pixel565) Display() error {
	return nil
}
