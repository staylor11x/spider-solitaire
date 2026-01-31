# CHANGELOG

### v.1.7.4 - Improving Card Selection UX

- User can no longer select a column and has to actually click on the pile to select the card.
- Added ESC key usage to cancel move.

### v1.7.3 - Move Undo Functionality

Implemented comprehensive undo system for move history tracking and recovery:

**Backend Engine:**
- Added GameState.snapshot() → deep copies entire game state for history preservation
- Added GameState.pushHistory() → saves snapshot after validation with FIFO max 25 moves
- Added GameState.Undo() → restores previous GameState from history
- Added Pile.Clone() → deep copy constructor for defensive undo snapshots
- Added ErrNoHistory error type for undo edge cases

**UI Integration:**
- Keyboard shortcut: [U] - Undo Move
- Selection automatically clears after undo for clean UX
- Error display when no undo history available
- Help overlay updated with undo instruction

**Testing:**
- 9 comprehensive test cases covering:
  - Empty history error handling
  - Snapshot isolation and deep copy verification
  - Undo after deal row operations
  - Undo after card moves
  - Undo after run auto-completion
  - Chained undo operations
  - History FIFO enforcement at 25 moves
  - State flag restoration (Won/Lost)

Quality-of-life improvement enabling players to recover from mistakes without resetting the game.

### v0.1.6 - Improved Errors handling

Implemented improved error handling in the game package
- Refactored most of the errors in the game package to use sentinel errors or typed errors.

### v0.1.5 - Implemented Auto Versioning

Implemented auto versioning into repo
- Version will increase every patch version on merge/push to main branch.


### v0.0.4 – Moving Sequences Between Piles

Added MoveSequence(srcIdx, startIdx, dstIdx) on GameState.

Validations:
- Source/destination pile indices valid.
- Start index valid and points to a face-up card.
- Sequence must be descending and same suit.
- Destination pile must be able to accept the sequence.
- Moves are atomic: if move fails, state restored.

Added Pile.RemoveCardsFrom(), Pile.AddCards(), Pile.FlipTopCardIfFaceDown().
Added internal helpers: isValidSequence(), sequenceEqual().

Unit tests:
- Valid moves succeed.
- Invalid moves fail (wrong suit, not descending, destination invalid, face-down card).
- Verified top card flips after move if previously face-down.

### v0.0.3 – Stock & Row Dealing

Added DealRow() → deals one face-up card to each tableau pile from stock.

Added CanDealRow() → check if stock has at least 10 cards.

Unit tests: row dealing reduces stock, ensures face-up cards, fails gracefully if not enough stock.

### v0.0.2 – Tableau & GameState Initialization

Added Pile struct:

- AddCard(), TopCard(), GetCards(), Size().
- Added CardInPile wrapper to track FaceUp state.
- Added Tableau (10 piles).
- Added GameState containing Tableau and Stock.
- Added DealInitialGame() → deals 54 cards to tableau (correct 6/5 split, only last card face-up), leaves 50 in stock.

Unit tests: validated initial setup, face-up rules.

### v0.0.1 – Initial Deck Support

Added Card, Suit, and Rank types with String() methods.

Implemented Deck with:

- NewStandardDeck() → 52-card deck.
- NewMultiDeck(n) → multiple decks combined.
- Shuffle() → randomize order.
- Draw() and DrawAll() → retrieve cards.
- Size() → count remaining cards.









