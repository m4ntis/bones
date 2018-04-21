package ppu

import (
	"image"
	"image/color"
)

// pixel de
type pixel struct {
	x int
	y int

	color color.RGBA
}

// frame is the glue between the ppu and it's renderers, letting it push pixels
// in each cycle, and finally creating the frame when rendering is done.
type frame struct {
	pixels []pixel
}

func newFrame() *frame {
	return &frame{
		pixels: make([]pixel, 256*240),
	}
}

func (f *frame) push(pix pixel) {
	f.pixels[pix.x+pix.y*256] = pix
}

func (f *frame) create() image.Image {
	frame := image.NewRGBA(image.Rect(0, 0, 256, 240))

	for _, pix := range f.pixels {
		frame.Set(pix.x, pix.y, pix.color)
	}

	f.pixels = make([]pixel, 256*240)
	return frame
}
