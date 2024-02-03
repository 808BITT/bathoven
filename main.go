package main

import (
	"bathoven/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// TODO: add a menu
	// TODO: add a pause menu
	e := game.NewEngine(3840, 2160)
	ebiten.SetWindowSize(e.Background.Width/2, e.Background.Height/2)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Bathoven")
	// ebiten.SetFullscreen(true)
	err := ebiten.RunGame(e)
	if err != nil {
		panic(err)
	}
}
