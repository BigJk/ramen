package main

import (
	"fmt"
	"github.com/BigJk/ramen/components"
	"github.com/BigJk/ramen/concolor"
	"github.com/BigJk/ramen/console"
	"github.com/BigJk/ramen/font"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	con, err := console.New(60, 30, font.DefaultFont, "ramen - component example")
	if err != nil {
		panic(err)
	}

	// Create a "Test Button" button.
	btn := components.NewButton(5, 5, 15, 5, "Test Button", func() {
		fmt.Println("Button Pressed!")
	})
	btn.SetBackground(concolor.RGB(50, 50, 50).P(), concolor.RGB(70, 70, 70).P(), concolor.RGB(30, 30, 30).P())
	btn.SetForeground(concolor.White.P(), concolor.White.P(), concolor.White.P())
	con.AddComponent(btn)

	// Create a Textbox.
	txtInput := components.NewTextbox(5, 12, 15, 1)
	txtInput.SetBackground(concolor.RGB(50, 50, 50).P(), concolor.RGB(70, 70, 70).P(), concolor.RGB(30, 30, 30).P())
	txtInput.SetForeground(concolor.White.P(), concolor.RGB(90, 90, 90).P())
	txtInput.SetEnterCallback(func(text string) {
		fmt.Println("Text:", text)
		txtInput.SetText("")
	})
	con.AddComponent(txtInput)

	// Create a "Toggle Button" button that will hide and un-hide the first button.
	btnHide := components.NewButton(5, 15, 15, 5, "Toggle Button", func() {
		btn.Show(!btn.ShouldDraw())
	})
	btnHide.SetBackground(concolor.RGB(50, 50, 50).P(), concolor.RGB(70, 70, 70).P(), concolor.RGB(30, 30, 30).P())
	btnHide.SetForeground(concolor.White.P(), concolor.White.P(), concolor.White.P())
	con.AddComponent(btnHide)

	// Just clear on each render
	con.SetPreRenderHook(func(screen *ebiten.Image, timeDelta float64) error {
		con.ClearAll()
		return nil
	})

	con.ShowFPS = true

	// start the console with a scaling of 1
	con.Start(1)
}
