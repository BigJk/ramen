package components

import (
	"sync"

	"github.com/BigJk/ramen/console"
	"github.com/BigJk/ramen/consolecolor"
	"github.com/BigJk/ramen/t"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// TextChangeCallback will be called if a text of the component has been changed.
type TextChangeCallback func(text string)

// EnterCallback will be called if the enter key has been pressed.
type EnterCallback func()

// TextBox represents a box that you can type in.
type TextBox struct {
	*console.ComponentBase

	mtx  sync.RWMutex
	text string

	textChangeCallback TextChangeCallback
	enterCallback      EnterCallback

	background        consolecolor.Color
	backgroundHover   consolecolor.Color
	backgroundClicked consolecolor.Color

	foreground         consolecolor.Color
	foregroundInactive consolecolor.Color

	blink float64

	state ComponentState
}

// NewTextbox creates a new textbox at the given position and size.
func NewTextbox(x, y, width, height int) *TextBox {
	tb := TextBox{
		ComponentBase:      console.NewComponentBase(x, y, width, height),
		background:         consolecolor.NewHex("#353a41"),
		backgroundHover:    consolecolor.NewHex("#3a4047"),
		backgroundClicked:  consolecolor.NewHex("#2c3036"),
		foreground:         consolecolor.NewHex("#e1e1e1"),
		foregroundInactive: consolecolor.NewHex("#949494"),
	}

	return &tb
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
		tb.enterCallback()
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

	var bgColor consolecolor.Color
	var fColor consolecolor.Color
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

	con.Clear(tb.X, tb.Y, tb.Width, tb.Height, t.Background(bgColor))

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
func (tb *TextBox) SetBackground(bg, bgHover, bgClicked *consolecolor.Color) {
	if bg != nil {
		tb.background = *bg
	}

	if bgHover != nil {
		tb.backgroundHover = *bgHover
	}

	if bgClicked != nil {
		tb.backgroundClicked = *bgClicked
	}
}

// SetForeground sets the foreground colors for the textbox states. Parameters that
// are nil will be ignored and not set.
func (tb *TextBox) SetForeground(f, fInactive *consolecolor.Color) {
	if f != nil {
		tb.foreground = *f
	}

	if fInactive != nil {
		tb.foregroundInactive = *fInactive
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
