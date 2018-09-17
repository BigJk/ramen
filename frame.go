package ramen

// DefaultFrame is a simple frame
var DefaultFrame = Frame{
	TopLeft:     int('+'),
	TopRight:    int('+'),
	BottomLeft:  int('+'),
	BottomRight: int('+'),
	Horizontal:  int('='),
	Vertical:    int('|'),

	Foreground: NewColor(255, 255, 255),
	Background: NewColor(0, 0, 0),
}

// Frame represents the look of a frame
type Frame struct {
	TopLeft     int
	TopRight    int
	BottomLeft  int
	BottomRight int
	Horizontal  int
	Vertical    int

	Foreground Color
	Background Color
}
