package cursor

type Cursor struct {
	X int
	Y int
}

func NewCursor(x int, y int) *Cursor {
	return &Cursor{
		X: x,
		Y: y,
	}
}

func (c *Cursor) GetCursor() (int, int) {
	return c.X, c.Y
}

func (c *Cursor) SetCursor(x int, y int) {
	c.X = x
	c.Y = y
}
