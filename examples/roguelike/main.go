package main

import (
	"fmt"
	"github.com/BigJk/ramen/concolor"
	"github.com/BigJk/ramen/console"
	"github.com/BigJk/ramen/font"
	"github.com/BigJk/ramen/t"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"strings"
)

var player = struct {
	X int
	Y int
}{3, 3}

var levelLayout = `
#####################
#         #    #    #
#    #    #         #
#    ######    #    #
#              #    #
##  #############  ##
#    #    #    #    #
#    #         #    #
#    ######         #
#              #    #
#####################
`

var world [][]byte

// converts levelLayout to world
func initWorld() {
	lines := strings.Split(levelLayout, "\n")

	for i := range lines {
		if len(lines[i]) == 0 {
			continue
		}

		world = append(world, []byte(lines[i]))
	}
}

// checks if a tile is solid (tile content is not a space ' ' character)
func isSolid(x int, y int) bool {
	return world[y][x] != ' '
}

func main() {
	initWorld()

	rootView, err := console.New(60, 35, font.DefaultFont, "ramen - roguelike example")
	if err != nil {
		panic(err)
	}

	worldView, err := rootView.CreateSubConsole(0, 1, rootView.Width-20, rootView.Height-1)
	if err != nil {
		panic(err)
	}

	playerInfoView, err := rootView.CreateSubConsole(worldView.Width, 1, 20, rootView.Height-1)
	if err != nil {
		panic(err)
	}

	/*
		move player on key press
	*/

	rootView.SetTickHook(func(timeElapsed float64) error {
		if inpututil.IsKeyJustPressed(ebiten.KeyW) && !isSolid(player.X, player.Y-1) {
			player.Y -= 1
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyS) && !isSolid(player.X, player.Y+1) {
			player.Y += 1
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyA) && !isSolid(player.X-1, player.Y) {
			player.X -= 1
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyD) && !isSolid(player.X+1, player.Y) {
			player.X += 1
		}

		return nil
	})

	/*
		render
	*/

	rootView.SetPreRenderHook(func(screen *ebiten.Image, timeDelta float64) error {
		/*
			clear console
		*/

		rootView.ClearAll()
		rootView.TransformAll(t.Background(concolor.RGB(50, 50, 50)))

		worldView.ClearAll()
		worldView.TransformAll(t.Background(concolor.RGB(55, 55, 55)), t.Char(0))

		playerInfoView.ClearAll()

		/*
			draw header
		*/

		rootView.TransformArea(0, 0, rootView.Width, 1, t.Background(concolor.RGB(80, 80, 80)))
		rootView.Print(2, 0, "World View", t.Foreground(concolor.White))
		rootView.Print(worldView.Width+2, 0, "Player Info", t.Foreground(concolor.White))

		/*
			draw world
		*/

		midX := worldView.Width / 2
		midY := worldView.Height / 2

		for y := range world {
			for x := range world[y] {
				if world[y][x] == ' ' {
					continue
				}

				worldView.Transform(midX-player.X+x, midY-player.Y+y, t.CharByte(world[y][x]))
			}
		}

		/*
			draw player in the middle
		*/

		worldView.Transform(midX, midY, t.CharByte('@'), t.Foreground(concolor.Green))

		/*
			draw player info
		*/
		playerInfoView.PrintBounded(1, 1, playerInfoView.Width-2, 2, fmt.Sprintf("X=%d Y=%d", player.X, player.Y))

		return nil
	})

	rootView.Start(2)
}
