package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	X      int
	Y      int
	dy     float64
	Sprite *ebiten.Image
}

func NewPlayer(x, y int, img *ebiten.Image) *Player {
	return &Player{
		X:      x,
		Y:      y,
		dy:     -50,
		Sprite: img,
	}
}

func (p *Player) Update(input *InputData) {
	if input.Jump && p.Y > 0 {
		p.dy = -20
	}

	p.dy += 1
	if p.dy > 20 {
		p.dy = 20
	}
	p.Y += int(p.dy)

	if p.Y > 2160-p.Sprite.Bounds().Dy()-200 {
		p.Y = 2160 - p.Sprite.Bounds().Dy() - 200
		p.dy = 0
	}

	if p.Y < 0 {
		p.Y = 0
		p.dy = 0
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.X), float64(p.Y))
	screen.DrawImage(p.Sprite, op)
}

func (p *Player) Bounds() (int, int) {
	return p.Sprite.Bounds().Dx() - 100, p.Sprite.Bounds().Dy()
}

func (p *Player) Position() (int, int) {
	return p.X, p.Y
}
