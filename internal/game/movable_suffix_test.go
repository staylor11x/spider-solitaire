package game

import (
	"testing"

	"github.com/staylor11x/spider-solitaire/internal/deck"
	"github.com/stretchr/testify/assert"
)

func TestMovableSuffix(t *testing.T) {
	tests := []struct {
		name string
		pile []CardInPile
		want []deck.Rank
	}{
		{
			name: "empty pile",
			pile: nil,
			want: nil,
		},
		{
			name: "top card face down",
			pile: []CardInPile{
				{Card: deck.Card{Suit: deck.Spades, Rank: deck.Ace}, FaceUp: false},
			},
			want: nil,
		},
		{
			name: "single face-up card",
			pile: []CardInPile{
				{Card: deck.Card{Suit: deck.Spades, Rank: deck.Ace}, FaceUp: true},
			},
			want: []deck.Rank{deck.Ace},
		},
		{
			name: "valid descending order suit run",
			pile: []CardInPile{
				{Card: deck.Card{Suit: deck.Spades, Rank: deck.Four}, FaceUp: true},
				{Card: deck.Card{Suit: deck.Spades, Rank: deck.Three}, FaceUp: true},
				{Card: deck.Card{Suit: deck.Spades, Rank: deck.Two}, FaceUp: true},
				{Card: deck.Card{Suit: deck.Spades, Rank: deck.Ace}, FaceUp: true},
			},
			want: []deck.Rank{
				deck.Four,
				deck.Three,
				deck.Two,
				deck.Ace,
			},
		},
		{
			name: "mixed suit breaks sequence",
			pile: []CardInPile{
				{Card: deck.Card{Suit: deck.Spades, Rank: deck.Three}, FaceUp: true},
				{Card: deck.Card{Suit: deck.Hearts, Rank: deck.Two}, FaceUp: true},
				{Card: deck.Card{Suit: deck.Hearts, Rank: deck.Ace}, FaceUp: true},
			},
			want: []deck.Rank{
				deck.Two,
				deck.Ace,
			},
		},
		{
			name: "rank gap between sequence",
			pile: []CardInPile{
				{Card: deck.Card{Suit: deck.Spades, Rank: deck.Four}, FaceUp: true},
				{Card: deck.Card{Suit: deck.Spades, Rank: deck.Two}, FaceUp: true},
				{Card: deck.Card{Suit: deck.Spades, Rank: deck.Ace}, FaceUp: true},
			},
			want: []deck.Rank{
				deck.Two,
				deck.Ace,
			},
		},
		{
			name: "face-down card stops suffix",
			pile: []CardInPile{
				{Card: deck.Card{Suit: deck.Spades, Rank: deck.Three}, FaceUp: false},
				{Card: deck.Card{Suit: deck.Spades, Rank: deck.Two}, FaceUp: true},
				{Card: deck.Card{Suit: deck.Spades, Rank: deck.Ace}, FaceUp: true},
			},
			want: []deck.Rank{
				deck.Two,
				deck.Ace,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := movableSuffix(tt.pile)

			assert.Equal(t, len(tt.want), len(got))

			for i, card := range got {
				assert.Equal(t, tt.want[i], card.Card.Rank)
			}
		})
	}
}
