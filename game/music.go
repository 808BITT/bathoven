package game

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

type MusicPlayer struct {
	SampleRate int
	Context    *audio.Context
	Sound      *audio.Player
}

func NewMusicPlayer(path string) *MusicPlayer {
	context := audio.NewContext(44100)
	player, err := mp3.DecodeWithoutResampling(bytes.NewReader(LoadSound(path)))
	if err != nil {
		log.Fatal(err)
	}
	sound, err := context.NewPlayer(player)
	if err != nil {
		log.Fatal(err)
	}

	return &MusicPlayer{
		SampleRate: 48000,
		Context:    context,
		Sound:      sound,
	}
}

func LoadSound(path string) []byte {
	f, err := assets.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return f
}
