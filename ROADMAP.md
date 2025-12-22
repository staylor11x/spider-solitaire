# spider-solitaire
Basic implementation of a spider solitaire game, designed for Ubuntu Linux


## 🛠️ Roadmap – Features Needed to Complete a Playable Game

### 🔹 Phase 1 – Technical Cleanup

- Refactor game package into smaller files: pile.go, tableau.go, gamestate.go.

- Introduce structured error types (e.g., ErrInvalidMove, ErrNotEnoughCards).

- Optional: Introduce structured logging (zap or zerolog) for debugging moves/deals.

### 🔹 Phase 2 – Core Gameplay Completion

Run Completion Detection:

- Detect when a full suit run from King → Ace exists in a pile.

- Automatically remove the run from the tableau.

Win/Loss Conditions:

- Win when all runs (8 total in 2-deck Spider) are completed.

- Lose when no valid moves and no stock left.

Move Undo (stretch goal for quality-of-life):

- Keep history of moves and allow undo (useful for testing too).

### 🔹 Phase 3 – UI Integration Prep

- Define rendering-friendly API:

- Methods to expose tableau/stock state in a form the UI can easily consume ([]CardDTO).

- Preserve face-up/face-down distinction.

- Add basic CLI printer for debugging before UI.

- Start Ebiten integration:

- Draw piles and cards.

- Handle simple interactions (click/drag or keyboard simulation).


## Current Status

- Game rules fully implemented
- Win and loss conditions complete
- Engine stable and test-covered
- Ready for UI integration via Ebiten