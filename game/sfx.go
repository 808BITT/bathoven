package game

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

type SoundEffect struct {
	SampleRate int
	Context    *audio.Context
	Sound      *audio.Player
}

func NewSoundEffect(context *audio.Context) *SoundEffect {
	player, err := mp3.DecodeWithoutResampling(bytes.NewReader(LoadSound("assets/death.mp3")))
	if err != nil {
		log.Fatal(err)
	}
	sound, err := context.NewPlayer(player)
	if err != nil {
		log.Fatal(err)
	}

	return &SoundEffect{
		SampleRate: 48000,
		Context:    context,
		Sound:      sound,
	}
}
