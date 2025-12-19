package game

// movableSuffix is a method to return the sequence of cards in a pile that can potentially be moved starting from the bottom and traversing the pile upwards.
func movableSuffix(cards []CardInPile) []CardInPile {
	if len(cards) == 0 {
		return nil
	}

	// start from the top of the pile
	top := cards[len(cards)-1]

	// i don't believe that this will ever happen
	if !top.FaceUp {
		return nil
	}

	suffix := []CardInPile{top}

	for i := len(cards) - 2; i >= 0; i-- {
		curr := cards[i]
		prev := suffix[0]

		if !curr.FaceUp {
			break
		}

		if curr.Card.Suit != prev.Card.Suit {
			break
		}

		if curr.Card.Rank != prev.Card.Rank+1 {
			break
		}

		// prepend
		suffix = append([]CardInPile{curr}, suffix...)
	}

	return suffix
}
