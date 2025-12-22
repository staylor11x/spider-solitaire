package game

// movableSuffix returns the longest movable sequence at the end of a pile.
// In Spider, you can only move sequences that are:
// - All face up
// - Same suit
// - Descending rank
func movableSuffix(cards []CardInPile) []CardInPile {
	if len(cards) == 0 {
		return nil
	}

	// Start from the last card and work backwards
	lastIdx := len(cards) - 1

	// Last card must be face up to be movable
	if !cards[lastIdx].FaceUp {
		return nil
	}

	// Find the longest valid sequence from the end
	startIdx := lastIdx
	for startIdx > 0 {
		current := cards[startIdx]
		prev := cards[startIdx-1]

		// Previous card must be face up
		if !prev.FaceUp {
			break
		}

		// Must be same suit
		if current.Card.Suit != prev.Card.Suit {
			break
		}

		// Must be descending (prev rank = current rank + 1)
		if prev.Card.Rank != current.Card.Rank+1 {
			break
		}

		startIdx--
	}

	return cards[startIdx:]
}
