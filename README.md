# Spider Solitaire – Project Summary

## Overview

This project implements a rules-complete Spider Solitaire game engine with a strict separation between game logic and presentation (UI). The engine is deterministic, fully testable, and UI-agnostic, enabling multiple frontends (CLI, Ebiten, etc.).

## Tech Stack 

- Language: Go
- UI (planned): Ebiten
- Testing: Go testing package#
- Architecture Style: Domain-driven, engine-first design

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