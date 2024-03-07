package tictacgoe

import (
	"bytes"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gitlab.com/soockee/tictacgoe/assets/fonts"
	"gitlab.com/soockee/tictacgoe/configs"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Config struct {
	screenWidth  int
	screenHeight int
	boardSize    int
	font         font.Face
}

func NewConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(configs.Config_yml))
	config := &Config{
		screenWidth:  viper.GetInt("screen.screenWidth"),
		screenHeight: viper.GetInt("screen.screenHeight"),
		boardSize:    viper.GetInt("screen.boardSize"),
	}

	const dpi = 72
	tt, err := opentype.Parse(fonts.MarioFont)
	if err != nil {
		log.Fatal().Err(err)
	}
	f, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal().Err(err)
	}
	config.font = f

	return config
}
