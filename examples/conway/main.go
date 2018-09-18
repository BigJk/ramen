package main

import (
	"math/rand"
	"time"

	"github.com/BigJk/ramen"
	"github.com/hajimehoshi/ebiten"
)

const (
	Width  = 50
	Height = 50
)

var board [][]bool

/*
	Conway's Game Of Life

	Controls:
		- Place cell: Left Mouse Button
		- Step: Space
*/

func main() {
	board = createBoard(Width, Height)

	for x := range board {
		for y := range board[x] {
			board[x][y] = rand.Intn(5) <= 1
		}
	}

	font, err := ramen.NewFont("../../fonts/ti84-6x8.png", 6, 8)
	if err != nil {
		panic(err)
	}

	con, err := ramen.NewConsole(Width, Height, font, "ramen - conway's game of life")
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			if ebiten.IsKeyPressed(ebiten.KeySpace) {
				step()
			}

			if ebiten.IsKeyPressed(ebiten.KeyK) {
				board = createBoard(Width, Height)
			}

			for x := range board {
				for y := range board[x] {
					if board[x][y] {
						con.PutColor(x, y, ramen.NewColor(255, 255, 255), ramen.ModifyBackgroundColor)
					} else {
						con.PutColor(x, y, ramen.NewColor(0, 0, 0), ramen.ModifyBackgroundColor)
					}
				}
			}

			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				cx, cy := ebiten.CursorPosition()

				bx := cx / font.TileWidth
				by := cy / font.TileHeight

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
	wx := x % Width
	wy := y % Height
	if wx < 0 {
		wx += Width
	}
	if wy < 0 {
		wy += Height
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
	nextBoard := createBoard(Width, Height)
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
