package ui

import (
	"testing"

	"github.com/staylor11x/spider-solitaire/internal/deck"
	"github.com/staylor11x/spider-solitaire/internal/game"
)

func testCard(rank deck.Rank, suit deck.Suit, faceUp bool) game.CardDTO {
	return game.CardDTO{
		Rank:   game.RankDTO(rank),
		Suit:   game.SuitDTO(suit),
		FaceUp: faceUp,
	}
}

func TestComputeMovableHoverEnd(t *testing.T) {
	tests := []struct {
		name     string
		cards    []game.CardDTO
		startIdx int
		wantEnd  int
	}{
		{
			name:     "invalid negative start",
			cards:    []game.CardDTO{testCard(deck.King, deck.Spades, true)},
			startIdx: -1,
			wantEnd:  -1,
		},
		{
			name:     "invalid start past end",
			cards:    []game.CardDTO{testCard(deck.King, deck.Spades, true)},
			startIdx: 2,
			wantEnd:  -1,
		},
		{
			name: "start card face down",
			cards: []game.CardDTO{
				testCard(deck.King, deck.Spades, false),
				testCard(deck.Queen, deck.Spades, true),
			},
			startIdx: 0,
			wantEnd:  -1,
		},
		{
			name: "single face-up card",
			cards: []game.CardDTO{
				testCard(deck.Ace, deck.Spades, true),
			},
			startIdx: 0,
			wantEnd:  0,
		},
		{
			name: "full valid descending same-suit sequence",
			cards: []game.CardDTO{
				testCard(deck.Four, deck.Spades, true),
				testCard(deck.Three, deck.Spades, true),
				testCard(deck.Two, deck.Spades, true),
				testCard(deck.Ace, deck.Spades, true),
			},
			startIdx: 0,
			wantEnd:  3,
		},
		{
			name: "mid-sequence hover returns 10 and 9 only",
			cards: []game.CardDTO{
				testCard(deck.King, deck.Spades, true),
				testCard(deck.Queen, deck.Spades, true),
				testCard(deck.Jack, deck.Spades, true),
				testCard(deck.Ten, deck.Hearts, true),
				testCard(deck.Nine, deck.Hearts, true),
			},
			startIdx: 3,
			wantEnd:  4,
		},
		{
			name: "blocked above means no movable sequence from hovered card",
			cards: []game.CardDTO{
				testCard(deck.King, deck.Spades, true),
				testCard(deck.Queen, deck.Spades, true),
				testCard(deck.Jack, deck.Spades, true),
				testCard(deck.Ten, deck.Spades, true),
				testCard(deck.Three, deck.Spades, true),
			},
			startIdx: 0,
			wantEnd:  -1,
		},
		{
			name: "suit break blocks sequence to top",
			cards: []game.CardDTO{
				testCard(deck.Four, deck.Spades, true),
				testCard(deck.Three, deck.Hearts, true),
				testCard(deck.Two, deck.Hearts, true),
			},
			startIdx: 0,
			wantEnd:  -1,
		},
		{
			name: "rank gap blocks sequence to top",
			cards: []game.CardDTO{
				testCard(deck.Four, deck.Spades, true),
				testCard(deck.Two, deck.Spades, true),
				testCard(deck.Ace, deck.Spades, true),
			},
			startIdx: 0,
			wantEnd:  -1,
		},
		{
			name: "face-down card above blocks sequence to top",
			cards: []game.CardDTO{
				testCard(deck.Three, deck.Spades, true),
				testCard(deck.Two, deck.Spades, false),
				testCard(deck.Ace, deck.Spades, true),
			},
			startIdx: 0,
			wantEnd:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := computeMovableHoverEnd(tt.cards, tt.startIdx)
			if got != tt.wantEnd {
				t.Fatalf("expected end index %d, got %d", tt.wantEnd, got)
			}
		})
	}
}
