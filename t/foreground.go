package t

import (
	"github.com/BigJk/ramen"
	"github.com/BigJk/ramen/concolor"
)

// ForegroundTransform sets the foreground color of a cell
type ForegroundTransform struct {
	color concolor.Color
}

// Transform sets the foreground color of a cell
func (f ForegroundTransform) Transform(cell *ramen.Cell) error {
	cell.Foreground = f.color
	return nil
}

// Foreground creates a new transformer that sets the foreground color of a cell to the given color
func Foreground(newBackground concolor.Color) ForegroundTransform {
	return ForegroundTransform{newBackground}
}
