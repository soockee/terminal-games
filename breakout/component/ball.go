package component

import (
	"time"

	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/engine"
	"github.com/yohamta/donburi"
)

type BallData struct {
	Speed                   float64
	Shape                   *resolv.Circle
	MaxSpeed                float64
	CollisionCooldownBlock  time.Duration
	CollisionCooldownPlayer time.Duration
	CooldownTimer           engine.Timer
}

var Ball = donburi.NewComponentType[BallData]()
