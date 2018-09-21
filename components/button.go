package components

import (
	"github.com/BigJk/ramen/console"
	"github.com/BigJk/ramen/consolecolor"
	"github.com/BigJk/ramen/t"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// ClickedCallback will be called if a click on the component happened.
type ClickedCallback func()

// Button represents a button that you can click.
type Button struct {
	*console.ComponentBase

	text            string
	clickedCallback ClickedCallback
	transformer     []t.Transformer

	background        consolecolor.Color
	backgroundHover   consolecolor.Color
	backgroundClicked consolecolor.Color

	foreground        consolecolor.Color
	foregroundHover   consolecolor.Color
	foregroundClicked consolecolor.Color

	state ComponentState
}

// NewButton creates a new button at the given position, size and text.
func NewButton(x, y, width, height int, text string, callback ClickedCallback) *Button {
	b := Button{
		ComponentBase:     console.NewComponentBase(x, y, width, height),
		text:              text,
		clickedCallback:   callback,
		background:        consolecolor.NewHex("#353a41"),
		backgroundHover:   consolecolor.NewHex("#3a4047"),
		backgroundClicked: consolecolor.NewHex("#2c3036"),
		foreground:        consolecolor.NewHex("#e1e1e1"),
		foregroundHover:   consolecolor.NewHex("#e1e1e1"),
		foregroundClicked: consolecolor.NewHex("#e1e1e1"),
	}

	return &b
}

// Update updates the button
func (b *Button) Update(con *console.Console, timeElapsed float64) bool {
	b.state = CalculateComponentState(con, b.X, b.Y, b.Width, b.Height)

	if b.state == ComponentHovered && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		b.clickedCallback()
	}

	return true
}

// Draw draws the button
func (b *Button) Draw(con *console.Console, timeElapsed float64) {
	tY := b.Y + b.Height/2
	tX := b.X + b.Width/2 - len(b.text)/2

	var bgColor consolecolor.Color
	var fColor consolecolor.Color
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

	con.Clear(b.X, b.Y, b.Width, b.Height, t.Background(bgColor))
	con.Print(tX, tY, b.text, t.Foreground(fColor))

}

// SetBackground sets the background colors for the button states. Parameters that
// are nil will be ignored and not set.
func (b *Button) SetBackground(bg, bgHover, bgClicked *consolecolor.Color) {
	if bg != nil {
		b.background = *bg
	}

	if bgHover != nil {
		b.backgroundHover = *bgHover
	}

	if bgClicked != nil {
		b.backgroundClicked = *bgClicked
	}
}

// SetForeground sets the foreground colors for the button states. Parameters that
// are nil will be ignored and not set.
func (b *Button) SetForeground(f, fHover, fClicked *consolecolor.Color) {
	if f != nil {
		b.foreground = *f
	}

	if fHover != nil {
		b.foregroundHover = *fHover
	}

	if fClicked != nil {
		b.foregroundClicked = *fClicked
	}
}
