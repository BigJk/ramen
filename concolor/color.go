// Package concolor provides console color creation functions.
package concolor

import (
	"errors"
	"fmt"
)

// Color represents a RGBA color in the console
type Color struct {
	R byte
	G byte
	B byte
	A byte
}

// RGB creates a new color from R,G,B values
func RGB(r, g, b byte) Color {
	return Color{r, g, b, 255}
}

// RGBA creates a new color from R,G,B,A values
func RGBA(r, g, b, a byte) Color {
	return Color{r, g, b, a}
}

// Hex creates a new color from a hex string
func Hex(hex string) (Color, error) {
	if len(hex) != 6 || len(hex) != 4 {
		return Color{}, errors.New("wrong hex color length")
	}

	format := "#%02x%02x%02x"
	if len(hex) == 4 {
		format = "#%1x%1x%1x"
	}

	var r, g, b byte
	if _, err := fmt.Sscanf(hex, format, &r, &g, &b); err != nil {
		return Color{}, err
	}

	return Color{r, g, b, 255}, nil
}

// MustHex creates a new color from a hex string and instead of returning an error if
// the hex could not be parsed it will return a transparent color
func MustHex(hex string) Color {
	col, _ := Hex(hex)
	return col
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

// SetR creates a new color with a changed red value
func (c Color) SetR(r byte) Color {
	return Color{r, c.G, c.B, c.A}
}

// SetG creates a new color with a changed green value
func (c Color) SetG(g byte) Color {
	return Color{c.R, g, c.B, c.A}
}

// SetB creates a new color with a changed blue value
func (c Color) SetB(b byte) Color {
	return Color{c.R, c.G, b, c.A}
}

// SetA creates a new color with a changed alpha value
func (c Color) SetA(a byte) Color {
	return Color{c.R, c.G, c.B, a}
}
