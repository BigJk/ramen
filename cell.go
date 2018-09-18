package ramen

var emptyCell = Cell{
	Foreground: NewColor(255, 255, 255),
}

// Cell represents a cell in the console
type Cell struct {
	Foreground Color
	Background Color

	Char int
}
