package ramen

// Color represents a ARGB color in the console
type Color struct {
	R byte
	G byte
	B byte
	A byte
}

// NewColor creates a new color from R,G,B values
func NewColor(r, g, b byte) Color {
	return Color{r, g, b, 255}
}

// NewColorTransparent creates a new color from R,G,B,A values
func NewColorTransparent(r, g, b, a byte) Color {
	return Color{r, g, b, a}
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
