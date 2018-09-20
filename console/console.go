package console

import (
	"fmt"
	"math"
	"sync"

	"sort"

	"github.com/BigJk/ramen"
	"github.com/BigJk/ramen/consolecolor"
	"github.com/BigJk/ramen/font"
	"github.com/BigJk/ramen/t"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var emptyCell = ramen.Cell{
	Foreground: consolecolor.New(255, 255, 255),
}

// Console represents a emulated console view
type Console struct {
	Title       string
	Width       int
	Height      int
	Font        *font.Font
	ShowFPS     bool
	SubConsoles []*Console

	parent       *Console
	x            int
	y            int
	priority     int
	isSubConsole bool

	mtx    *sync.RWMutex
	buffer [][]ramen.Cell

	tickHook       func(timeElapsed float64) error
	preRenderHook  func(screen *ebiten.Image, timeElapsed float64) error
	postRenderHook func(screen *ebiten.Image, timeElapsed float64) error
}

// New creates a new console
func New(width, height int, font *font.Font, title string) (*Console, error) {
	buf := make([][]ramen.Cell, width)
	for x := range buf {
		buf[x] = make([]ramen.Cell, height)
		for y := range buf[x] {
			buf[x][y] = emptyCell
		}
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
		Title:       title,
		Width:       width,
		Height:      height,
		Font:        font,
		SubConsoles: make([]*Console, 0),
		mtx:         new(sync.RWMutex),
		buffer:      buf,
	}, nil
}

// SetTickHook will apply a hook that gets triggered every tick, even if drawing is skipped in this tick.
// This is a good place for game logic as it runs disconnected from the fps.
func (c *Console) SetTickHook(hook func(timeElapsed float64) error) error {
	if c.isSubConsole {
		return fmt.Errorf("can't hook into sub-console")
	}
	c.tickHook = hook
	return nil
}

// SetPreRenderHook will apply a hook that gets triggered before the console started rendering.
// This is a good place to change the console or to draw extra content under the console.
func (c *Console) SetPreRenderHook(hook func(screen *ebiten.Image, timeElapsed float64) error) error {
	if c.isSubConsole {
		return fmt.Errorf("can't hook into sub-console")
	}
	c.preRenderHook = hook
	return nil
}

// SetPostRenderHook will apply a hook that gets triggered after the console is finished rendering.
// This is a good place if you want to draw some extra content over the console.
func (c *Console) SetPostRenderHook(hook func(screen *ebiten.Image, timeElapsed float64) error) error {
	if c.isSubConsole {
		return fmt.Errorf("can't hook into sub-console")
	}
	c.postRenderHook = hook
	return nil
}

// SetPriority sets the priority of the console. A higher priority will result in the console
// being drawn on top of all the ones with lower priority.
func (c *Console) SetPriority(priority int) error {
	if !c.isSubConsole {
		return fmt.Errorf("priority of the main console can't be changed")
	}
	c.priority = priority
	c.parent.sortSubConsoles()
	return nil
}

// CreateSubConsole creates a new sub-console
func (c *Console) CreateSubConsole(x, y, width, height int) (*Console, error) {
	if x < 0 || y < 0 || x+width > c.Width || y+height > c.Height || width <= 0 || height <= 0 {
		return nil, fmt.Errorf("sub-console is out of bounds")
	}

	c.mtx.Lock()

	sub, err := New(width, height, c.Font, "")
	if err != nil {
		return nil, err
	}

	sub.parent = c
	sub.x = x
	sub.y = y
	sub.isSubConsole = true

	c.SubConsoles = append(c.SubConsoles, sub)

	c.mtx.Unlock()

	c.sortSubConsoles()

	return sub, nil
}

// RemoveSubConsole removes a sub-console from his parent
func (c *Console) RemoveSubConsole(con *Console) error {
	c.mtx.Lock()
	for i := range c.SubConsoles {
		if c.SubConsoles[i] == con {
			c.SubConsoles[i] = c.SubConsoles[len(c.SubConsoles)-1]
			c.SubConsoles[len(c.SubConsoles)-1] = nil
			c.SubConsoles = c.SubConsoles[:len(c.SubConsoles)-1]
			c.mtx.Unlock()

			c.sortSubConsoles()

			return nil
		}
	}
	c.mtx.Unlock()
	return fmt.Errorf("sub-console not found")
}

// Start will open the console window with the given scale
func (c *Console) Start(scale float64) error {
	if c.isSubConsole {
		return fmt.Errorf("only the main console can be started")
	}
	return ebiten.Run(c.update, c.Width*c.Font.TileWidth, c.Height*c.Font.TileHeight, scale, c.Title)
}

// ClearAll clears the whole console. If no transformer are specified the console will be cleared
// to the default cell look.
func (c *Console) ClearAll(transformer ...t.Transformer) {
	c.Clear(0, 0, c.Width, c.Height, transformer...)
}

// Clear clears part of the console. If no transformer are specified the console will be cleared
// to the default cell look.
func (c *Console) Clear(x, y, width, height int, transformer ...t.Transformer) error {
	c.mtx.Lock()

	for px := 0; px < width; px++ {
		for py := 0; py < height; py++ {
			if err := c.checkOutOfBounds(px+x, py+y); err != nil {
				return err
			}

			if len(transformer) == 0 {
				if c.buffer[px+x][py+y] != emptyCell {
					c.buffer[px+x][py+y] = emptyCell
				}
			} else {
				for i := range transformer {
					err := transformer[i].Transform(&c.buffer[px+x][py+y])
					if err != nil {
						return err
					}
				}
			}
		}
	}

	c.mtx.Unlock()

	return nil
}

// Transform transforms a cell. This can be used to change the character, foreground and
// background of a cell or apply custom transformers onto a cell.
func (c *Console) Transform(x, y int, transformer ...t.Transformer) error {
	if len(transformer) == 0 {
		return fmt.Errorf("no transformer given")
	} else if err := c.checkOutOfBounds(x, y); err != nil {
		return err
	}

	c.mtx.Lock()

	for i := range transformer {
		err := transformer[i].Transform(&c.buffer[x][y])
		if err != nil {
			return err
		}
	}

	c.mtx.Unlock()

	return nil
}

// Print prints a text onto the console. To give the text a different foreground or
// background color use transformer.
func (c *Console) Print(x, y int, text string, transformer ...t.Transformer) {
	if y >= c.Height {
		return
	}

	linePos := 0
	for i := range text {
		if x+i >= c.Width {
			continue
		}

		if text[i] == '\n' {
			y++
			linePos = 0
			continue
		}

		c.Transform(linePos+x, y, append(transformer, t.CharByte(text[i]))...)
		linePos++
	}
}

func (c *Console) sortSubConsoles() {
	c.mtx.Lock()
	sort.Slice(c.SubConsoles, func(i, j int) bool {
		return c.SubConsoles[i].priority < c.SubConsoles[j].priority
	})
	c.mtx.Unlock()
}

func (c *Console) checkOutOfBounds(x, y int) error {
	if x < 0 || y < 0 || x >= c.Width || y >= c.Height {
		return fmt.Errorf("position out of bounds")
	}
	return nil
}

func (c *Console) draw(screen *ebiten.Image, offsetX, offsetY int) {
	for x := range c.buffer {
		for y := range c.buffer[x] {
			if c.buffer[x][y].Background.A == 0 {
				continue
			}

			ebitenutil.DrawRect(screen, float64((offsetX+c.x+x)*c.Font.TileWidth), float64((offsetY+c.y+y)*c.Font.TileHeight), float64(c.Font.TileWidth), float64(c.Font.TileHeight), c.buffer[x][y].Background)
		}
	}

	for x := range c.buffer {
		for y := range c.buffer[x] {
			op := c.Font.ToOptions(c.buffer[x][y].Char)
			if !c.Font.IsTile(c.buffer[x][y].Char) {
				op.ColorM.Scale(c.buffer[x][y].Foreground.Floats())
			}
			op.GeoM.Translate(float64((offsetX+c.x+x)*c.Font.TileWidth), float64((offsetY+c.y+y)*c.Font.TileHeight))
			screen.DrawImage(c.Font.Image, op)
		}
	}

	for i := range c.SubConsoles {
		c.SubConsoles[i].draw(screen, offsetX+c.x, offsetY+c.y)
	}
}

func (c *Console) elapsedTPS() float64 {
	e := 1.0 / math.Min(float64(ebiten.MaxTPS()), ebiten.CurrentTPS())
	if e > math.MaxFloat64 {
		e = 0
	}
	return e
}

func (c *Console) elapsedFPS() float64 {
	e := 1.0 / math.Min(float64(ebiten.FPS), ebiten.CurrentFPS())
	if e > math.MaxFloat64 {
		e = 0
	}
	return e
}

func (c *Console) update(screen *ebiten.Image) error {
	if c.tickHook != nil {
		if err := c.tickHook(c.elapsedTPS()); err != nil {
			return err
		}
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	timeElapsed := c.elapsedFPS()

	if c.preRenderHook != nil {
		if err := c.preRenderHook(screen, timeElapsed); err != nil {
			return err
		}
	}

	c.mtx.RLock()
	c.draw(screen, 0, 0)
	c.mtx.RUnlock()

	if c.postRenderHook != nil {
		if err := c.postRenderHook(screen, timeElapsed); err != nil {
			return err
		}
	}

	if c.ShowFPS {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
	}

	return nil
}
