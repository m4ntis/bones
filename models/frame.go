package models

import (
	"image"
)

type Frame struct {
	pixels []Pixel
}

func (f *Frame) Push(pix Pixel) {
	f.pixels = append(f.pixels, pix)
}

func (f *Frame) Create() image.Image {
	defer func() { f.pixels = make([]Pixel, 256*240) }()

	frame := image.NewRGBA(image.Rect(0, 0, 256, 240))
	for _, pix := range f.pixels {
		frame.Set(pix.X, pix.Y, pix.Color)
	}

	return frame
}
