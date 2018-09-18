package ramen

import "github.com/BigJk/ramen/consolecolor"

// Cell represents a cell in the console
type Cell struct {
	Foreground consolecolor.Color
	Background consolecolor.Color

	Char int
}
