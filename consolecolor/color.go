package consolecolor

import "fmt"

// Color represents a ARGB color in the console
type Color struct {
	R byte
	G byte
	B byte
	A byte
}

// New creates a new color from R,G,B values
func New(r, g, b byte) Color {
	return Color{r, g, b, 255}
}

// NewTransparent creates a new color from R,G,B,A values
func NewTransparent(r, g, b, a byte) Color {
	return Color{r, g, b, a}
}

// NewHex creates a new color from a hex string
func NewHex(hex string) Color {
	format := "#%02x%02x%02x"
	if len(hex) == 4 {
		format = "#%1x%1x%1x"
	}

	var r, g, b byte
	fmt.Sscanf(hex, format, &r, &g, &b)

	return Color{r, g, b, 255}
}

// RGBA returns the color values as uint32s
func (c Color) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = uint32(c.A)
	a |= a << 8
	return
}

// Floats returns the color values as floats (0f - 1f)
func (c Color) Floats() (r, g, b, a float64) {
	return float64(c.R) / 0xff, float64(c.G) / 0xff, float64(c.B) / 0xff, float64(c.A) / 0xff
}
