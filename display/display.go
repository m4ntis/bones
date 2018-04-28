package display

import (
	"image"
	"image/draw"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/m4ntis/bones/controller"
	"golang.org/x/image/colornames"
)

var (
	width  = 256
	height = 240
	scale  = float64(4)
)

type Display struct {
	imgc chan image.Image

	ctrl *controller.Controller
}

func New(ctrl *controller.Controller) *Display {
	return &Display{
		imgc: make(chan image.Image),

		ctrl: ctrl,
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
		d.updateCtrl(win)

		img := <-d.imgc
		p := pixel.PictureDataFromImage(img)
		s := pixel.NewSprite(p, p.Bounds())

		win.Clear(colornames.White)
		s.Draw(win, pixel.IM.Moved(c).Scaled(c, scale))
		win.Update()
	}
}

func (d *Display) updateCtrl(win *pixelgl.Window) {
	if win.Pressed(pixelgl.KeyX) {
		d.ctrl.PressA()
	} else {
		d.ctrl.ReleaseA()
	}
	if win.Pressed(pixelgl.KeyZ) {
		d.ctrl.PressB()
	} else {
		d.ctrl.ReleaseB()
	}
	if win.Pressed(pixelgl.KeyA) {
		d.ctrl.PressSelect()
	} else {
		d.ctrl.ReleaseSelect()
	}
	if win.Pressed(pixelgl.KeyS) {
		d.ctrl.PressStart()
	} else {
		d.ctrl.ReleaseStart()
	}
	if win.Pressed(pixelgl.KeyUp) {
		d.ctrl.PressUp()
	} else {
		d.ctrl.ReleaseUp()
	}
	if win.Pressed(pixelgl.KeyDown) {
		d.ctrl.PressDown()
	} else {
		d.ctrl.ReleaseDown()
	}
	if win.Pressed(pixelgl.KeyLeft) {
		d.ctrl.PressLeft()
	} else {
		d.ctrl.ReleaseLeft()
	}
	if win.Pressed(pixelgl.KeyRight) {
		d.ctrl.PressRight()
	} else {
		d.ctrl.ReleaseRight()
	}
}
