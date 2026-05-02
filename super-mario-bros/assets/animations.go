package assets

import (
	"time"

	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/yohamta/ganim8/v2"
)

// Tileset paths relative to the embedded FS root.
const (
	tilesetCharacter  = "ldtk/tilesets/character.png"
	tilesetKoopaWalk  = "ldtk/tilesets/koopa_troopa_walk.png"
	tilesetKoopaShell = "ldtk/tilesets/koopa_troopa_shell.png"
	tilesetGoombaWalk = "ldtk/tilesets/goomba_walk.png"
	tilesetGoombaDead = "ldtk/tilesets/goomba_squashed.png"
)

// Frame duration for walk cycles.
const walkFrameDuration = 200 * time.Millisecond

// PlayerAnimations returns an AnimationData for the player character.
// character.png: 176×20, 16×20 frames, 11 columns, 1 row, no border.
//
// Frame map:
//
//	1      – idle (centered, neutral)
//	2      – idle (directional, facing right)
//	3      – jumping upward
//	4      – falling downward
//	5      – stomp (colliding with enemy from above)
//	6-7    – climb
//	8-9    – (unused)
//	10-11  – walk right
func PlayerAnimations() component.AnimationData {
	img := LoadImage(tilesetCharacter, false) // already has transparency
	grid := ganim8.NewGrid(16, 20, 176, 20)

	anims := map[string]*ganim8.Animation{
		"idle":     ganim8.New(img, grid.GetFrames(1, 1), walkFrameDuration),
		"idle_dir": ganim8.New(img, grid.GetFrames(2, 1), walkFrameDuration),
		"walk":     ganim8.New(img, grid.GetFrames(10, 1, 11, 1), walkFrameDuration),
		"jump":     ganim8.New(img, grid.GetFrames(3, 1), walkFrameDuration, ganim8.PauseAtEnd),
		"fall":     ganim8.New(img, grid.GetFrames(4, 1), walkFrameDuration, ganim8.PauseAtEnd),
		"stomp":    ganim8.New(img, grid.GetFrames(5, 1), walkFrameDuration, ganim8.PauseAtEnd),
		"climb":    ganim8.New(img, grid.GetFrames(6, 1, 7, 1), walkFrameDuration),
	}

	return component.AnimationData{
		Animations: anims,
		Current:    "idle",
	}
}

// EnemyAnimations returns an AnimationData for the given enemy type.
func EnemyAnimations(enemyType component.EnemyType) component.AnimationData {
	switch enemyType {
	case component.EnemyKooper:
		return kooperAnimations()
	case component.EnemyGoomba:
		return goombaAnimations()
	default:
		return kooperAnimations()
	}
}

// kooperAnimations builds animations from koopa_troopa_walk.png (32×24)
// and koopa_troopa_shell.png (32×16). No border, transparent background.
func kooperAnimations() component.AnimationData {
	walkImg := LoadImage(tilesetKoopaWalk, false)
	shellImg := LoadImage(tilesetKoopaShell, false)

	walkGrid := ganim8.NewGrid(16, 24, 32, 24)
	shellGrid := ganim8.NewGrid(16, 16, 32, 16)

	anims := map[string]*ganim8.Animation{
		"walk":  ganim8.New(walkImg, walkGrid.GetFrames(1, 1, 2, 1), walkFrameDuration),
		"shell": ganim8.New(shellImg, shellGrid.GetFrames(1, 1), walkFrameDuration, ganim8.PauseAtEnd),
	}

	return component.AnimationData{
		Animations: anims,
		Current:    "walk",
	}
}

// goombaAnimations builds animations from goomba_walk.png (32×16)
// and goomba_squashed.png (16×8). No border, transparent background.
func goombaAnimations() component.AnimationData {
	walkImg := LoadImage(tilesetGoombaWalk, false)
	deadImg := LoadImage(tilesetGoombaDead, false)

	walkGrid := ganim8.NewGrid(16, 16, 32, 16)
	deadGrid := ganim8.NewGrid(16, 8, 16, 8)

	anims := map[string]*ganim8.Animation{
		"walk":     ganim8.New(walkImg, walkGrid.GetFrames(1, 1, 2, 1), walkFrameDuration),
		"squashed": ganim8.New(deadImg, deadGrid.GetFrames(1, 1), walkFrameDuration, ganim8.PauseAtEnd),
	}

	return component.AnimationData{
		Animations: anims,
		Current:    "walk",
	}
}
