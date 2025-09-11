# CHANGELOG

### v0.0.1 – Initial Deck Support

Added Card, Suit, and Rank types with String() methods.

Implemented Deck with:

- NewStandardDeck() → 52-card deck.
- NewMultiDeck(n) → multiple decks combined.
- Shuffle() → randomize order.
- Draw() and DrawAll() → retrieve cards.
- Size() → count remaining cards.

### v0.0.2 – Tableau & GameState Initialization

Added Pile struct:

- AddCard(), TopCard(), GetCards(), Size().
- Added CardInPile wrapper to track FaceUp state.
- Added Tableau (10 piles).
- Added GameState containing Tableau and Stock.
- Added DealInitialGame() → deals 54 cards to tableau (correct 6/5 split, only last card face-up), leaves 50 in stock.

Unit tests: validated initial setup, face-up rules.

### v0.0.3 – Stock & Row Dealing

Added DealRow() → deals one face-up card to each tableau pile from stock.

Added CanDealRow() → check if stock has at least 10 cards.

Unit tests: row dealing reduces stock, ensures face-up cards, fails gracefully if not enough stock.

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

### v0.1.5 - Implemented Auto Versioning

Implemented auto versioning into repo
- Version will increase every patch version on merge/push to main branch.

### v0.1.6 - Improved Errors handling

Implemented improved error handling in the game package
- Refactored most of the errors in the game package to use sentinel errors or typed errors.
