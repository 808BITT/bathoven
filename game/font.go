package game

import (
	"image/color"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type FontManager struct {
	Normal font.Face
	Big    font.Face
	Green  color.Color
	White  color.Color
	Red    color.Color
}

func NewFontManager() *FontManager {
	return &FontManager{
		Normal: LoadFont("assets/SuperLegendBoy.ttf", 100),
		Big:    LoadFont("assets/SuperLegendBoy.ttf", 200),
		White:  color.RGBA{255, 255, 255, 255},
		Green:  color.RGBA{34, 139, 34, 255},
		Red:    color.RGBA{255, 0, 0, 255},
	}
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
