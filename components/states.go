package components

import (
	"github.com/BigJk/ramen/console"
	"github.com/hajimehoshi/ebiten"
)

// ComponentState represents generic states that components can use to define logic.
type ComponentState int

const (
	// ComponentIdle means the component is neither clicked not hovered over.
	ComponentIdle = ComponentState(0)
	// ComponentHovered means the component is hovered over
	ComponentHovered = ComponentState(1)
	// ComponentClicked means the component is hovered over and the left mouse button is currently pressed.
	ComponentClicked = ComponentState(2)
)

// CalculateComponentState is generic helper function to calculate the state of a given area in the console.
func CalculateComponentState(con *console.Console, x, y, w, h int) ComponentState {
	if con.MouseInArea(x, y, w, h) {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			return ComponentClicked
		} else {
			return ComponentHovered
		}
	}
	return ComponentIdle
}
