package io

import (
	"fmt"
	"image"
	"time"
)

type BenchDisplay struct {
	frameCount int

	lastFPSUpdate time.Time
}

func NewBenchDisplay() *BenchDisplay {
	return &BenchDisplay{0, time.Now()}
}

// Display increments frame count.
func (d *BenchDisplay) Display(img image.Image) {
	d.frameCount++

	if time.Now().Sub(d.lastFPSUpdate) >= time.Second {
		d.lastFPSUpdate = d.lastFPSUpdate.Add(time.Second)
		fmt.Println("Frames per second:", d.frameCount)
		d.frameCount = 0
	}
}
