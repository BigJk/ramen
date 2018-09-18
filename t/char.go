package t

import "github.com/BigJk/ramen"

// CharTransform sets the char value of a cell
type CharTransform struct {
	char int
}

// Transform sets the char value of a cell
func (c CharTransform) Transform(cell *ramen.Cell) (bool, error) {
	if cell.Char == c.char {
		return false, nil
	}
	cell.Char = c.char
	return true, nil
}

// Char creates a transformer that sets the cells char value to the given int
func Char(newValue int) CharTransform {
	return CharTransform{newValue}
}

// CharRune creates a transformer that sets the cells char value to the given rune
func CharRune(newValue rune) CharTransform {
	return Char(int(newValue))
}

// CharByte creates a transformer that sets the cells char value to the given byte
func CharByte(newValue byte) CharTransform {
	return Char(int(newValue))
}
