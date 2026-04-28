# Match-3

A Match-3 puzzle game built with Ebitengine + Donburi ECS, with levels designed in LDtk.

## Design Overview

### Architecture

- **Hybrid board:** Singleton `Board` component owns grid structure + per-cell Donburi entity pointers
- **LDtk-driven:** IntGrid defines board shape; level custom fields for config
- **Variable dimensions:** Board size/shape derived from IntGrid per level

### Components

| Component | Scope | Fields |
|-----------|-------|--------|
| `Board` | Singleton | `Cols, Rows, CellType [][]int, Cells [][]*donburi.Entry, Phase, SelectedCol, SelectedRow, SwapA, SwapB, ChainDepth, NumColors, ScoreTarget, TimeLimit, TimeRemaining, OffsetX, OffsetY, TileSize` |
| `GridPos` | Per-tile | `Col, Row int` |
| `GemType` | Per-tile | `Color int` (0–5) |
| `PixelPos` | Per-tile | `X, Y float64` |
| `Tween` | Per-tile | `StartX, StartY, EndX, EndY, Elapsed, Duration float64, Active bool, Ease EaseFunc` |
| `Sprite` | Per-tile | `*ebiten.Image` |
| `Score` | Singleton | `Value, Target int` |
| `GameState` | Singleton | `Dead, Started, Restart, Won, Paused bool` |
| `Camera` | Singleton | `X, ScaleX, ScaleY float64` |
| `Debug` | Singleton | `Enabled bool` |
| `Audio` | Singleton | `Ctx, BGMusic, SFX` |

### Tags

- `Tile` — applied to all gem entities

### Phase State Machine

```
Idle → Selected → Swapping → Checking → Removing → Falling → Refilling → Checking → ... → Idle
                                    ↘ Reversing → Idle (no match found)
```

### Systems (update order)

1. `UpdateInput` — tap→grid coord, selection, swap initiation
2. `UpdateBoard` — state machine driver with per-phase helpers
3. `UpdateTween` — advance all active tweens, interpolate PixelPos
4. `UpdateScore` — check score target, timer countdown
5. `ProcessEvents` — drain audio event queue

### Renderers (draw order)

1. `DrawEntities` — tile sprites at PixelPos + selection highlight
2. `DrawScore` — score/target/timer/chain multiplier
3. `DrawHUD` — start/pause/win overlays
4. `DrawDebug` — grid lines, phase name, FPS

### Gameplay Rules

- **Input:** Direct tap/click selection, keyboard arrow+enter fallback
- **Swap:** Lenient — animate swap, reverse if no match
- **Matching:** 3+ same color horizontally or vertically, separate matches with dedup
- **Scoring:** 10 pts/tile × chain multiplier (resets on return to Idle)
- **Win:** Reach score target
- **Loss:** None for v1; reshuffle board on deadlock
- **Initial fill:** Regenerate until match-free
- **Gravity:** Per-column, blocked cells (IntGrid value 2) act as barriers
- **Colors:** 6 default, configurable per level (up to 8)

### Tween Timing

| Animation | Duration | Easing |
|-----------|----------|--------|
| Swap | 0.15s | Ease-out quad |
| Reverse | 0.12s | Linear |
| Fall | 0.08s × distance | Ease-out quad |
| Spawn | 0.12s | Ease-out quad |
| Remove (fade) | 0.15s | Linear |

### Tileset

- **File:** `assets/tilesets/match_3_art.png` (96×256, 3 cols × 8 rows of 32×32 tiles)
- **Mapping:** Color 0–5 → rows 0–5, column 0 (first variant per row)
- **Reserved:** Rows 6–7 for future 7th/8th colors on harder levels

---

## Implementation Plan

### Phase 1: LDtk Level Setup

These steps are performed in the LDtk editor (`assets/ldtk/match-3.ldtk`):

1. **Create level custom fields** (Project Settings → Level fields):
   - `NumColors` — Int, default `6`
   - `ScoreTarget` — Int, default `1000`
   - `TimeLimit` — Int, default `0` (0 = unlimited, value in seconds)

2. **Create IntGrid layer** named `Board`:
   - Cell size: 32×32
   - Define values:
     - `0` — Empty (default, not part of board)
     - `1` — Playable cell
     - `2` — Blocked / obstacle
   - Assign colors in LDtk for visual clarity (e.g., green=playable, red=blocked, transparent=empty)

3. **Design Level 0** (first playable level):
   - Paint an 8×8 square of value `1` cells in the `Board` IntGrid layer
   - Set level custom fields: `NumColors=6`, `ScoreTarget=1000`, `TimeLimit=0`
   - Level pixel size should accommodate the board (e.g., 256×256 for 8×8 at 32px)

4. **Design Level 1** (second level — introduce obstacles):
   - Paint an 8×8 grid with some `2` (blocked) cells in the middle
   - Set `NumColors=6`, `ScoreTarget=1500`, `TimeLimit=0`

5. **(Optional) Design a win screen level:**
   - No IntGrid — just a background and entity for "congratulations" display

6. **Add the tileset** to the LDtk project:
   - Import `match_3_art.png` as a tileset (32×32 tile size)
   - This is for LDtk preview only — rendering in-game uses code-based quad slicing

7. **Save and export** simplified data (Project Settings → enable "Save levels to separate files" if desired)

### Phase 2: Module Setup

1. Rename `go.mod` module from `match-3` to `match-3`
2. Update all import paths from `match-3` to `match-3`
3. Verify builds: `go build ./...`

### Phase 3: Components

1. Replace `component/component.go` with match-3 specific components:
   - Remove: `Shape`, `SpawnPos`, `Color` (the old fallback color)
   - Add: `Board`, `GridPos`, `GemType`, `PixelPos`, `Tween`
   - Keep: `Sprite`, `Score`, `GameState` (add `Paused`), `Camera`, `Debug`, `Audio`
2. Define `Phase` enum and `EaseFunc` enum in the component package
3. Define `GemQuads` lookup table

### Phase 4: Tags

1. Replace `tags/tags.go`:
   - Remove: `Player`, `Collidable`, `Ground`
   - Add: `Tile`

### Phase 5: Board Initialization (game/game.go)

1. In `Build()`, after loading the LDtk level:
   - Read `Board` IntGrid layer → derive `Cols`, `Rows`, `CellType`
   - Read level custom fields → `NumColors`, `ScoreTarget`, `TimeLimit`
   - Create `Board` singleton entity
   - For each `CellType[c][r] == 1` cell: create a tile entity with `GridPos`, `GemType` (random, match-free), `PixelPos`, `Sprite`, `Tween`
   - Store entity pointers in `Board.Cells[c][r]`

### Phase 6: Systems — Input

1. Implement `UpdateInput`:
   - Map click/touch position → grid col/row (using Board.OffsetX/Y and TileSize)
   - Phase `Idle` + valid cell tapped → set `SelectedCol/Row`, transition to `Selected`
   - Phase `Selected` + adjacent cell tapped → set `SwapA/B`, start swap tweens, transition to `Swapping`
   - Phase `Selected` + non-adjacent tapped → update selection
   - Phase `Selected` + same cell tapped → deselect, back to `Idle`
   - Keyboard: arrow keys move cursor, Enter/Space selects

### Phase 7: Systems — Tween

1. Implement `UpdateTween`:
   - Query all entities with `Tween` component where `Active == true`
   - Advance `Elapsed` by `1.0/60.0` (or use Ebitengine's TPS)
   - Compute `t = Elapsed / Duration`, clamp to 1.0
   - Apply easing function to `t`
   - Interpolate `PixelPos.X` and `PixelPos.Y` from Start to End
   - If `t >= 1.0`: set `Active = false`, snap to End position

### Phase 8: Systems — Board Logic

1. Implement `UpdateBoard` with a switch on `Phase`:
   - **Swapping:** Check if both swap tweens are done → transition to `Checking`
   - **Checking:** Call `checkMatches()` → if matches found, mark tiles, increment `ChainDepth`, transition to `Removing`. If no matches and came from `Swapping`, transition to `Reversing`. If no matches and came from `Refilling`, reset `ChainDepth`, transition to `Idle`.
   - **Removing:** Start fade tweens on matched tiles. When done, destroy entities, set `Cells[c][r] = nil`, transition to `Falling`
   - **Falling:** Call `applyGravity()` — for each column, shift tiles down into gaps, start fall tweens. When all done → transition to `Refilling`
   - **Refilling:** Call `refillBoard()` — for each column, count empty cells at top, spawn new tile entities above board, start fall tweens. When all done → transition to `Checking` (chain detection)
   - **Reversing:** Start reverse tweens on SwapA/B. When done, swap grid data back, transition to `Idle`

2. Implement helpers:
   - `checkMatches() map[[2]int]bool` — horizontal + vertical scan
   - `applyGravity()` — per-column bottom-up scan, respects `CellType == 2` barriers
   - `refillBoard()` — per-column top-fill with random match-free colors
   - `hasValidMoves() bool` — check if any swap would produce a match; if not, reshuffle

### Phase 9: Systems — Score

1. Implement `UpdateScore`:
   - Decrement `TimeRemaining` if `TimeLimit > 0` and game is active
   - Check `TimeRemaining <= 0` → game over (future: for now just freeze)
   - Check `Score.Value >= Score.Target` → set `GameState.Won`

### Phase 10: Renderers

1. `DrawEntities`: iterate all `Tile`-tagged entities, draw `Sprite` at `PixelPos`
2. `DrawScore`: show score/target, chain multiplier popup, timer if active
3. `DrawHUD`: start prompt, win overlay, pause screen
4. `DrawDebug`: grid lines for `Playable` cells, current phase name, FPS

### Phase 11: Polish & Wiring

1. Wire systems and renderers in `main.go` GameConfig
2. Add selection highlight rendering (border/glow on selected tile)
3. Add audio events for match, swap, chain
4. Add pause button (reuse flappy-bird pattern)
5. Test with multiple LDtk levels and level switching

### Phase 12: Build & Release

1. Add `.github/workflows/match-3.yml` (mirror flappy-bird.yml, tag format `v*-match3`)
2. Add WASM build with `wasm_exec.js` bundled
3. Test in browser
