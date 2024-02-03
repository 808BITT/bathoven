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
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

//go:embed assets/*.png assets/*.ttf assets/*.mp3
var assets embed.FS

const DeathScreenTime = 60

type Engine struct {
	Background      *ScrollingBackground
	MusicPlayer     *MusicPlayer
	DeathSound      *SoundEffect
	Entities        *EntityManager
	Input           *InputData
	Fonts           *FontManager
	Score           int
	ScoreTimer      int
	HighScore       int
	CurrentCongrats string
	Paused          bool
	DeathScreen     bool
	DeathTimer      int
	SpikeDistance   float64
}

func NewEngine(w, h int) *Engine {
	context := audio.NewContext(44100)
	return &Engine{
		Background:      NewViewport(w, h),
		MusicPlayer:     NewMusicPlayer(context),
		DeathSound:      NewSoundEffect(context),
		Entities:        NewEntityManager(),
		Input:           NewInputData(),
		Fonts:           NewFontManager(),
		CurrentCongrats: RandomGratz(),
		Paused:          false,
		DeathScreen:     false,
		DeathTimer:      0,
		SpikeDistance:   1800,
		Score:           0,
		ScoreTimer:      0,
		HighScore:       LoadHighScore(),
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

func (e *Engine) Update() error {

	input := e.Input.Update()

	if input.CycleOutfit {
		e.Entities.Player.CycleOutfit()
	}

	if e.Input.ActivateFullscreen {
		if !ebiten.IsFullscreen() {
			ebiten.SetFullscreen(true)
		}
	}

	if e.Input.ExitFullscreen {
		if ebiten.IsFullscreen() {
			ebiten.SetFullscreen(false)
		}
	}

	if e.DeathScreen {
		if e.MusicPlayer.Sound.IsPlaying() {
			e.DeathTimer = 0
			e.MusicPlayer.Sound.Pause()
			e.DeathSound.Sound.Play()
		}
		if e.Input.Jump && e.DeathTimer > DeathScreenTime {
			// TODO: replace this with a restart function.. (probably in the engine struct)
			e.DeathScreen = false
			e.DeathTimer = 0
			e.Score = 0
			e.ScoreTimer = 0
			e.Entities.Player.Y = 1000
			e.Entities.Player.dy = -50
			e.Background.X = 0
			*e.Entities = *NewEntityManager()
			e.SpikeDistance = 1800
			e.DeathSound.Sound.Rewind()
			e.MusicPlayer.Sound.Rewind()
		}
		e.DeathTimer++
		return nil
	}

	if !e.MusicPlayer.Sound.IsPlaying() {
		e.MusicPlayer.Sound.Play()
	}

	if input.Pause {
		e.Paused = !e.Paused
	}

	if e.Paused {
		e.MusicPlayer.Sound.Pause()
		if e.Input.Jump {
			e.Paused = false
			e.MusicPlayer.Sound.Play()
		}
		return nil
	}

	e.Background.Update()
	e.Entities.Update(input)

	if e.Entities.Collide() {
		e.DeathScreen = true
		if e.Score > e.HighScore {
			e.HighScore = e.Score
			e.SaveHighScore()
		}
	}

	if !e.DeathScreen {
		e.ScoreTimer += 1
		e.Score = e.ScoreTimer / 60
	}

	return nil
}

func (e *Engine) Draw(screen *ebiten.Image) {
	e.Background.Draw(screen)
	e.Entities.Draw(screen)

	if e.DeathScreen {
		vector.DrawFilledRect(screen, float32(0), float32(0), float32(3840), float32(2160), color.RGBA{0, 0, 0, 150}, true)
		text.Draw(screen, "You Died!", e.Fonts.Big, 1250, 980, e.Fonts.Red)
		text.Draw(screen, fmt.Sprintf("High Score: %d", e.HighScore), e.Fonts.Normal, 1300, 1200, e.Fonts.White)
		if e.DeathTimer > DeathScreenTime {
			text.Draw(screen, "Press Space to Restart", e.Fonts.Normal, 1100, 1400, e.Fonts.White)
		}
		return
	}

	vector.DrawFilledRect(screen, float32(0), float32(0), float32(3840), float32(150), color.RGBA{0, 0, 0, 100}, true)
	text.Draw(screen, fmt.Sprintf("Score: %d", e.Score), e.Fonts.Normal, 100, 120, e.Fonts.White)
	text.Draw(screen, fmt.Sprintf("High Score: %d", e.HighScore), e.Fonts.Normal, 2600, 100, e.Fonts.White)

	if e.Paused {
		vector.DrawFilledRect(screen, float32(0), float32(0), float32(3840), float32(2160), color.RGBA{0, 0, 0, 150}, true)
		text.Draw(screen, "Paused", e.Fonts.Big, 1400, 960, e.Fonts.White)
		return
	}

	if e.Score > e.HighScore {
		text.Draw(screen, "New High Score!", e.Fonts.Normal, 1300, 120, e.Fonts.Green)
	}

	if e.Score > 0 && e.Score%10 == 0 {
		text.Draw(screen, e.CurrentCongrats, e.Fonts.Big, 1920-len(e.CurrentCongrats)*80, 600, e.Fonts.Green)
		go func() {
			time.Sleep(3 * time.Second)
			e.CurrentCongrats = RandomGratz()
		}()
	}

	if e.Score == 0 {
		text.Draw(screen, "Press Space to Jump", e.Fonts.Normal, 1200, 1100, e.Fonts.White)
	}
}

func (e *Engine) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 3840, 2160
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
