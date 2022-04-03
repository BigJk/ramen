package t

import "github.com/BigJk/ramen"

// CellTransform sets the whole value of a cell
type CellTransform struct {
	cell ramen.Cell
}

// Transform sets the whole value of a cell
func (c CellTransform) Transform(cell *ramen.Cell) error {
	*cell = c.cell
	return nil
}

// Cell creates a transformer that sets the whole value of the cell
func Cell(newValue ramen.Cell) CellTransform {
	return CellTransform{newValue}
}
