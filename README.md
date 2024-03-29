<img src="https://cdn.rawgit.com/BigJk/7e61395616df18c9b6003aa90c77e829/raw/ec7bc03e02015deb0c96c6914f5c0460af773b59/ramen.svg" width="145" align="left" />

<img src="https://i.imgur.com/glpKbxk.png" width="10" height="145" align="left" />

[![Documentation](https://godoc.org/github.com/BigJk/ramen/console?status.svg)](http://godoc.org/github.com/BigJk/ramen/console) [![Go Report Card](https://goreportcard.com/badge/github.com/BigJk/ramen)](https://goreportcard.com/report/github.com/BigJk/ramen) [![Hex.pm](https://img.shields.io/hexpm/l/plug.svg)](LICENSE)

**ramen** is a simple console emulator written in go that can be used to create various ascii / text (roguelike) games. It's based on the great **[ebiten](https://github.com/hajimehoshi/ebiten)** library and inspired by libraries like **[libtcod](https://github.com/libtcod/libtcod)**.

**Warning:** API and features are not fixed yet. Bugs will happen!

<br>

## Features

- PNG Fonts with more than 256 chars possible
- Fonts can contain chars and colored tiles
- Create sub-consoles to organize rendering
- Component based ui system
- Inlined color definitions in strings
- Pre-build components ready to use
  - TextBox
  - Button
- REXPaint file parsing
- Everything **ebiten** can do
  - Input: Mouse, Keyboard, Gamepads, Touches
  - Audio: MP3, Ogg/Vorbis, WAV, PCM
  - ...

## Get Started

```
go get github.com/BigJk/ramen/...
```

## Transformer

In ramen you change the content of the console by applying transformations to cells. Examples would be:

```go
// set one cell at position 10,15 to a green @:
con.Transform(10, 15, t.CharByte('@'), t.Foreground(concolor.RGB(0, 255, 0)))

// change the background of the area 0,0 with the width and height of 25,25:
con.TransformArea(0, 0, 25, 25, t.Background(concolor.RGBA(255, 255, 255, 20)))

// change the background of all the cells:
con.TransformAll(t.Background(concolor.RGBA(255, 255, 255, 10)))
```

All transformer functions accept objects that implement the **t.Transformer** interface, so it's also possible to create transformers with custom behaviour by implementing that interface.

## Inlined Color Definitions

There are also convenient string printing functions. The **console.Print** function supports parsing of inlined color definitions that can change the forground and background color of parts of the string.

``[[f:#ff0000]]red foreground\n[[f:#ffffff|b:#000000]]white foreground and black background\n[[b:#00ff00]]green background``

<img src="./.github/screen_colored_string.png" width="400">

## Example

```go
package main

import (
  "fmt"
  "github.com/BigJk/ramen/concolor"
  "github.com/BigJk/ramen/console"
  "github.com/BigJk/ramen/font"
  "github.com/BigJk/ramen/t"
  "github.com/hajimehoshi/ebiten/v2"
)

func main() {
  // create a 50x30 cells console with the title 'ramen example'
  con, err := console.New(50, 30, font.DefaultFont, "ramen example")
  if err != nil {
    panic(err)
  }

  // set a tick hook. This function will be executed
  // each tick (60 ticks per second by default) even
  // when the fps is lower than 60fps. This is a good
  // place for your game logic.
  //
  // The timeDelta parameter is the elapsed time in seconds
  // since the last tick.
  con.SetTickHook(func(timeElapsed float64) error {
    // your game logic
    return nil
  })

  // set a pre-render hook. This function will be executed
  // each frame before the drawing happens. This is a good
  // place to draw onto the console, because it only executes
  // if a draw is really about to happen.
  //
  // The timeDelta parameter is the elapsed time in seconds
  // since the last frame.
  con.SetPreRenderHook(func(screen *ebiten.Image, timeDelta float64) error {
    con.ClearAll() // clear console 
    con.TransformAll(t.Background(concolor.RGB(50, 50, 50))) // set the background
	
    con.Print(2, 2, "Hello!\nTEST\n Line 3", t.Foreground(concolor.RGB(0, 255, 0)), t.Background(concolor.RGB(255, 0, 0)))
    con.Print(2, 7, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\nElapsed: %0.4f", ebiten.CurrentFPS(), ebiten.CurrentFPS(), timeDelta))
	
    return nil
  })

  // start the console with a scaling of 1
  con.Start(1)
}
```

## Screenshots

<img src="./.github/screen_colored_tiles.png" width="538">
<img src="./.github/screen_text.png" width="200">
<img src="./.github/screen_comp_shaded.png" width="314">
