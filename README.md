# Spider Solitaire – Project Summary

## Overview

This project implements a rules-complete Spider Solitaire game engine with a strict separation between game logic and presentation (UI). The engine is deterministic, fully testable, and UI-agnostic, enabling multiple frontends (CLI, Ebiten, etc.).

## Tech Stack 

- Language: Go
- UI (planned): Ebiten
- Testing: Go testing package
- Architecture Style: Domain-driven, engine-first design

## Installation

### Download

Download the latest release for your platform from the [Releases page](https://github.com/staylor11x/spider-solitaire/releases/latest):
- **Windows (64-bit):** `spider-solitaire-vX.X.X-windows-amd64.exe`
- **Linux (64-bit):** `spider-solitaire-vX.X.X-linux-amd64`

### Verify Download (Recommended)

To verify the integrity of your download:

**Windows (PowerShell):**
```powershell
# Download checksums.txt from the release
# Then run:
Get-FileHash spider-solitaire-v1.5.0-windows-amd64.exe -Algorithm SHA256

# Compare the output hash with the one in checksums.txt
```

**Linux:**
```bash
# Download checksums.txt and the binary
# Then run:
sha256sum -c checksums.txt

# Should output: spider-solitaire-v1.5.0-linux-amd64: OK
```

### Running the Game

**Windows:**
```cmd
.\spider-solitaire-v1.5.0-windows-amd64.exe
```

**Linux:**
```bash
chmod +x spider-solitaire-v1.5.0-linux-amd64
./spider-solitaire-v1.5.0-linux-amd64
```

**Note for Windows users:** You may see a "Windows protected your PC" warning because this software is not code-signed. This is expected for open-source software. Click "More info" → "Run anyway" to proceed.

## Core Design Principles

- Pure game engine (no UI assumptions)
- Deterministic state transitions
- Explicit validation before mutation
- Defensive copying for safety
- Strong test coverage using table-driven tests

## Game Engine Structure

### Key Types

- GameState
  - Owns full game state
    - Tableau, Stock, Completed runs
    - Win/Loss flags
- Tableau → fixed array of Pile
- Pile → ordered stack of CardInPile
- CardInPile
  - Card (Rank, Suit)
  - FaceUp flag

### Core Responsibilities

- Deal initial game
- Deal rows from stock
- Move card sequences atomically
- Flip cards when exposed
- Detect completed runs
- Detect win condition (8 completed runs)
- Detect loss condition (no stock, no empty piles, no valid moves)

### Move & Rule Validation

- All moves are validated before execution:
  - Indices
  - Face-up cards only
  - Proper descending, same-suit sequences
  - Destination acceptance rules
- Moves are atomic: invalid execution restores state

### Loss Detection Logic

A loss occurs when:
1. Stock is empty
2. No empty tableau piles
3. No valid moves exist

Move detection uses:
- A movable suffix helper
- Detection of partial and full sequence moves
- King-only moves to empty piles

## Testing Strategy 

Table-driven unit tests for:
- Move validation
- Completed runs
- Valid move detection

Integration-style tests for:
- Win condition
- Loss condition

Tests focus on behavior, not implementation details

## UI Integration Strategy

- Introduce DTO projection layer:
  - CardDTO, PileDTO, GameViewDTO
- UI consumes read-only snapshots
- No UI code depends on engine internals

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for:
- Development workflow
- Branching and versioning strategy
- How to create releases
- Testing and code quality guidelines
