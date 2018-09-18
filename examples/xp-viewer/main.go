package main

import (
	"flag"

	"os"

	"github.com/BigJk/ramen"
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

	font, err := ramen.NewFont("../../fonts/terminus-11x11.png", 11, 11)
	if err != nil {
		panic(err)
	}

	con, err := ramen.NewConsole(xp.Width, xp.Height, font, "REXPaint Viewer")
	if err != nil {
		panic(err)
	}

	layers := []*ramen.Console{con}
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
				layers[l].PutCharInt(x, y, xp.Layers[l].Cells[x][y].Char)
				layers[l].PutColor(x, y, xp.Layers[l].Cells[x][y].Foreground, ramen.ModifyForegroundColor)
				layers[l].PutColor(x, y, xp.Layers[l].Cells[x][y].Background, ramen.ModifyBackgroundColor)
			}
		}
	}

	con.Start(1)
}
