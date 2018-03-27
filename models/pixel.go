package models

import "image/color"

type Pixel struct {
	X int
	Y int

	Color color.RGBA
}
