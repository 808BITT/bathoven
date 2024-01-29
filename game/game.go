package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Engine *Engine
}

func Initialize() Game {
	return Game{
		Engine: NewEngine(3840, 2160),
	}
}

func (g *Game) Update() error {
	g.Engine.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Engine.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.Engine.Background.Width, g.Engine.Background.Height
}

func (g *Game) SaveHighScore() {
	g.Engine.SaveHighScore()
}
