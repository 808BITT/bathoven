package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InputData struct {
	Jump               bool
	CycleOutfit        bool
	Pause              bool
	ActivateFullscreen bool
	ExitFullscreen     bool
}

func NewInputData() *InputData {
	return &InputData{
		Jump:               false,
		CycleOutfit:        false,
		Pause:              false,
		ActivateFullscreen: false,
		ExitFullscreen:     false,
	}
}

func (i *InputData) Update() *InputData {
	i.Reset()

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		i.SetJump()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyO) {
		i.SetCycleOutfit()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		i.SetPause()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		i.SetActivateFullscreen()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		i.SetExitFullscreen()
	}

	return i
}

func (i *InputData) SetJump() {
	i.Jump = true
}

func (i *InputData) SetCycleOutfit() {
	i.CycleOutfit = true
}

func (i *InputData) SetPause() {
	i.Pause = true
}

func (i *InputData) SetActivateFullscreen() {
	i.ActivateFullscreen = true
}

func (i *InputData) SetExitFullscreen() {
	i.ExitFullscreen = true
}

func (i *InputData) Reset() {
	if i == nil {
		return
	}
	i.Jump = false
	i.CycleOutfit = false
	i.Pause = false
	i.ActivateFullscreen = false
	i.ExitFullscreen = false
}
