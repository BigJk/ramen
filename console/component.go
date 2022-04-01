package console

import (
	"fmt"
	"math/rand"
	"sync"
)

// ComponentAttributes represents a closable object with a position and size.
type ComponentAttributes interface {
	ID() string
	Position() (int, int)
	Size() (int, int)
	ShouldClose() bool
	ShouldDraw() bool
	IsFocused() bool
	SetFocus(value bool)
}

// ComponentLogic represents a object that can be updated and drawn on a console.
type ComponentLogic interface {
	Update(con *Console, timeElapsed float64) bool
	Draw(con *Console, timeElapsed float64)
	FocusOnClick() bool
}

// ComponentBase represents the base for a ui element on the console.
type ComponentBase struct {
	X      int
	Y      int
	Width  int
	Height int

	mtx   sync.Mutex
	id    string
	show  bool
	close bool
	focus bool
}

func (cb *ComponentBase) ID() string {
	return cb.id
}

// Position returns the position of the component.
func (cb *ComponentBase) Position() (int, int) {
	return cb.X, cb.Y
}

// Size returns the size of the component.
func (cb *ComponentBase) Size() (int, int) {
	return cb.Width, cb.Height
}

// ShouldClose returns true if the component should be closed and deleted from the console.
func (cb *ComponentBase) ShouldClose() bool {
	cb.mtx.Lock()
	defer cb.mtx.Unlock()
	return cb.close
}

// ShouldDraw returns true if the component should be drawn.
func (cb *ComponentBase) ShouldDraw() bool {
	cb.mtx.Lock()
	defer cb.mtx.Unlock()
	return cb.show
}

// Close tells the component to close and remove itself from the parents component list on the next update.
func (cb *ComponentBase) Close() {
	cb.mtx.Lock()
	defer cb.mtx.Unlock()
	cb.close = true
}

// Show shows or hides the component.
func (cb *ComponentBase) Show(value bool) {
	cb.mtx.Lock()
	defer cb.mtx.Unlock()
	cb.show = value
}

// IsFocused returns true if the component is active, which means it was clicked on.
func (cb *ComponentBase) IsFocused() bool {
	cb.mtx.Lock()
	defer cb.mtx.Unlock()
	return cb.focus
}

// SetFocus adds or remove focus from component.
func (cb *ComponentBase) SetFocus(value bool) {
	cb.mtx.Lock()
	defer cb.mtx.Unlock()
	cb.focus = value
}

// NewComponentBase creates a new component base for ease of use.
func NewComponentBase(x, y, width, height int) *ComponentBase {
	return &ComponentBase{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
		id:     fmt.Sprint(rand.Int63()),
		show:   true,
	}
}

// Component represents a ui element on the console.
type Component interface {
	ComponentAttributes
	ComponentLogic
}
