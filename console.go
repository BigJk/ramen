package ramen

import (
	"fmt"
	"sync"

	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Console represents a emulated console view
type Console struct {
	Title   string
	Width   int
	Height  int
	Font    *Font
	ShowFPS bool

	mtx     *sync.RWMutex
	updates []int
	buffer  [][]Cell

	lines []*ebiten.Image

	preRenderHook  func(screen *ebiten.Image) error
	postRenderHook func(screen *ebiten.Image) error
}

// NewConsole creates a new console
func NewConsole(width, height int, font *Font, title string) (*Console, error) {
	buf := make([][]Cell, width)
	for i := range buf {
		buf[i] = make([]Cell, height)
	}

	lines := make([]*ebiten.Image, width)
	for i := range lines {
		line, err := ebiten.NewImage(font.TileWidth, height*font.TileHeight, ebiten.FilterNearest)
		if err != nil {
			return nil, err
		}
		lines[i] = line
	}

	return &Console{
		Title:   title,
		Width:   width,
		Height:  height,
		Font:    font,
		mtx:     new(sync.RWMutex),
		updates: make([]int, 0),
		buffer:  buf,
		lines:   lines,
	}, nil
}

// SetPreRenderHook will apply a hook that gets triggered before the console started rendering
func (c *Console) SetPreRenderHook(hook func(screen *ebiten.Image) error) {
	c.preRenderHook = hook
}

// SetPostRenderHook will apply a hook that gets triggered after the console is finished rendering
func (c *Console) SetPostRenderHook(hook func(screen *ebiten.Image) error) {
	c.preRenderHook = hook
}

// Start will open the console window with the given scale
func (c *Console) Start(scale float64) error {
	return ebiten.Run(c.update, c.Width*c.Font.TileWidth, c.Height*c.Font.TileHeight, scale, c.Title)
}

// PutCharRune will set the char value of the cell at the given position to the value of the rune
func (c *Console) PutCharRune(x, y int, char rune) error {
	return c.PutCharInt(x, y, int(char))
}

// PutCharByte will set the char value of the cell at the given position to the value of the byte
func (c *Console) PutCharByte(x, y int, char byte) error {
	return c.PutCharInt(x, y, int(char))
}

// PutCharInt will set the char value of the cell at the given position to the value of the int
func (c *Console) PutCharInt(x, y int, char int) error {
	if err := c.checkOutOfBounds(x, y); err != nil {
		return err
	}

	c.mtx.Lock()
	if c.buffer[x][y].Char == char {
		c.mtx.Unlock()
		return nil
	}
	c.buffer[x][y].Char = char
	c.queueUpdate(x)
	c.mtx.Unlock()

	return nil
}

// PutColor will either change the foreground or background color of the tile at the given position
func (c *Console) PutColor(x, y int, color Color, colorType ColorType) error {
	if err := c.checkOutOfBounds(x, y); err != nil {
		return err
	}

	c.mtx.Lock()
	if colorType == ForegroundColor {
		if c.buffer[x][y].Foreground == color {
			c.mtx.Unlock()
			return nil
		}
		c.buffer[x][y].Foreground = color
	} else if colorType == BackgroundColor {
		if c.buffer[x][y].Background == color {
			c.mtx.Unlock()
			return nil
		}
		c.buffer[x][y].Background = color
	}
	c.queueUpdate(x)
	c.mtx.Unlock()

	return nil
}

// GetCell returns the cell data of the given position
func (c *Console) GetCell(x, y int) (Cell, error) {
	if err := c.checkOutOfBounds(x, y); err != nil {
		return Cell{}, err
	}

	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.buffer[x][y], nil
}

// Print prints a string on the console with a given foreground color
func (c *Console) Print(x, y int, text string, foreground Color) {
	for i := range text {
		c.PutCharByte(x+i, y, text[i])
		c.PutColor(x+i, y, foreground, ForegroundColor)
	}
}

// PrintEx prints a string on the console with a given foreground and background color
func (c *Console) PrintEx(x, y int, text string, foreground Color, background Color) {
	for i := range text {
		c.PutColor(x+i, y, background, BackgroundColor)
	}
	c.Print(x, y, text, foreground)
}

// PrintFmt prints a formatted string on the console with a given foreground color
func (c *Console) PrintFmt(x, y int, format string, foreground Color, a ...interface{}) {
	text := fmt.Sprintf(format, a...)
	for i := range text {
		c.PutCharByte(x+i, y, text[i])
		c.PutColor(x+i, y, foreground, ForegroundColor)
	}
}

// PrintFmt prints a formatted string on the console with a given foreground and background color
func (c *Console) PrintFmtEx(x, y int, format string, foreground Color, background Color, a ...interface{}) {
	text := fmt.Sprintf(format, a...)
	for i := range text {
		c.PutColor(x+i, y, background, BackgroundColor)
	}
	c.Print(x, y, text, foreground)
}

// PrintFrame prints a frame on the console
func (c *Console) PrintFrame(x, y, width, height int, frame Frame) {
	c.PutCharInt(x, y, frame.TopLeft)
	c.PutColor(x, y, frame.Foreground, ForegroundColor)
	c.PutColor(x, y, frame.Background, BackgroundColor)

	c.PutCharInt(x+width-1, y, frame.TopRight)
	c.PutColor(x+width-1, y, frame.Foreground, ForegroundColor)
	c.PutColor(x+width-1, y, frame.Background, BackgroundColor)

	c.PutCharInt(x, y+height-1, frame.BottomLeft)
	c.PutColor(x, y+height-1, frame.Foreground, ForegroundColor)
	c.PutColor(x, y+height-1, frame.Background, BackgroundColor)

	c.PutCharInt(x+width-1, y+height-1, frame.BottomRight)
	c.PutColor(x+width-1, y+height-1, frame.Foreground, ForegroundColor)
	c.PutColor(x+width-1, y+height-1, frame.Background, BackgroundColor)

	for i := 1; i < width-1; i++ {
		c.PutCharInt(x+i, y, frame.Horizontal)
		c.PutColor(x+i, y, frame.Foreground, ForegroundColor)
		c.PutColor(x+i, y, frame.Background, BackgroundColor)

		c.PutCharInt(x+i, y+height-1, frame.Horizontal)
		c.PutColor(x+i, y+height-1, frame.Foreground, ForegroundColor)
		c.PutColor(x+i, y+height-1, frame.Background, BackgroundColor)
	}

	for i := 1; i < height-1; i++ {
		c.PutCharInt(x, y+i, frame.Vertical)
		c.PutColor(x, y+i, frame.Foreground, ForegroundColor)
		c.PutColor(x, y+i, frame.Background, BackgroundColor)

		c.PutCharInt(x+width-1, y+i, frame.Vertical)
		c.PutColor(x+width-1, y+i, frame.Foreground, ForegroundColor)
		c.PutColor(x+width-1, y+i, frame.Background, BackgroundColor)
	}
}

// PrintFrameEx prints a frame with a title on the console
func (c *Console) PrintFrameEx(x, y, width, height int, frame Frame, title string) {
	c.PrintFrame(x, y, width, height, frame)
	c.Print(x+5, y, "["+title+"]", frame.Foreground)
}

// Clear clears the whole console
func (c *Console) ClearAll() {
	c.mtx.Lock()
	for x := range c.buffer {
		for y := range c.buffer[x] {
			c.buffer[x][y] = emptyCell
			c.updates = append(c.updates, x)
		}
	}
	c.mtx.Unlock()
}

// Clear clears part of the console
func (c *Console) Clear(x, y, width, height int) {
	c.mtx.Lock()
	for px := 0; px < width; px++ {
		mustUpdate := false
		for py := 0; py < height; py++ {
			if c.buffer[px+x][py+y] != emptyCell {
				c.buffer[px+x][py+y] = emptyCell
				mustUpdate = true
			}
		}
		if mustUpdate {
			c.updates = append(c.updates, px+x)
		}
	}
	c.mtx.Unlock()
}

func (c *Console) queueUpdate(x int) {
	for i := range c.updates {
		if c.updates[i] == x {
			return
		}
	}
	c.updates = append(c.updates, x)
}

func (c *Console) checkOutOfBounds(x, y int) error {
	if x < 0 || y < 0 || x >= c.Width || y >= c.Height {
		return fmt.Errorf("position out of bounds")
	}
	return nil
}

func (c *Console) updateLine(x int) {
	c.lines[x].Fill(color.NRGBA{0, 0, 0, 0})
	for y := range c.buffer[x] {
		if c.buffer[x][y].Background.A > 0 {
			ebitenutil.DrawRect(c.lines[x], 0, float64(y*c.Font.TileHeight), float64(c.Font.TileWidth), float64(c.Font.TileHeight), c.buffer[x][y].Background)
		}

		if c.buffer[x][y].Char == 0 {
			continue
		}

		op := c.Font.ToOptions(c.buffer[x][y].Char)
		op.GeoM.Translate(0, float64(y*c.Font.TileHeight))

		if !c.Font.IsTile(c.buffer[x][y].Char) {
			op.ColorM.Scale(c.buffer[x][y].Foreground.Floats())
		}

		c.lines[x].DrawImage(c.Font.Image, op)
	}
}

func (c *Console) update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	if c.preRenderHook != nil {
		if err := c.preRenderHook(screen); err != nil {
			return err
		}
	}

	c.mtx.RLock()

	for i := range c.updates {
		c.updateLine(c.updates[i])
	}

	if len(c.updates) > 0 {
		c.updates = make([]int, 0)
	}

	for x := range c.buffer {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x*c.Font.TileWidth), 0)
		screen.DrawImage(c.lines[x], op)
	}

	c.mtx.RUnlock()

	if c.postRenderHook != nil {
		if err := c.postRenderHook(screen); err != nil {
			return err
		}
	}

	if c.ShowFPS {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
	}

	return nil
}
