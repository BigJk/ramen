package ramen

// ColorModifier represents a identifier to differentiate between a foreground and background color change
type ColorModifier int

const (
	// ModifyForegroundColor defines that the foreground color should be changed
	ModifyForegroundColor = ColorModifier(1)
	// ModifyBackgroundColor defines that the background color should be changed
	ModifyBackgroundColor = ColorModifier(2)
)

func (c ColorModifier) HasFlag(flag ColorModifier) bool {
	return c|flag == c
}
