package t

import (
	"github.com/BigJk/ramen"
	"github.com/BigJk/ramen/consolecolor"
)

// BackgroundTransform sets the background color of a cell
type BackgroundTransform struct {
	color consolecolor.Color
}

// Transform sets the background color of a cell
func (b BackgroundTransform) Transform(cell *ramen.Cell) (bool, error) {
	if cell.Background == b.color {
		return false, nil
	}
	cell.Background = b.color
	return true, nil
}

// Background creates a new transformer that sets the background color of a cell to the given color
func Background(newBackground consolecolor.Color) BackgroundTransform {
	return BackgroundTransform{newBackground}
}
