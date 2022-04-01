package ramen

import "github.com/BigJk/ramen/concolor"

// Cell represents a cell in the console
type Cell struct {
	Foreground concolor.Color
	Background concolor.Color

	Char int
}
