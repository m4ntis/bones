// Package display implements a simple OpenGL PPU Display.
package display

import (
	"fmt"
	"image"
	"image/draw"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/m4ntis/bones/controller"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

var (
	width  = 256
	height = 240
)

// Display is BoNES' implementation of a simple OpenGL PPU display.
type Display struct {
	imgc chan image.Image

	ctrl *controller.Controller

	scale float64

	fps bool

	frameCount    int
	lastFPSUpdate time.Time
}

// New returns an instance of a Display, initialized with a controller.
//
// ctrl is controlled by the display and read by the NES.
//
// fps determines whether to display a small fps counter in the bottom of the
// display.
func New(ctrl *controller.Controller, fps bool, scale float64) *Display {
	return &Display{
		imgc: make(chan image.Image),

		ctrl: ctrl,

		scale: scale,

		fps: fps,
	}
}

// Display sets the image to be displayed.
func (d *Display) Display(img image.Image) {
	r := image.Rect(0, 0, width, height)
	cropped := image.NewRGBA(r)
	draw.Draw(cropped, r, img, image.ZP, draw.Src)
	d.imgc <- cropped
}

// Run starts displaying the set images.
//
// IMPORTANT: As the implementation of Display uses OpenGL, Run must be called
// from the main goroutine. This method blocks forever, and the remaining logic
// must be run in a different goroutine.
func (d *Display) Run() {
	pixelgl.Run(d.run)
}

func (d *Display) run() {
	d.lastFPSUpdate = time.Now()
	d.frameCount = 0

	win := d.createWindow()
	center := win.Bounds().Center()

	if d.fps {
		d.runWithFPS(win, center)
		return
	}

	d.runWithoutFPS(win, center)
}

func (d *Display) runWithFPS(win *pixelgl.Window, center pixel.Vec) {
	fpsTxt := initTxt()
	fpsTxtBgr := initTxtRect()

	for !win.Closed() {
		d.updateFPS(fpsTxt)
		d.updateCtrl(win)
		d.displayNextFrameWithFPS(win, center, fpsTxt, fpsTxtBgr)
	}
}

func (d *Display) runWithoutFPS(win *pixelgl.Window, center pixel.Vec) {
	for !win.Closed() {
		d.updateCtrl(win)
		d.displayNextFrame(win, center)
	}
}

func (d *Display) updateFPS(txt *text.Text) {
	d.frameCount++

	if time.Now().Sub(d.lastFPSUpdate) >= time.Second {
		d.lastFPSUpdate = d.lastFPSUpdate.Add(time.Second)
		txt.Clear()
		fmt.Fprintln(txt, d.frameCount)
		d.frameCount = 0
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

func (d *Display) displayNextFrame(win *pixelgl.Window, center pixel.Vec) {
	img := <-d.imgc

	p := pixel.PictureDataFromImage(img)
	s := pixel.NewSprite(p, p.Bounds())

	win.Clear(colornames.White)
	s.Draw(win, pixel.IM.Moved(center).Scaled(center, d.scale))
	win.Update()
}

func (d *Display) displayNextFrameWithFPS(win *pixelgl.Window, center pixel.Vec,
	fpsTxt *text.Text, fpsTxtBgr *imdraw.IMDraw) {

	img := <-d.imgc

	p := pixel.PictureDataFromImage(img)
	s := pixel.NewSprite(p, p.Bounds())

	win.Clear(colornames.White)

	s.Draw(win, pixel.IM.Moved(center).Scaled(center, d.scale))
	fpsTxtBgr.Draw(win)
	fpsTxt.Draw(win, pixel.IM)

	win.Update()
}

func (d *Display) createWindow() *pixelgl.Window {
	cfg := pixelgl.WindowConfig{
		Title:  "BoNES",
		Bounds: pixel.R(0, 0, float64(width)*d.scale, float64(height)*d.scale),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	return win
}

func initTxt() *text.Text {
	txt := text.New(pixel.V(3, 5),
		text.NewAtlas(basicfont.Face7x13, text.ASCII))
	txt.Color = pixel.RGB(200, 0, 0)

	return txt
}

func initTxtRect() *imdraw.IMDraw {
	fpsTxtBgr := imdraw.New(nil)
	fpsTxtBgr.Push(pixel.V(0, 0))
	fpsTxtBgr.Push(pixel.V(20, 20))
	fpsTxtBgr.Color = pixel.RGB(1, 1, 1)
	fpsTxtBgr.Rectangle(0)

	return fpsTxtBgr
}
