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
