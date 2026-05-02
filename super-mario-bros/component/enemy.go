package component

import "github.com/yohamta/donburi"

// EnemyType identifies the kind of enemy. Matches LDtk entity identifiers.
type EnemyType string

const (
	EnemyKooper EnemyType = "Kooper"
	EnemyGoomba EnemyType = "Goomba"
)

// EnemyState tracks the lifecycle of an enemy.
type EnemyState int

const (
	EnemyAlive EnemyState = iota
	EnemyShell            // Kooper only: shell state, can be stomped again
	EnemyDead             // death animation playing, awaiting cleanup
)

// EnemyData marks an entity as an enemy and stores its type and state.
type EnemyData struct {
	Type         EnemyType
	State        EnemyState
	CleanupTimer float64 // seconds remaining before entity removal (0 = no timer)
}

var Enemy = donburi.NewComponentType[EnemyData]()
