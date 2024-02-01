package game

import (
	"bufio"
	"crypto/rand"
	"embed"
	"fmt"
	"image/color"
	"image/png"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed assets/*.png assets/*.ttf assets/*.mp3
var assets embed.FS

type Engine struct {
	Background        *Viewport
	Player            *Player
	Spikes            *[]Spike
	Input             *InputManager
	Score             int
	HighScore         int
	Font              font.Face
	ScoreFontColor    color.Color
	CongratsFontColor color.Color
	CurrentCongrats   string
	DeathFont         font.Face
	DeathFontColor    color.Color
	DeathScreen       bool
	SpikeDistance     float64
	MusicPlayer       *MusicPlayer
}

func NewEngine(w, h int) *Engine {
	return &Engine{
		Background:        NewViewport(w, h, LoadImage("assets/bg3.png")),
		Player:            NewPlayer(300, 1000, LoadImage(outfits[0])),
		Spikes:            &[]Spike{},
		Input:             &InputManager{},
		Score:             0,
		HighScore:         LoadHighScore(),
		Font:              LoadFont("assets/SuperLegendBoy.ttf", 100),
		ScoreFontColor:    color.RGBA{255, 255, 255, 255},
		CongratsFontColor: color.RGBA{34, 139, 34, 255},
		CurrentCongrats:   "Nice!",
		DeathFont:         LoadFont("assets/SuperLegendBoy.ttf", 200),
		DeathFontColor:    color.RGBA{255, 0, 0, 255},
		DeathScreen:       false,
		SpikeDistance:     1800,
		MusicPlayer:       NewMusicPlayer("assets/8bithovben.mp3"),
	}
}

func LoadHighScore() int {
	if _, err := os.Stat("highscore.txt"); os.IsNotExist(err) {
		file, err := os.Create("highscore.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		file.WriteString("0")
	}
	file, err := os.Open("highscore.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	highScore, err := strconv.Atoi(scanner.Text())
	if err != nil {
		log.Fatal(err)
	}

	return highScore
}

func LoadFont(fontPath string, fontSize int) font.Face {
	fontBytes, err := assets.ReadFile(fontPath)
	if err != nil {
		log.Fatal(err)
	}
	ttfFont, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}
	fontFace, _ := opentype.NewFace(ttfFont, &opentype.FaceOptions{
		Size:    float64(fontSize),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	return fontFace
}

func (e *Engine) Update() {
	e.SpikeDistance -= 0.1
	if e.SpikeDistance < 400 {
		e.SpikeDistance = 400
	}

	e.Player.Update(e.Input.Update())

	if e.DeathScreen {
		if e.MusicPlayer.Sound.IsPlaying() {
			e.MusicPlayer.Sound.Pause()
		}
		if e.Input.InputData.Jump {
			e.DeathScreen = false
			e.Score = 0
			e.Player.Y = 1000
			e.Player.dy = -50
			e.Background.X = 0
			*e.Spikes = []Spike{}
			e.SpikeDistance = 1800
			e.MusicPlayer.Sound.Rewind()
		}
		return
	}
	if !e.MusicPlayer.Sound.IsPlaying() {
		e.MusicPlayer.Sound.Play()
	}

	// check for collisions
	for _, spike := range *e.Spikes {
		if spike.CollidesWith(e.Player) {
			if e.Score > e.HighScore {
				e.HighScore = e.Score
				e.SaveHighScore()
			}
			e.DeathScreen = true
		}
	}

	e.Background.Move()
	for i, spike := range *e.Spikes {
		spike.Update()
		if *spike.X < -400 {
			if spike.Top {
				e.Score++
			}
			*e.Spikes = append((*e.Spikes)[:i], (*e.Spikes)[i+1:]...)
			break
		}
	}

	// add spikes if needed
	if len(*e.Spikes) < 40 {
		//make sure the spikes are not too close to each other
		for _, spike := range *e.Spikes {
			if *spike.X > 3840-int(e.SpikeDistance) {
				return
			}
		}
		offset, err := rand.Int(rand.Reader, big.NewInt(900))
		if err != nil {
			log.Fatal(err)
		}
		offsetInt := int(offset.Int64() - 450 - int64(2000/e.SpikeDistance))
		*e.Spikes = append(*e.Spikes, *NewSpike(3840, 500-offsetInt, LoadImage("assets/tall_mite2.png"), true))
		*e.Spikes = append(*e.Spikes, *NewSpike(3840, 1410-offsetInt, LoadImage("assets/tall_mite1.png"), false))
	}

}

func (e *Engine) Draw(screen *ebiten.Image) {
	e.Background.Draw(screen)
	e.Player.Draw(screen)
	for _, spike := range *e.Spikes {
		spike.Draw(screen)
	}

	// paint slightly transparent black over the score area to make it easier to read
	vector.DrawFilledRect(screen, float32(0), float32(0), float32(3840), float32(150), color.RGBA{0, 0, 0, 100}, true)

	if e.DeathScreen {
		// paint slightly transparent black over the score area to make it easier to read
		vector.DrawFilledRect(screen, float32(0), float32(0), float32(3840), float32(2160), color.RGBA{0, 0, 0, 150}, true)
		text.Draw(screen, "You Died!", e.DeathFont, 1250, 980, e.DeathFontColor)
		text.Draw(screen, fmt.Sprintf("High Score: %d", e.HighScore), e.Font, 1300, 1200, e.ScoreFontColor)
		text.Draw(screen, "Press Space to Restart", e.Font, 1100, 1400, e.ScoreFontColor)
		return
	}

	// text.Draw(screen, fmt.Sprintf("High Score: %d", e.HighScore), e.ScoreFont, 2600, 100, e.ScoreFontColor)
	text.Draw(screen, fmt.Sprintf("Score: %d", e.Score), e.Font, 100, 120, e.ScoreFontColor)

	if e.Score > e.HighScore {
		text.Draw(screen, "New High Score!", e.Font, 1300, 120, e.CongratsFontColor)
	}

	if e.Score > 0 && e.Score%10 == 0 {
		text.Draw(screen, e.CurrentCongrats, e.DeathFont, 1920-len(e.CurrentCongrats)*80, 600, e.CongratsFontColor)
		go func() {
			time.Sleep(1 * time.Second)
			e.CurrentCongrats = RandomGratz()
		}()
	}

	if e.Score == 0 {
		text.Draw(screen, "Press Space to Jump", e.Font, 1100, 1100, e.ScoreFontColor)
	}
}

func RandomGratz() string {
	gratz := []string{
		"Nice!",
		"Good Job!",
		"Keep Going!",
		"Great!",
		"Awesome!",
		"Amazing!",
		"Unbelievable!",
		"Fantastic!",
		"Excellent!",
		"Outstanding!",
		"Stupendous!",
		"Stunning!",
		"Superb!",
		"Terrific!",
		"Phenomenal!",
		"Marvelous!",
		"Brilliant!",
		"Extraordinary!",
		"Remarkable!",
		"Keep flying!",
		"Great job!",
		"Keep it up!",
		"Soaring high!",
		"Awesome skills!",
		"You've got this!",
		"Fantastic!",
		"Unstoppable!",
		"Superb flying!",
		"Impressive!",
		"You're a star!",
		"Sky's the limit!",
		"Amazing work!",
		"Flying master!",
		"Unbelievable!",
	}
	randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(gratz))))
	if err != nil {
		log.Fatal(err)
	}
	return gratz[randIndex.Int64()]
}

func LoadImage(path string) *ebiten.Image {
	f, err := assets.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		log.Fatalf("Failed to decode image file: %v", err)
	}

	image := ebiten.NewImageFromImage(img)
	return image
}

func (e *Engine) SaveHighScore() {
	f, err := os.Create("highscore.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	fmt.Fprintf(w, "%d", e.HighScore)
}
