package game

import "github.com/hajimehoshi/ebiten/v2"

type ScrollingBackground struct {
	Width  int
	Height int
	X      int
	Image  *ebiten.Image
}

func NewViewport(width, height int) *ScrollingBackground {
	return &ScrollingBackground{
		Width:  width,
		Height: height,
		X:      0,
		Image:  LoadImage("assets/bg3.png"),
	}
}

func (v *ScrollingBackground) Update() {
	s := v.Image.Bounds().Size()
	maxX := s.X * 16

	v.X += s.X / 200
	v.X %= maxX
}

func (p *ScrollingBackground) Position() int {
	return p.X
}

func (v *ScrollingBackground) Draw(screen *ebiten.Image) {
	x := v.Position()
	bgScrollOffset := float64(-x) / 8

	const repeat = 3
	w, _ := v.Image.Bounds().Dx(), v.Image.Bounds().Dy()
	for i := 0; i < repeat; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(w*i), 0)
		op.GeoM.Translate(bgScrollOffset, 0)
		screen.DrawImage(v.Image, op)
	}

}
