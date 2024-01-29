package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

func NewSpike(x, y int, img *ebiten.Image, top bool) *Spike {
	return &Spike{
		X:      &x,
		Y:      &y,
		Sprite: img,
		Top:    top,
	}
}

type Spike struct {
	X      *int
	Y      *int
	Sprite *ebiten.Image
	Top    bool
}

func (s *Spike) Update() {
	// move the spike to the left
	*s.X -= 30
}

func (s *Spike) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	if s.Top {
		if *s.Y > 0 {
			fmt.Println("s.Y: ", *s.Y, " s.Sprite.Bounds().Dy(): ", s.Sprite.Bounds().Dy())

			// scale proportionally so the spike bottom is at the same place but the top is at the top of the screen
			op.GeoM.Scale(1, float64(*s.Y)/float64(s.Sprite.Bounds().Dy())*0.2)
			*s.Y = -10
		}
	} else {
		if *s.Y < 2160 {
			// scale proportionally so the spike top is at the same place but the bottom is at the bottom of the screen
			op.GeoM.Scale(1, float64(2160)/float64(s.Sprite.Bounds().Dy())+0.2)
		}
	}
	op.GeoM.Scale(2, 1)
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

	// adjust the spike hitbox to be centered on the sprite
	sx += sw / 2
	sw /= 2

	if px < sx+sw && px+pw > sx && py < sy+sh && py+ph > sy {
		return true
	}
	return false
}
