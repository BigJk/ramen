// Package font provides functionality to load ascii fonts from png files.
package font

import (
	"image"

	"io"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Font represents a console font.
type Font struct {
	File       string
	Image      *ebiten.Image
	TileWidth  int
	TileHeight int
	TileSizeX  int
	TileSizeY  int
	Tiles      map[int]bool
}

// New creates a new font.
func New(filePath string, tileWidth, tileHeight int) (*Font, error) {
	file, err := ebitenutil.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	f, err := NewFromReader(file, tileWidth, tileHeight)
	if err != nil {
		return nil, err
	}
	f.File = filePath

	return f, nil
}

// NewFromReader creates a new font from a reader.
func NewFromReader(reader io.Reader, tileWidth, tileHeight int) (*Font, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	fontImage, err := ebiten.NewImageFromImage(img, ebiten.FilterNearest)
	if err != nil {
		return nil, err
	}

	return &Font{"", fontImage, tileWidth, tileHeight, img.Bounds().Max.X / tileWidth, img.Bounds().Max.Y / tileHeight, make(map[int]bool)}, nil
}

// ToOptions extracts the rectangle of a given char from the base image of the font.
func (f *Font) ToOptions(char int) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}

	x := (int(char) % f.TileSizeX) * f.TileWidth
	y := (int(char) / f.TileSizeY) * f.TileHeight

	r := image.Rect(x, y, x+f.TileWidth, y+f.TileHeight)
	op.SourceRect = &r

	return op
}

// SetTiles changes if a char is a colored tile or not.
// start specifies the char where SetTiles should start setting the value
// and count is the amount of chars after start that should also be set.
func (f *Font) SetTiles(start, count int, value bool) {
	for i := start; i <= start+count; i++ {
		f.Tiles[i] = value
	}
}

// IsTile checks if the given char represents a colored tile
func (f *Font) IsTile(char int) bool {
	val, ok := f.Tiles[char]
	return ok && val
}
