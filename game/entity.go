package game

import (
	"crypto/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type EntityManager struct {
	Player        *Player
	Spikes        *[]Spike
	SpikeDistance float64
	SpikeGap      int
	SpikeSpeed    float64
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		Player:        NewPlayer(500, 1000, LoadImage(outfits[0])),
		Spikes:        &[]Spike{},
		SpikeDistance: 1800,
		SpikeGap:      300,
		SpikeSpeed:    15,
	}
}

func (e *EntityManager) Collide() bool {
	for _, spike := range *e.Spikes {
		if spike.CollidesWith(e.Player) {
			return true
		}
	}
	return false
}

func (e *EntityManager) Update(input *InputData) int {
	e.Player.Update(input)

	// update spikes
	for i, spike := range *e.Spikes {
		spike.Update()
		if *spike.X < -spike.Sprite.Bounds().Dx() {
			*e.Spikes = append((*e.Spikes)[:i], (*e.Spikes)[i+1:]...)
			return 1
		}
	}

	// add spikes if needed
	if len(*e.Spikes) < 40 {
		//make sure the spikes are not too close to each other
		for _, spike := range *e.Spikes {
			if *spike.X > 3840-int(e.SpikeDistance) {
				return 0
			}
		}
		difficulty := int(1800/e.SpikeDistance) * 15
		random := 0
		b := make([]byte, 1)
		_, err := rand.Read(b)
		if err != nil {
			panic(err)
		}
		random = int(b[0]) - 64
		random *= 2
		*e.Spikes = append(*e.Spikes, *NewSpike(e.SpikeSpeed, 3840+random*2, -500-e.SpikeGap+(difficulty*4)+random, LoadImage("assets/tall_mite2.png"), true))
		*e.Spikes = append(*e.Spikes, *NewSpike(e.SpikeSpeed, 3840+random*2, 1000+e.SpikeGap-(difficulty*4)+random, LoadImage("assets/tall_mite1.png"), false))

		if e.SpikeDistance >= 900 {
			e.SpikeDistance -= 100
		}

		if e.SpikeSpeed <= 30 {
			e.SpikeSpeed += 0.1
		}
	}
	return 0
}

func (e *EntityManager) Draw(screen *ebiten.Image) {
	e.Player.Draw(screen)
	for _, spike := range *e.Spikes {
		spike.Draw(screen)
	}
}
