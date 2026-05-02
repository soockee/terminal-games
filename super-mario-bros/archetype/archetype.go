package archetype

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/jakecoffman/cp/v2"
	ldtkgo "github.com/soockee/ldtk-super-simple-loader"
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/yohamta/donburi"
)

// SpawnEntities reads all entity instances from the LDtk level and creates
// corresponding ECS entities in the donburi world.
// Only entities tagged "collidable" in LDtk receive static physics shapes.
// Enemies and the player get their own dynamic bodies from their spawn functions.
func SpawnEntities(w donburi.World, level *ldtkgo.Level) {
	// Build an IID → Entity lookup for resolving entity references.
	iidMap := make(map[string]*ldtkgo.Entity)
	for _, ent := range level.AllEntities() {
		iidMap[ent.IID] = ent
	}

	if player := level.Entity("Player"); player != nil {
		NewPlayer(w, player)
	}
	for _, ent := range level.EntitiesByID("Kooper") {
		NewKooper(w, ent, iidMap)
	}
	for _, ent := range level.EntitiesByID("Goomba") {
		NewGoomba(w, ent, iidMap)
	}
	for _, ent := range level.EntitiesByID("Block") {
		NewBlock(w, ent)
	}
	for _, ent := range level.EntitiesByID("JumpBlock") {
		NewJumpBlock(w, ent)
	}
	for _, ent := range level.EntitiesByID("Ground") {
		NewGround(w, ent)
	}
	for _, ent := range level.EntitiesByID("Gem") {
		NewGem(w, ent)
	}

	// Merge horizontally adjacent collidable tiles into wider shapes to
	// eliminate ghost collisions at tile seams.
	addMergedStaticCollision(w, level.EntitiesByTag("collidable"))
}

// NewDebug creates the singleton debug toggle entity.
func NewDebug(w donburi.World) {
	e := w.Entry(w.Create(component.Debug))
	component.Debug.Set(e, &component.DebugData{Enabled: false})
}

// NewCamera creates the singleton camera entity with the given virtual screen scale.
func NewCamera(w donburi.World, scaleX, scaleY float64) {
	e := w.Entry(w.Create(component.Camera))
	component.Camera.Set(e, &component.CameraData{ScaleX: scaleX, ScaleY: scaleY})
}

// NewScore creates the singleton score entity with a target for win detection.
func NewScore(w donburi.World, target int) {
	e := w.Entry(w.Create(component.Score))
	component.Score.Set(e, &component.ScoreData{Target: target})
}

// NewGameState creates the singleton game state entity.
func NewGameState(w donburi.World) {
	e := w.Entry(w.Create(component.GameState))
	component.GameState.Set(e, &component.GameStateData{})
}

// NewAudio creates the singleton audio entity with pre-decoded players.
func NewAudio(w donburi.World, ctx *audio.Context, bgm *audio.Player, sfx map[string]*audio.Player) {
	e := w.Entry(w.Create(component.Audio))
	component.Audio.Set(e, &component.AudioData{Ctx: ctx, BGMusic: bgm, SFX: sfx})
}

// NewPhysicsSpace creates the singleton physics space with the given gravity.
func NewPhysicsSpace(w donburi.World, gravity cp.Vector) *cp.Space {
	space := cp.NewSpace()
	space.SetGravity(gravity)
	space.SetDamping(1.0)
	space.Iterations = 20 // more iterations for robust overlap resolution

	e := w.Entry(w.Create(component.PhysicsSpace))
	component.PhysicsSpace.Set(e, &component.PhysicsSpaceData{Space: space})
	return space
}
