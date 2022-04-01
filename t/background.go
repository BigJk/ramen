package t

import (
	"github.com/BigJk/ramen"
	"github.com/BigJk/ramen/concolor"
)

// BackgroundTransform sets the background color of a cell
type BackgroundTransform struct {
	color concolor.Color
}

// Transform sets the background color of a cell
func (b BackgroundTransform) Transform(cell *ramen.Cell) error {
	cell.Background = b.color
	return nil
}

// Background creates a new transformer that sets the background color of a cell to the given color
func Background(newBackground concolor.Color) BackgroundTransform {
	return BackgroundTransform{newBackground}
}
