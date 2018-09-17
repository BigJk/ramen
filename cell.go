package ramen

type ColorType int

const (
	// ForegroundColor defines that the foreground color should be addressed
	ForegroundColor = ColorType(0)
	// BackgroundColor defines that the background color should be addressed
	BackgroundColor = ColorType(1)
)

var emptyCell = Cell{}

// Cell represents a cell in the console
type Cell struct {
	Foreground Color
	Background Color

	Char int
}
