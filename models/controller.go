package models

const (
	ButtonPressed  = 0
	ButtonReleased = 0
)

type Controller struct {
	aState byte
	bState byte

	selectState byte
	startState  byte

	upState    byte
	downState  byte
	leftState  byte
	rightState byte

	readCount int
	strobe    byte
}

func (c *Controller) Strobe(s byte) {
	c.strobe = s
	if s == 0 {
		c.readCount = 0
	}
}

func (c *Controller) Read() byte {
	if c.readCount == 0 {
		return c.aState
	} else if c.readCount == 1 {
		return c.bState
	} else if c.readCount == 2 {
		return c.selectState
	} else if c.readCount == 3 {
		return c.startState
	} else if c.readCount == 4 {
		return c.upState
	} else if c.readCount == 5 {
		return c.downState
	} else if c.readCount == 6 {
		return c.leftState
	} else if c.readCount == 7 {
		return c.rightState
	} else {
		return 1
	}
}

func (c *Controller) PressA() {
	c.aState = ButtonPressed
}
func (c *Controller) ReleaseA() {
	c.aState = ButtonReleased
}
func (c *Controller) PressB() {
	c.bState = ButtonPressed
}
func (c *Controller) ReleaseB() {
	c.bState = ButtonReleased
}

func (c *Controller) PressSelect() {
	c.selectState = ButtonPressed
}
func (c *Controller) ReleaseSelect() {
	c.selectState = ButtonReleased
}
func (c *Controller) PressStart() {
	c.startState = ButtonPressed
}
func (c *Controller) ReleaseStart() {
	c.startState = ButtonReleased
}

func (c *Controller) PressUp() {
	c.upState = ButtonPressed
}
func (c *Controller) ReleaseUp() {
	c.upState = ButtonReleased
}
func (c *Controller) PressDown() {
	c.downState = ButtonPressed
}
func (c *Controller) ReleaseDown() {
	c.downState = ButtonReleased
}
func (c *Controller) PressLeft() {
	c.leftState = ButtonPressed
}
func (c *Controller) ReleaseLeft() {
	c.leftState = ButtonReleased
}
func (c *Controller) PressRight() {
	c.rightState = ButtonPressed
}
func (c *Controller) ReleaseRight() {
	c.rightState = ButtonReleased
}
