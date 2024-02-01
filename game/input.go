package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InputManager struct {
	InputData *InputData
}

type InputData struct {
	Jump        bool
	CycleOutfit bool
}

func NewInputData() *InputData {
	return &InputData{
		Jump:        false,
		CycleOutfit: false,
	}
}

func (i *InputData) SetJump() {
	i.Jump = true
}

func (i *InputData) SetCycleOutfit() {
	i.CycleOutfit = true
}

func (i *InputData) Reset() {
	if i == nil {
		return
	}
	i.Jump = false
	i.CycleOutfit = false
}

func NewInputManager() *InputManager {
	return &InputManager{
		InputData: NewInputData(),
	}
}

func (i *InputManager) Update() *InputData {
	if i.InputData == nil {
		i.InputData = NewInputData()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		i.InputData.SetJump()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyO) {
		i.InputData.SetCycleOutfit()
	} else {
		i.InputData.Reset()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(true)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) && ebiten.IsFullscreen() {
		ebiten.SetFullscreen(false)
	}
	return i.InputData
}
