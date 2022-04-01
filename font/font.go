// Package font provides functionality to load ascii fonts from png files.
package font

import (
	"image"

	"io"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
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

	return &Font{"", ebiten.NewImageFromImage(img), tileWidth, tileHeight, img.Bounds().Max.X / tileWidth, img.Bounds().Max.Y / tileHeight, make(map[int]bool)}, nil
}

// ToSubImage extracts the image of a given char from the base image of the font.
func (f *Font) ToSubImage(char int) *ebiten.Image {
	x := (int(char) % f.TileSizeX) * f.TileWidth
	y := (int(char) / f.TileSizeY) * f.TileHeight

	r := image.Rect(x, y, x+f.TileWidth, y+f.TileHeight)

	if r.Max.X > f.Image.Bounds().Max.X || r.Max.Y > f.Image.Bounds().Max.Y {
		return nil
	}

	return f.Image.SubImage(r).(*ebiten.Image)
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
