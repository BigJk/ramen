package main

import (
	"image/color"

	"github.com/BigJk/ramen"
	"github.com/hajimehoshi/ebiten"
)

func main() {
	/*
		load a font that contains colored tiles
	*/
	font, err := ramen.NewFont("./custom_font.png", 8, 8)
	if err != nil {
		panic(err)
	}

	/*
		specify which chars in the fonts are colored tiles
	*/
	font.SetTiles(272, 16, true)

	/*
		create a console
	*/
	con, err := ramen.NewConsole(40, 20, font, "ramen - colored tiles")
	if err != nil {
		panic(err)
	}

	/*
		set a pre-render hook
	*/
	con.SetPreRenderHook(func(screen *ebiten.Image, deltaTime float64) error {
		screen.Fill(color.NRGBA{69, 40, 60, 255})

		/*
			use colored tiles to draw a frame
		*/
		con.PutCharInt(0, 0, 272)
		con.PutCharInt(1, 0, 276)
		con.PutCharInt(0, 1, 280)

		con.PutCharInt(con.Width-1, 0, 273)
		con.PutCharInt(con.Width-2, 0, 277)
		con.PutCharInt(con.Width-1, 1, 281)

		con.PutCharInt(0, con.Height-1, 274)
		con.PutCharInt(1, con.Height-1, 278)
		con.PutCharInt(0, con.Height-2, 282)

		con.PutCharInt(con.Width-1, con.Height-1, 275)
		con.PutCharInt(con.Width-2, con.Height-1, 279)
		con.PutCharInt(con.Width-1, con.Height-2, 283)

		/*
			draw all chars in font
		*/
		i := 0
		for y := 1; y < con.Height-1; y++ {
			for x := 1; x < con.Width-1; x++ {
				con.PutCharInt(x, y, i)
				con.PutColor(x, y, ramen.NewColor(255, 0, 0), ramen.ForegroundColor)
				i++
			}
		}
		return nil
	})

	con.Start(2)
}
