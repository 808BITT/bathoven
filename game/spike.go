package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func NewSpike(speed float64, x, y int, img *ebiten.Image, top bool) *Spike {
	return &Spike{
		X:      &x,
		Y:      &y,
		Sprite: img,
		Top:    top,
		Speed:  speed,
	}
}

type Spike struct {
	X      *int
	Y      *int
	Sprite *ebiten.Image
	Top    bool
	Speed  float64
}

func (s *Spike) Update() {
	*s.X -= int(s.Speed)
}

func (s *Spike) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(*s.X), float64(*s.Y))
	screen.DrawImage(s.Sprite, op)
}

func (s *Spike) Bounds() (int, int) {
	return s.Sprite.Bounds().Dx(), s.Sprite.Bounds().Dy()
}

func (s *Spike) Position() (int, int) {
	return *s.X, *s.Y
}

func (s *Spike) CollidesWith(p *Player) bool {
	px, py := p.Position()
	pw, ph := p.Bounds()

	sx, sy := s.Position()
	sw, sh := s.Bounds()

	sx += sw / 4
	sw /= 4

	px += (pw / 10) * 2
	pw -= (pw / 10) * 4
	py += (ph / 10) * 2
	ph -= (ph / 10) * 4

	if px < sx+sw && px+pw > sx && py < sy+sh && py+ph > sy {
		return true
	}
	return false
}
