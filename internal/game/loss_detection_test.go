package game

import (
	"testing"

	"github.com/staylor11x/spider-solitaire/internal/deck"
	"github.com/stretchr/testify/assert"
)

func TestHasAnyValidMove(t *testing.T) {
	tests := []struct {
		name  string
		piles [TableauPiles]Pile
		want  bool
	}{
		{
			name: "simple valid placement exists",
			piles: func() [TableauPiles]Pile {
				var p [TableauPiles]Pile
				p[0] = newPile(
					makeCardInPile(deck.Spades, deck.Ten, true),
				)
				p[1] = newPile(
					makeCardInPile(deck.Spades, deck.Jack, true),
				)
				return p
			}(),
			want: true,
		},
		{
			name: "king can move to empty pile",
			piles: func() [TableauPiles]Pile {
				var p [TableauPiles]Pile
				p[0] = newPile(
					makeCardInPile(deck.Hearts, deck.King, true),
				)
				//everything else empty
				return p
			}(),
			want: true,
		},
		{
			name: "non king can move to empty pile",
			piles: func() [TableauPiles]Pile {
				var p [TableauPiles]Pile
				p[0] = newPile(
					makeCardInPile(deck.Spades, deck.Ten, true),
					makeCardInPile(deck.Spades, deck.Nine, true),
				)
				return p
			}(),
			want: true,
		},
		{
			name: "valid run buried under invalid card",
			piles: func() [TableauPiles]Pile {
				var p [TableauPiles]Pile
				p[0] = newPile(
					makeCardInPile(deck.Spades, deck.Four, true),
					makeCardInPile(deck.Spades, deck.Three, true),
					makeCardInPile(deck.Spades, deck.Two, true),
					makeCardInPile(deck.Spades, deck.Ace, true),
					makeCardInPile(deck.Hearts, deck.Queen, true), // blocks the run
				)
				// the potentially valid move
				p[1] = newPile(
					makeCardInPile(deck.Clubs, deck.Five, true),
				)
				// fill the rest of the piles
				for i := 2; i < TableauPiles; i++ {
					p[i] = newPile(
						makeCardInPile(deck.Clubs, deck.Ace, true),
					)
				}
				return p
			}(),
			want: false,
		},
		{
			name: "no valid moves - all piles blocked",
			piles: func() [TableauPiles]Pile {
				var p [TableauPiles]Pile
				// All the piles have aces (lowest rank, can't place on anything)
				for i := range TableauPiles {
					p[i] = newPile(
						makeCardInPile(deck.Spades, deck.Ace, true),
					)
				}
				return p
			}(),
			want: false,
		},
		{
			name: "cross-suit stacking allowed",
			piles: func() [TableauPiles]Pile {
				var p [TableauPiles]Pile
				p[0] = newPile(
					makeCardInPile(deck.Hearts, deck.Seven, true),
				)
				p[1] = newPile(
					makeCardInPile(deck.Spades, deck.Eight, true),
				)
				return p
			}(),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasAnyValidMove(tt.piles)
			assert.Equal(t, tt.want, got)
		})
	}
}
