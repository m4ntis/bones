// Package controller provides a controller struct with an api for
// controlling and reading its state
package controller

type button byte

const (
	released button = iota
	pressed
)

// Controller implements the NES controller, providing an API both for button
// manipulation as well as reading it's state by the NES.
type Controller struct {
	a button
	b button

	sel   button
	start button

	up    button
	down  button
	left  button
	right button

	readCount int
}

// Strobe is called by the CPU to reset the controller's internal read counter.
func (c *Controller) Strobe(s byte) {
	if s == 0 {
		c.readCount = 0
	}
}

// Read returns buttpon state determined by read count. Each Read increments
// internal read count by 1.
func (c *Controller) Read() byte {
	defer func() { c.readCount++ }()

	if c.readCount == 0 {
		return byte(c.a)
	} else if c.readCount == 1 {
		return byte(c.b)
	} else if c.readCount == 2 {
		return byte(c.sel)
	} else if c.readCount == 3 {
		return byte(c.start)
	} else if c.readCount == 4 {
		return byte(c.up)
	} else if c.readCount == 5 {
		return byte(c.down)
	} else if c.readCount == 6 {
		return byte(c.left)
	} else if c.readCount == 7 {
		return byte(c.right)
	}

	return 1
}

func (c *Controller) PressA() {
	c.a = pressed
}
func (c *Controller) ReleaseA() {
	c.a = released
}
func (c *Controller) PressB() {
	c.b = pressed
}
func (c *Controller) ReleaseB() {
	c.b = released
}

func (c *Controller) PressSelect() {
	c.sel = pressed
}
func (c *Controller) ReleaseSelect() {
	c.sel = released
}
func (c *Controller) PressStart() {
	c.start = pressed
}
func (c *Controller) ReleaseStart() {
	c.start = released
}

func (c *Controller) PressUp() {
	c.up = pressed
}
func (c *Controller) ReleaseUp() {
	c.up = released
}
func (c *Controller) PressDown() {
	c.down = pressed
}
func (c *Controller) ReleaseDown() {
	c.down = released
}
func (c *Controller) PressLeft() {
	c.left = pressed
}
func (c *Controller) ReleaseLeft() {
	c.left = released
}
func (c *Controller) PressRight() {
	c.right = pressed
}
func (c *Controller) ReleaseRight() {
	c.right = released
}
