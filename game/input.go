package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InputManager struct {
	InputData *InputData
}

type InputData struct {
	Jump bool
}

func NewInputData() *InputData {
	return &InputData{
		Jump: false,
	}
}

func (i *InputData) SetJump() {
	i.Jump = true
}

func (i *InputData) Reset() {
	if i == nil {
		return
	}
	i.Jump = false
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
	} else {
		i.InputData.Reset()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	return i.InputData
}
