package main

import (
	"math/rand"
	"time"

	"github.com/BigJk/ramen/concolor"
	"github.com/BigJk/ramen/console"
	"github.com/BigJk/ramen/font"
	"github.com/BigJk/ramen/t"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	width  = 50
	height = 50
)

var board [][]bool

/*
	Conway's Game Of Life

	Controls:
		- Place cell: Left Mouse Button
		- Step: Space
*/

func main() {
	board = createBoard(width, height)

	for x := range board {
		for y := range board[x] {
			board[x][y] = rand.Intn(5) <= 1
		}
	}

	con, err := console.New(width, height, font.DefaultFont, "ramen - conway's game of life")
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			if ebiten.IsKeyPressed(ebiten.KeySpace) {
				step()
			}

			if ebiten.IsKeyPressed(ebiten.KeyK) {
				board = createBoard(width, height)
			}

			for x := range board {
				for y := range board[x] {
					if board[x][y] {
						con.Transform(x, y, t.Background(concolor.RGB(255, 255, 255)))
					} else {
						con.Transform(x, y, t.Background(concolor.RGB(0, 0, 0)))
					}
				}
			}

			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				cx, cy := ebiten.CursorPosition()

				bx := cx / font.DefaultFont.TileWidth
				by := cy / font.DefaultFont.TileHeight

				board[bx][by] = true
			} else {
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()

	con.Start(1)
}

func createBoard(w, h int) [][]bool {
	b := make([][]bool, w)
	for i := range b {
		b[i] = make([]bool, h)
	}
	return b
}

func getCell(x, y int) bool {
	wx := x % width
	wy := y % height
	if wx < 0 {
		wx += width
	}
	if wy < 0 {
		wy += height
	}
	return board[wx][wy]
}

func livingNeighbors(x, y int) int {
	c := 0

	if getCell(x-1, y) {
		c++
	}

	if getCell(x+1, y) {
		c++
	}

	if getCell(x, y-1) {
		c++
	}

	if getCell(x, y+1) {
		c++
	}

	if getCell(x-1, y-1) {
		c++
	}

	if getCell(x+1, y+1) {
		c++
	}

	if getCell(x-1, y+1) {
		c++
	}

	if getCell(x+1, y-1) {
		c++
	}

	return c
}

func step() {
	nextBoard := createBoard(width, height)
	for x := range board {
		for y := range board[x] {
			c := livingNeighbors(x, y)
			if c < 2 || c > 3 {
				nextBoard[x][y] = false
			} else if board[x][y] && (c == 2 || c == 3) {
				nextBoard[x][y] = true
			} else if !board[x][y] && c == 3 {
				nextBoard[x][y] = true
			}
		}
	}
	board = nextBoard
}
