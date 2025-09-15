package testtools

import (
	"github.com/staylor11x/spider-solitaire/internal/deck"
	"github.com/staylor11x/spider-solitaire/internal/game"
)

func NewPile(cards ...game.CardInPile) game.Pile {
	var p game.Pile
	p.AddCards(cards)
	return p
}

func NewSequence(s deck.Suit, ranks []deck.Rank) []game.CardInPile {
	seq := make([]game.CardInPile, 0, len(ranks))
	for _, r := range ranks {
		seq = append(seq, game.CardInPile{Card: deck.Card{Suit: s, Rank: r}, FaceUp: true})
	}
	return seq
}

// helper to build a completed run, n parameter is how many
func buildCompleteRun(n int) []game.CardInPile {
	run := make([]game.CardInPile, 13)
	for r := deck.King; r >= deck.Two; r-- {
		run = append(run, game.CardInPile{
			Card:   deck.Card{Suit: deck.Spades, Rank: r},
			FaceUp: true,
		})
	}
	return run
}

// NewSequenceWithIgnoreRank is a method that can be used to build a sequence with a card missing
func NewSequenceWithIgnoreRank(s deck.Suit, rankToIgnore deck.Rank) []game.CardInPile {
	seq := make([]game.CardInPile, 0, 13)
	for r := deck.King; r >= deck.Ace; r-- {
		if r == rankToIgnore {
			continue
		}
		seq = append(seq, game.CardInPile{
			Card:   deck.Card{Suit: s, Rank: r},
			FaceUp: true,
		})
	}
	return seq
}
