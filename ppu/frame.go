package ppu

import (
	"image"

	"github.com/m4ntis/bones/models"
)

// frame is the glue between the ppu and it's renderers, letting it push pixels
// in each cycle, and finally creating the frame when rendering is done.
type frame struct {
	pixels []models.Pixel
}

func newFrame() *frame {
	return &frame{
		pixels: make([]models.Pixel, 256*240),
	}
}

func (f *frame) push(pix models.Pixel) {
	f.pixels[pix.X+pix.Y*256] = pix
}

func (f *frame) create() image.Image {
	frame := image.NewRGBA(image.Rect(0, 0, 256, 240))

	for _, pix := range f.pixels {
		frame.Set(pix.X, pix.Y, pix.Color)
	}

	f.pixels = make([]models.Pixel, 256*240)
	return frame
}
