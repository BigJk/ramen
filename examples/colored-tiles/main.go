package main

import (
	"image/color"

	"github.com/BigJk/ramen/console"
	"github.com/BigJk/ramen/consolecolor"
	"github.com/BigJk/ramen/font"
	"github.com/BigJk/ramen/t"
	"github.com/hajimehoshi/ebiten"
)

func main() {
	/*
		load a font that contains colored tiles
	*/
	font, err := font.New("./custom_font.png", 8, 8)
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
	con, err := console.New(40, 20, font, "ramen - colored tiles")
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
		con.Transform(0, 0, t.Char(272))
		con.Transform(1, 0, t.Char(276))
		con.Transform(0, 1, t.Char(280))

		con.Transform(con.Width-1, 0, t.Char(273))
		con.Transform(con.Width-2, 0, t.Char(277))
		con.Transform(con.Width-1, 1, t.Char(281))

		con.Transform(0, con.Height-1, t.Char(274))
		con.Transform(1, con.Height-1, t.Char(278))
		con.Transform(0, con.Height-2, t.Char(282))

		con.Transform(con.Width-1, con.Height-1, t.Char(275))
		con.Transform(con.Width-2, con.Height-1, t.Char(279))
		con.Transform(con.Width-1, con.Height-2, t.Char(283))

		/*
			draw all chars in font
		*/
		i := 0
		for y := 1; y < con.Height-1; y++ {
			for x := 1; x < con.Width-1; x++ {
				con.Transform(x, y, t.Char(i), t.Foreground(consolecolor.New(255, 0, 0)))
				i++
			}
		}
		return nil
	})

	con.Start(2)
}
