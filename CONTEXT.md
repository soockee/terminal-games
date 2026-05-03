# Architecture Context

## Stack

Go · Ebiten (game engine) · Donburi (ECS) · Resolv (physics, optional) · ldtk-super-simple-loader (level loading, optional) · [NATS](https://nats.io/) (networking / net code)

## Core Values

1. **Discoverability** — predictable structure; you always know where to look.
2. **Low import friction** — minimize packages to avoid circular dependencies and import tangles.
3. **Resist premature abstraction** — a new folder or package must earn its existence. A file should outgrow its folder before you split.
4. **Simplicity of intent** — every feature and line has a clear *why*. No speculative code, no unused features. The code should communicate its intent, not just its mechanics.
5. **Minimal dependencies** — only depend on external libraries when the domain genuinely demands expertise (game engine, ECS framework, physics). Otherwise, own the code.

## Repo Structure

Monorepo of standalone games. Each game is an independent Go module. No shared code between games — duplication is acceptable. `game-template/` is the starting point for new games.

## Canonical Game Structure

```
game-name/
  main.go        — trivial bootstrap (parse flags, call game.Run())
  game/          — thin orchestration: init ECS world, register systems, implement ebiten.Game
  archetype/     — entity constructors: create entity + set initial values (one file per entity type)
  component/     — data types + zero-size tag components (one file per component)
  system/        — game logic, named by verb/concern (one file per system)
  event/         — world-level pub/sub event types (one file per event)
  assets/        — raw files (.png, .wav, .ldtk) + embed declarations, no logic
```

## Design Rules

### Packages

Start with the canonical folders. Additional packages may exist only when they wrap an external library boundary (e.g., a `resolv/` physics wrapper) or solve a concern entirely outside the game loop (e.g., a `network/` package handling WebSocket connections and message serialization for multiplayer — the transport lives in its own package, but game state sync happens through components and systems). A package is never created for mere organizational comfort.

### Components & Tags

Components are data types. Tags are zero-size components. If it has data, it's a component. If it's just a query filter, it's a tag. Both are defined with `donburi.NewComponentType`.

### Archetypes

Archetypes are entity constructors. They create an entity, attach components, and set initial values. They return a `*donburi.Entry`. They are not containers or structs — they are factory functions.

### Systems

Systems are named by what they *do* (verb), not what they act on (noun). `movement.go`, `collision.go`, `render.go` — not `ball.go`, `player.go`. If a system grows too large, split by sub-concern, not by entity type.

### State & Levels

Game state lives as ECS data (components on entities), not as control flow in a scene manager. Levels are spatial regions of the same ECS world. Transitions are camera/viewport changes, not world resets. Special-purpose screens (menu, end screen) are just levels whose entities respond to different input.


### LDtk

LDtk defines the authored world, not the living world. Use it for level layout, static tiles, IntGrid semantics, entity placement, entity size, custom fields, spawn points, triggers, doors, camera zones, pickups, and authored paths.

LDtk data is imported into ECS once. The importer converts levels, tiles, IntGrid cells, and entity instances into game-owned data and archetype calls. It is not a gameplay system.

Entity visuals in LDtk are editor previews. Runtime sprites, spritesheets, animation clips, frame timing, current animation state, and rendering behavior belong to Go.

Paths in LDtk describe authored movement intent: points, mode, pauses, speed overrides, and initial direction. Runtime path state — current target, wait timers, interruption, collision-aware movement, and animation changes — belongs to ECS components and systems.


### Networking

All networking and net code must use [NATS](https://nats.io/). No other messaging or transport layer (raw TCP/UDP, WebSockets, gRPC, etc.) should be introduced. NATS handles pub/sub, request/reply, and any multiplayer communication. Transport logic lives in its own package (e.g., `network/`); game state sync flows through components and systems.

## Implementation Conventions

- `Update()` mutates state; `Draw()` only renders. Never mutate in Draw.
- Ebiten runs at fixed TPS. Don't use delta-time for game logic.
- Minimize allocations in the hot loop — pre-allocate, reuse `DrawImageOptions`.
- Systems are stateless — they read/write components, they don't hold their own state.
- Load assets once at startup; store as `*ebiten.Image` / `[]byte`.
- LDtk is load-time data. Normalize positions during import; systems should not depend on LDtk-specific structures.
