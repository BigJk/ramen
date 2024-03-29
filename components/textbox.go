package components

import (
	"sync"

	"github.com/BigJk/ramen/concolor"
	"github.com/BigJk/ramen/console"
	"github.com/BigJk/ramen/t"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TextChangeCallback will be called if a text of the component has been changed.
type TextChangeCallback func(text string)

// EnterCallback will be called if the enter key has been pressed.
type EnterCallback func(text string)

// TextBox represents a box that you can type in.
type TextBox struct {
	*console.ComponentBase

	mtx  sync.RWMutex
	text string

	textChangeCallback TextChangeCallback
	enterCallback      EnterCallback

	background        concolor.Color
	backgroundHover   concolor.Color
	backgroundClicked concolor.Color

	foreground         concolor.Color
	foregroundInactive concolor.Color

	blink float64

	state ComponentState
}

// NewTextbox creates a new textbox at the given position and size.
func NewTextbox(x, y, width, height int) *TextBox {
	tb := TextBox{
		ComponentBase:      console.NewComponentBase(x, y, width, height),
		background:         colorBg,
		backgroundHover:    colorBgHover,
		backgroundClicked:  colorBgClicked,
		foreground:         colorFg,
		foregroundInactive: colorFgInactive,
	}

	return &tb
}

// FocusOnClick returns true if a click should focus the textbox
func (tb *TextBox) FocusOnClick() bool {
	return true
}

// Update updates the textbox.
func (tb *TextBox) Update(con *console.Console, timeElapsed float64) bool {
	tb.state = CalculateComponentState(con, tb.X, tb.Y, tb.Width, tb.Height)

	if !tb.IsFocused() {
		return true
	}

	textChanged := false

	tb.mtx.Lock()
	if len(ebiten.InputChars()) > 0 {
		textChanged = true
	}
	tb.text += string(ebiten.InputChars())

	if tb.repeatingKeyPressed(ebiten.KeyBackspace) {
		if len(tb.text) >= 1 {
			tb.text = tb.text[:len(tb.text)-1]
			textChanged = true
		}
	}

	if textChanged && tb.textChangeCallback != nil {
		tb.textChangeCallback(tb.text)
	}
	tb.mtx.Unlock()

	if tb.enterCallback != nil && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		tb.enterCallback(tb.text)
	}

	return true
}

// Draw draws the textbox.
func (tb *TextBox) Draw(con *console.Console, timeElapsed float64) {
	// TODO: multi-line support

	tb.blink += timeElapsed

	tb.mtx.RLock()
	text := tb.text
	tb.mtx.RUnlock()

	var bgColor concolor.Color
	var fColor concolor.Color
	switch tb.state {
	case ComponentIdle:
		bgColor = tb.background
		if tb.IsFocused() {
			fColor = tb.foreground
		} else {
			fColor = tb.foregroundInactive
		}
	case ComponentHovered:
		bgColor = tb.backgroundHover
		fColor = tb.foreground
	case ComponentClicked:
		bgColor = tb.backgroundClicked
		fColor = tb.foreground
	}

	if tb.IsFocused() {
		bgColor = tb.backgroundHover
		fColor = tb.foreground
	}

	con.TransformArea(tb.X, tb.Y, tb.Width, tb.Height, t.Background(bgColor))

	if tb.blink < 0.5 && tb.IsFocused() {
		text += "_"
	} else {
		text += " "
	}

	if tb.Height == 1 && len(text) >= tb.Width {
		text = text[len(text)-tb.Width:]
	}

	con.Print(tb.X, tb.Y, text, t.Foreground(fColor))

	if tb.blink > 1 {
		tb.blink = 0
	}
}

// SetTextChanged sets the text change callback.
func (tb *TextBox) SetTextChanged(callback TextChangeCallback) {
	tb.textChangeCallback = callback
}

// SetEnterCallback sets the enter callback.
func (tb *TextBox) SetEnterCallback(callback EnterCallback) {
	tb.enterCallback = callback
}

// SetText changes the text of the textbox.
func (tb *TextBox) SetText(newText string) {
	tb.mtx.Lock()
	tb.text = newText
	tb.mtx.Unlock()
}

// GetText returns the text of the textbox.
func (tb *TextBox) GetText() string {
	tb.mtx.RLock()
	defer tb.mtx.RUnlock()
	return tb.text
}

// SetBackground sets the background colors for textbox. Parameters that
// are nil will be ignored and not set.
func (tb *TextBox) SetBackground(idle, hover, clicked *concolor.Color) {
	if idle != nil {
		tb.background = *idle
	}

	if hover != nil {
		tb.backgroundHover = *hover
	}

	if clicked != nil {
		tb.backgroundClicked = *clicked
	}
}

// SetForeground sets the foreground colors for the textbox states. Parameters that
// are nil will be ignored and not set.
func (tb *TextBox) SetForeground(active, inactive *concolor.Color) {
	if active != nil {
		tb.foreground = *active
	}

	if inactive != nil {
		tb.foregroundInactive = *inactive
	}
}

func (tb *TextBox) repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}
