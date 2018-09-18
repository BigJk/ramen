package t

import "github.com/BigJk/ramen"

// Transformer is a interface that specifies transformations (changes) that
// can be applied to cells in  a console
type Transformer interface {
	Transform(cell *ramen.Cell) (bool, error)
}
