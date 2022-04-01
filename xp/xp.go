// Package xp provides REXPaint (.xp) file parsing.
package xp

import (
	"compress/gzip"
	"io"

	"encoding/binary"

	"fmt"

	"github.com/BigJk/ramen/concolor"
)

// XP represents a REXPaint file
type XP struct {
	Version int
	Width   int
	Height  int
	Layers  []Layer
}

// Cell represents a cell in the REXPaint file
type Cell struct {
	Char       int
	Foreground concolor.Color
	Background concolor.Color
}

// Layer represents a layer of cells in the REXPaint file
type Layer struct {
	Cells [][]Cell
}

// Read parses a xp file from the reader
func Read(reader io.Reader) (*XP, error) {
	deflated, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}

	var xp XP

	var version int32
	if err := binary.Read(deflated, binary.LittleEndian, &version); err != nil {
		return nil, fmt.Errorf("error while reading version")
	}
	xp.Version = int(version) * -1

	var layers int32
	if err := binary.Read(deflated, binary.LittleEndian, &layers); err != nil {
		return nil, fmt.Errorf("error while reading layer count")
	}

	var width, height uint32
	for i := 0; i < int(layers); i++ {
		var newLayer Layer

		if err := binary.Read(deflated, binary.LittleEndian, &width); err != nil {
			return nil, fmt.Errorf("error while reading layer width")
		}

		if err := binary.Read(deflated, binary.LittleEndian, &height); err != nil {
			return nil, fmt.Errorf("error while reading layer height")
		}

		if i == 0 {
			xp.Width = int(width)
			xp.Height = int(height)
		}

		newLayer.Cells = make([][]Cell, xp.Width)
		for x := 0; x < int(xp.Width); x++ {
			newLayer.Cells[x] = make([]Cell, xp.Height)
		}

		for j := 0; j < int(xp.Width*xp.Height); j++ {
			var cell Cell

			var char uint32
			if err := binary.Read(deflated, binary.LittleEndian, &char); err != nil {
				return nil, fmt.Errorf("error while reading char code")
			}
			cell.Char = int(char)

			if err := readRGB(deflated, &cell.Foreground); err != nil {
				return nil, err
			}

			if err := readRGB(deflated, &cell.Background); err != nil {
				return nil, err
			}

			if cell.Background.R == 255 && cell.Background.G == 0 && cell.Background.B == 255 {
				cell.Background.A = 0
				cell.Foreground.A = 0
			}

			x, y := j/int(xp.Height), j%int(xp.Height)
			newLayer.Cells[x][y] = cell
		}

		xp.Layers = append(xp.Layers, newLayer)
	}

	return &xp, nil
}

func readRGB(reader io.Reader, target *concolor.Color) error {
	rgb := make([]byte, 3)
	num, err := io.ReadAtLeast(reader, rgb, 3)

	if num != 3 || err != nil {
		return fmt.Errorf("error while reading rgb value")
	}

	target.R = rgb[0]
	target.G = rgb[1]
	target.B = rgb[2]
	target.A = 255

	return nil
}
