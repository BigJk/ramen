package main

import (
	"flag"

	"os"

	"github.com/BigJk/ramen/console"
	"github.com/BigJk/ramen/font"
	"github.com/BigJk/ramen/t"
	"github.com/BigJk/ramen/xp"
)

func main() {
	filePath := flag.String("file", "", "the REXPaint file to show")
	flag.Parse()

	if len(*filePath) == 0 {
		flag.PrintDefaults()
		return
	}

	file, err := os.Open(*filePath)
	if err != nil {
		panic(err)
	}

	xp, err := xp.Read(file)
	if err != nil {
		panic(err)
	}

	con, err := console.New(xp.Width, xp.Height, font.DefaultFont, "ramen - REXPaint Viewer")
	if err != nil {
		panic(err)
	}

	layers := []*console.Console{con}
	for len(layers) < len(xp.Layers) {
		newCon, err := layers[len(layers)-1].CreateSubConsole(0, 0, layers[len(layers)-1].Width, layers[len(layers)-1].Height)
		if err != nil {
			panic(err)
		}
		layers = append(layers, newCon)
	}

	for l := range xp.Layers {
		for x := range xp.Layers[l].Cells {
			for y := range xp.Layers[l].Cells[x] {
				layers[l].Transform(x, y, t.Char(xp.Layers[l].Cells[x][y].Char), t.Foreground(xp.Layers[l].Cells[x][y].Foreground), t.Background(xp.Layers[l].Cells[x][y].Background))
			}
		}
	}

	con.Start(1)
}
