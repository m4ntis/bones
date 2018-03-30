package display

import (
	"image"
	"image/draw"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	width  = 256
	height = 240
	scale  = float64(0.5)
)

type Display struct {
	imgc chan image.Image
}

func New() *Display {
	return &Display{
		imgc: make(chan image.Image),
	}
}

func (d *Display) Display(img image.Image) {
	r := image.Rect(0, 0, width, height)
	cropped := image.NewRGBA(r)
	draw.Draw(cropped, r, img, image.ZP, draw.Src)
	d.imgc <- cropped
}

// Run must be called from the main goroutine
func (d *Display) Run() {
	pixelgl.Run(d.run)
}

func (d *Display) run() {
	cfg := pixelgl.WindowConfig{
		Title:  "BoNES",
		Bounds: pixel.R(0, 0, float64(width)*scale, float64(height)*scale),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	c := win.Bounds().Center()

	for !win.Closed() {
		win.Clear(colornames.White)

		img := <-d.imgc
		p := pixel.PictureDataFromImage(img)
		pixel.NewSprite(p, p.Bounds()).
			Draw(win, pixel.IM.Moved(c).Scaled(c, scale))

		win.Update()
	}
}
