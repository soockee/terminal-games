package assets

import (
	_ "embed"
	_ "image/png"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	//go:embed fonts/kenney_pixel.ttf
	NormalFontData []byte
	//go:embed fonts/kenney_pixel_square.ttf
	SquareFontData []byte

	SmallFont  font.Face
	NormalFont font.Face
	SqaureFont font.Face
)

func MustLoadAssets() {
	SmallFont = MustLoadFont(NormalFontData, 10)
	NormalFont = MustLoadFont(NormalFontData, 24)
	SqaureFont = MustLoadFont(SquareFontData, 24)
}

func MustLoadFont(data []byte, size int) font.Face {
	f, err := opentype.Parse(data)
	if err != nil {
		panic(err)
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	return face
}
