package models

import (
	"image"
)

type Frame struct {
	pixels []Pixel
}

func NewFrame() *Frame {
	return &Frame{
		pixels: make([]Pixel, 256*240),
	}
}

func (f *Frame) Push(pix Pixel) {
	f.pixels[pix.X+pix.Y*256] = pix
}

func (f *Frame) Create() image.Image {
	frame := image.NewRGBA(image.Rect(0, 0, 256, 240))

	for _, pix := range f.pixels {
		frame.Set(pix.X, pix.Y, pix.Color)
	}

	f.pixels = make([]Pixel, 256*240)
	return frame
}
