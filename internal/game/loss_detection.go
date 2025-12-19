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
	return moving.Suit == target.Suit && moving.Rank+1 == target.Rank
}

// hasAnyValidMove is a method to check if there are any valid moves currently available in the game
func hasAnyValidMove(piles [TableauPiles]Pile) bool {
	// Build a list of movable sequences once
	type movableSeq struct {
		pileIdx    int
		topCard    deck.Card
		canMoveAll bool // true if entire sequence starting with King
	}

	var sequences []movableSeq

	for i, pile := range piles {
		suffix := movableSuffix(pile.Cards())
		if len(suffix) == 0 {
			continue
		}

		sequences = append(sequences, movableSeq{
			pileIdx:    i,
			topCard:    suffix[0].Card,
			canMoveAll: suffix[0].Card.Rank == deck.King,
		})
	}

	// Check if any sequence can be placed somewhere
	for _, seq := range sequences {
		// Kings can go to empty piles (Apparently only kings?)
		if seq.canMoveAll {
			for j, pile := range piles {
				if j != seq.pileIdx && pile.Size() == 0 {
					return true
				}
			}
		}

		// Check if sequence can be placed on another pile
		for j, pile := range piles {
			if j == seq.pileIdx {
				continue
			}

			if pile.Size() == 0 {
				continue // Already checked King moves above
			}

			top, err := pile.TopCard()
			if err != nil || !top.FaceUp {
				continue
			}

			if canPlaceOn(seq.topCard, top.Card) {
				return true
			}
		}
	}

	return false
}
