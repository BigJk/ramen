package components

import (
	"github.com/BigJk/ramen/concolor"
	"github.com/BigJk/ramen/console"
	"github.com/BigJk/ramen/t"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// ClickedCallback will be called if a click on the component happened.
type ClickedCallback func()

// Button represents a button that you can click.
type Button struct {
	*console.ComponentBase

	text            string
	clickedCallback ClickedCallback
	transformer     []t.Transformer

	background        concolor.Color
	backgroundHover   concolor.Color
	backgroundClicked concolor.Color

	foreground        concolor.Color
	foregroundHover   concolor.Color
	foregroundClicked concolor.Color

	state ComponentState
}

// NewButton creates a new button at the given position, size and text.
func NewButton(x, y, width, height int, text string, callback ClickedCallback) *Button {
	b := Button{
		ComponentBase:     console.NewComponentBase(x, y, width, height),
		text:              text,
		clickedCallback:   callback,
		background:        colorBg,
		backgroundHover:   colorBgHover,
		backgroundClicked: colorBgClicked,
		foreground:        colorFg,
		foregroundHover:   colorFgHover,
		foregroundClicked: colorFgClicked,
	}

	return &b
}

// FocusOnClick returns true if a click should focus the button
func (b *Button) FocusOnClick() bool {
	return false
}

// Update updates the button
func (b *Button) Update(con *console.Console, timeElapsed float64) bool {
	b.state = CalculateComponentState(con, b.X, b.Y, b.Width, b.Height)

	if b.state == ComponentHovered && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		b.clickedCallback()
	}

	if b.IsFocused() && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		b.clickedCallback()
	}

	return true
}

// Draw draws the button
func (b *Button) Draw(con *console.Console, timeElapsed float64) {
	tY := b.Y + b.Height/2
	tX := b.X + b.Width/2 - len(b.text)/2

	var bgColor concolor.Color
	var fColor concolor.Color
	switch b.state {
	case ComponentIdle:
		bgColor = b.background
		fColor = b.foreground
	case ComponentHovered:
		bgColor = b.backgroundHover
		fColor = b.foregroundHover
	case ComponentClicked:
		bgColor = b.backgroundClicked
		fColor = b.foregroundClicked
	}

	if b.IsFocused() {
		bgColor = b.backgroundHover
		fColor = b.foregroundHover
	}

	_ = con.Clear(b.X, b.Y, b.Width, b.Height, t.Background(bgColor))
	con.Print(tX, tY, b.text, t.Foreground(fColor))
}

// SetBackground sets the background colors for the button states. Parameters that
// are nil will be ignored and not set.
func (b *Button) SetBackground(idle, hover, clicked *concolor.Color) {
	if idle != nil {
		b.background = *idle
	}

	if hover != nil {
		b.backgroundHover = *hover
	}

	if clicked != nil {
		b.backgroundClicked = *clicked
	}
}

// SetForeground sets the foreground colors for the button states. Parameters that
// are nil will be ignored and not set.
func (b *Button) SetForeground(idle, hover, clicked *concolor.Color) {
	if idle != nil {
		b.foreground = *idle
	}

	if hover != nil {
		b.foregroundHover = *hover
	}

	if clicked != nil {
		b.foregroundClicked = *clicked
	}
}
