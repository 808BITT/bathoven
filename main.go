package main

import (
	"bathoven/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// TODO: add a menu
	// TODO: add a pause menu
	g := game.Initialize()
	ebiten.SetWindowSize(g.Engine.Background.Width/4, g.Engine.Background.Height/4)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Bathoven")
	// ebiten.SetFullscreen(true)
	err := ebiten.RunGame(&g)
	if err != nil {
		panic(err)
	}
}
