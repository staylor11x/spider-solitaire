# Copilot Prompt Guidelines – Spider Solitaire (Go + Ebiten)

Project context:
- Overview in [README.md](README.md); roadmap in [ROADMAP.md](ROADMAP.md)
- Engine-first, UI-agnostic design with deterministic state transitions
- Primary entrypoint: [cmd/game/main.go](cmd/game/main.go)
- Core engine: [internal/game](internal/game), deck: [internal/deck](internal/deck)

Environment & targets:
- Go version: 1.25 (update `go.mod` when upgrading)
- Current OS target: Windows
- Future OS targets: Windows, Linux, Android (via Gomobile for Ebiten)

Principles:
- Efficient, enterprise-quality, idiomatic Go
- Do not over-engineer; prefer small, composable units
- Design for testability; validate-before-mutate; deterministic operations
- Test behavior, not small implementation details; keep tests resilient
- Maintain strict separation of engine and UI (Ebiten consumes DTOs)

Coding standards:
- Accept interfaces where helpful; return concrete types; avoid needless abstractions
- Use clear naming and small functions; document exported symbols (focus on “why”)
- Defensive copying at boundaries; don’t expose internal mutable slices
- Error handling: prefer sentinel/typed errors in [internal/game/errors.go](internal/game/errors.go); wrap with `%w`
- Keep moves atomic; revert on invalid execution

Testing standards:
- Table-driven tests for move validation, movable suffix, completed runs, win/loss
- Integration-style tests for end-to-end scenarios via `GameState`

UI integration (Ebiten):
- Provide read-only DTOs: CardDTO, PileDTO, GameViewDTO
- UI reads snapshots; engine never depends on Ebiten
- Keep rendering/input logic in `cmd/game` or `internal/ui`

Performance guidance:
- Minimize allocations in hot paths (move detection, suffix computation)
- Measure with benchmarks before optimizing

Build & run (Windows focus):
- Run tests: 
  - `go test ./...`
- Build engine/CLI:
  - `go build ./cmd/game`
- Linux builds: ensure fonts/assets and path separators are handled
- Android: plan for Gomobile/Ebiten build flow later

Prompt patterns:
- Propose the smallest viable change; list affected files; include code + tests
- Call out trade-offs and choose the simplest workable path
- Ask targeted clarifying questions when requirements are ambiguous

Acceptance checklist per change:
- Idiomatic, readable, deterministic
- Tests added/updated for behavior; `go test ./...` passes
- No UI assumptions inside engine; DTOs for UI
- Minimal allocations in critical paths; no needless complexity
- Update docs when behavior changes (README/CHANGELOG if requested)
