package game

import (
	"github.com/staylor11x/spider-solitaire/internal/deck"
)

// hadEmptyPile is a method to check if there are any empty piles in the tableau
func hasEmptyPile(piles [TableauPiles]Pile) bool {
	for _, p := range piles {
		if p.Size() == 0 {
			return true
		}
	}
	return false
}

// hasStock is a method to check if there are any card left in the stock
func hasStock(stock []deck.Card) bool {
	return (len(stock) > 0)
}

// canPlaceOn is a method to check if the moving card can be placed on the target
func canPlaceOn(moving deck.Card, target deck.Card) bool {
	return moving.Rank+1 == target.Rank
}

// hasAnyValidMove is a method to check if there are any valid moves currently available in the game
// This includes:
// - Moving any card/sequence to an empty pile
// - Moving any partial or full sequence to stack on another pile
// - Cross-suit stacking (any card one rank higher)
// TODO: Check the cyclomatic complexity of this function, is this clean code??
func hasAnyValidMove(piles [TableauPiles]Pile) bool {
	for i, pile := range piles {
		cards := pile.Cards()
		if len(cards) == 0 {
			continue
		}

		// Get the movable suffix (face-up, descending, same-suit run from the top)
		suffix := movableSuffix(cards)
		if len(suffix) == 0 {
			continue
		}

		// Check if any sequence can be placed somewhere
		// suffix[0] is the deepest movable card (e.g. 9S)
		// suffix[len-1] is the top card (e.g. 7S)
		// We can move: full (9-8-7), partial (8-7), or just top (7)
		for startIdx := range suffix {
			movingCard := suffix[startIdx].Card

			for j, targetPile := range piles {
				if j == i {
					continue
				}

				// Any sequence can move to an empty pile
				if targetPile.Size() == 0 {
					return true
				}

				// Check if this sub-sequence can stack on the target pile
				top, err := targetPile.TopCard()
				if err != nil || !top.FaceUp {
					continue
				}

				if canPlaceOn(movingCard, top.Card) {
					return true
				}
			}
		}
	}
	return false
}
