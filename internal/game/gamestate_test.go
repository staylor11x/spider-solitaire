package game

import (
	"testing"

	"github.com/staylor11x/spider-solitaire/internal/deck"
	"github.com/stretchr/testify/assert"
)

// Helper Functions

func newPile(cards ...CardInPile) Pile {
	var p Pile
	p.AddCards(cards)
	return p
}

func makeCardInPile(s deck.Suit, r deck.Rank, faceUp bool) CardInPile {
	return struct {
		Card   deck.Card
		FaceUp bool
	}{
		Card:   deck.Card{Suit: s, Rank: r},
		FaceUp: faceUp,
	}
}

// newSequence is a method that can be used to build a full completed sequence
func newSequence(s deck.Suit) []CardInPile {
	seq := make([]CardInPile, 0, 13)
	for r := deck.King; r >= deck.Ace; r-- {
		seq = append(seq, CardInPile{
			Card:   deck.Card{Suit: s, Rank: r},
			FaceUp: true,
		})
	}
	return seq
}

// newSequenceWithIgnoreRank is a method that can be used to build a sequence with a card missing
func newSequenceWithIgnoreRank(s deck.Suit, rankToIgnore deck.Rank) []CardInPile {
	seq := make([]CardInPile, 0, 13)
	for r := deck.King; r >= deck.Ace; r-- {
		if r == rankToIgnore {
			continue
		}
		seq = append(seq, CardInPile{
			Card:   deck.Card{Suit: s, Rank: r},
			FaceUp: true,
		})
	}
	return seq
}

// Unit Tests

// Note: Do I care that this helper works, or do I care that moves and run detection work?
func TestIsValidSequence(t *testing.T) {
	tests := []struct {
		name   string
		cards  []CardInPile
		expect bool
	}{
		{
			"Valid descending run",
			[]CardInPile{
				makeCardInPile(deck.Spades, deck.King, true),
				makeCardInPile(deck.Spades, deck.Queen, true),
				makeCardInPile(deck.Spades, deck.Jack, true),
			},
			true,
		},
		{
			"Wrong suit",
			[]CardInPile{
				makeCardInPile(deck.Spades, deck.King, true),
				makeCardInPile(deck.Diamonds, deck.Queen, true),
			},
			false,
		},
		{
			"Not descending",
			[]CardInPile{
				makeCardInPile(deck.Spades, deck.King, true),
				makeCardInPile(deck.Spades, deck.Ace, true),
			},
			false,
		},
		{
			"Edge case: Ace King",
			[]CardInPile{
				makeCardInPile(deck.Clubs, deck.Ace, true),
				makeCardInPile(deck.Clubs, deck.King, true),
			},
			false,
		},
		{
			"Edge case: skip card",
			[]CardInPile{
				makeCardInPile(deck.Clubs, deck.Seven, true),
				makeCardInPile(deck.Clubs, deck.Five, true),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, isValidSequence(tt.cards))
		})
	}
}

func TestIsValidRun(t *testing.T) {
	validRun := make([]CardInPile, 0, 13)
	for r := deck.King; r >= deck.Ace; r-- {
		validRun = append(validRun, makeCardInPile(deck.Spades, r, true))
	}
	assert.True(t, isValidRun(validRun), "Expected valid King > Ace run")

	invalid := append(validRun[:5], makeCardInPile(deck.Hearts, deck.Nine, true))
	assert.False(t, isValidRun(invalid), "Mixed suits should be invalid")
}

func TestCheckWinCondition(t *testing.T) {
	tests := []struct {
		name          string
		completedRuns int
		expectWon     bool
	}{
		{"not won at 0 runs", 0, false},
		{"not won at 7 runs", 7, false},
		{"won at 8 runs", 8, true},
		{"won beyond 8 runs", 9, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GameState{}

			for i := 0; i < tt.completedRuns; i++ {
				g.Completed = append(g.Completed, newSequence(deck.Spades))
			}
			g.checkWinCondition()

			if g.Won != tt.expectWon {
				t.Fatalf("expected won=%v, got %v", tt.expectWon, g.Won)
			}
		})
	}
}

// Integration type tests (Public API behavior)

func TestDealInitialGame(t *testing.T) {
	state, err := DealInitialGame()
	assert.NoError(t, err)

	totalCards := 0
	for i, pile := range state.Tableau.Piles {
		if i < 4 {
			assert.Equal(t, 6, pile.Size())
		} else {
			assert.Equal(t, 5, pile.Size())
		}
		totalCards += pile.Size()

		top, _ := pile.TopCard()
		assert.True(t, top.FaceUp)
	}
	assert.Equal(t, 54, totalCards)
	assert.Len(t, state.Stock, 50)
}

func TestDealRow(t *testing.T) {
	state, err := DealInitialGame()
	assert.NoError(t, err)

	originalStock := len(state.Stock)

	err = state.DealRow()
	assert.NoError(t, err)

	assert.Equal(t, originalStock-10, len(state.Stock))

	for _, pile := range state.Tableau.Piles {
		top := pile.Cards()[len(pile.Cards())-1]
		assert.True(t, top.FaceUp)
	}
}

func TestDealRow_FailsWhenStockIsInsufficient(t *testing.T) {
	state, _ := DealInitialGame()
	state.Stock = state.Stock[:5]

	err := state.DealRow()
	assert.ErrorIs(t, err, ErrInsufficientStock)
}

func TestMoveSequence(t *testing.T) {

	var cfde CardFaceDownError

	tests := []struct {
		name       string
		src        Pile
		dst        Pile
		startIdx   int
		expectErr  error
		validateFn func(t *testing.T, g *GameState)
	}{
		{
			name: "Valid move descending same suit",
			src: newPile(
				makeCardInPile(deck.Spades, deck.Ten, true),
				makeCardInPile(deck.Spades, deck.Nine, true),
			),
			dst:       newPile(makeCardInPile(deck.Spades, deck.Jack, true)),
			expectErr: nil,
			validateFn: func(t *testing.T, g *GameState) {
				assert.Equal(t, 0, g.Tableau.Piles[0].Size())
				assert.Equal(t, 3, g.Tableau.Piles[1].Size())
			},
		},
		{
			name: "Invalid sequence not descending",
			src: newPile(
				makeCardInPile(deck.Spades, deck.Ten, true),
				makeCardInPile(deck.Spades, deck.Eight, true),
			),
			dst:       newPile(makeCardInPile(deck.Spades, deck.Jack, true)),
			expectErr: ErrInvalidSequence,
		},
		{
			name: "Invalid sequence wrong suit",
			src: newPile(
				makeCardInPile(deck.Spades, deck.Ten, true),
				makeCardInPile(deck.Hearts, deck.Nine, true),
			),
			dst:       newPile(makeCardInPile(deck.Spades, deck.Jack, true)),
			expectErr: ErrInvalidSequence,
		},
		{
			name:      "Invalid destination wrong suit",
			src:       newPile(makeCardInPile(deck.Spades, deck.Ten, true)),
			dst:       newPile(makeCardInPile(deck.Hearts, deck.Jack, true)),
			expectErr: ErrDestinationNotAccepting,
		},
		{
			name:      "Move into empty pile allowed",
			src:       newPile(makeCardInPile(deck.Spades, deck.Ten, true)),
			dst:       newPile(),
			expectErr: nil,
			validateFn: func(t *testing.T, g *GameState) {
				assert.Equal(t, 1, g.Tableau.Piles[1].Size())
			},
		},
		{
			name: "Face down card disallowed",
			src: newPile(
				makeCardInPile(deck.Spades, deck.Jack, false),
				makeCardInPile(deck.Spades, deck.Ten, true),
			),
			dst:       newPile(),
			expectErr: cfde,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GameState{Tableau: Tableau{Piles: [10]Pile{tt.src, tt.dst}}}
			err := g.MoveSequence(0, 0, 1)

			if tt.expectErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.expectErr)
			}

			if tt.validateFn != nil {
				tt.validateFn(t, g)
			}
		})
	}
}

func TestMoveSequence_FlipsTopCard(t *testing.T) {

	src := newPile(
		makeCardInPile(deck.Spades, deck.Ace, false),
		makeCardInPile(deck.Spades, deck.Ten, true),
	)
	dst := newPile(makeCardInPile(deck.Spades, deck.Jack, true))

	g := &GameState{Tableau: Tableau{Piles: [10]Pile{src, dst}}}

	err := g.MoveSequence(0, 1, 1)
	assert.NoError(t, err)

	top, _ := g.Tableau.Piles[0].TopCard()
	assert.True(t, top.FaceUp)
}

func TestMoveSequence_CompletedRun(t *testing.T) {
	g := &GameState{Tableau: Tableau{Piles: [10]Pile{}}}

	dst := newSequenceWithIgnoreRank(deck.Spades, deck.Ace)
	g.Tableau.Piles[0].AddCards(dst)
	g.Tableau.Piles[1].AddCard(deck.Card{Suit: deck.Spades, Rank: deck.Ace}, true)

	err := g.MoveSequence(1, 0, 0)
	assert.NoError(t, err)

	assert.Equal(t, 0, g.Tableau.Piles[1].Size())
	assert.Equal(t, 0, g.Tableau.Piles[0].Size())
	assert.Len(t, g.Completed, 1)
	assert.Len(t, g.Completed[0], 13)
}

func TestDealRow_CompletesSingleRun(t *testing.T) {
	g := &GameState{Tableau: Tableau{Piles: [10]Pile{}}}

	// build k -> 2 (missing ace)
	almostRun := newSequenceWithIgnoreRank(deck.Spades, deck.Ace)
	g.Tableau.Piles[0].AddCards(almostRun)

	// stock contains exactly one Ace that will land on pile 0
	g.Stock = []deck.Card{
		{Suit: deck.Clubs, Rank: deck.King}, // pile 9
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Spades, Rank: deck.Ace}, // pile 0
	}

	err := g.DealRow()
	assert.NoError(t, err)

	assert.Len(t, g.Completed, 1)
	assert.Len(t, g.Completed[0], 13)
	assert.Equal(t, 0, g.Tableau.Piles[0].Size())
}

func TestDealRow_CompletesMultipleRuns(t *testing.T) {
	g := &GameState{Tableau: Tableau{Piles: [10]Pile{}}}

	// two piles both missing ace
	g.Tableau.Piles[0].AddCards(newSequenceWithIgnoreRank(deck.Spades, deck.Ace))
	g.Tableau.Piles[1].AddCards(newSequenceWithIgnoreRank(deck.Hearts, deck.Ace))

	// stock will deal ace to piles 0 and 1
	g.Stock = []deck.Card{
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Hearts, Rank: deck.Ace}, // pile 1
		{Suit: deck.Spades, Rank: deck.Ace}, // pile 0
	}

	err := g.DealRow()
	assert.NoError(t, err)

	assert.Len(t, g.Completed, 2)
	assert.Equal(t, 0, g.Tableau.Piles[0].Size())
	assert.Equal(t, 0, g.Tableau.Piles[1].Size())
}

func TestDealRow_DoesNotCompleteRunWithFaceDownCard(t *testing.T) {
	g := &GameState{Tableau: Tableau{Piles: [10]Pile{}}}

	// build K -> 2 but make one card face down
	seq := newSequenceWithIgnoreRank(deck.Spades, deck.Ace)
	seq[5].FaceUp = false
	g.Tableau.Piles[0].AddCards(seq)

	g.Stock = make([]deck.Card, 10)
	g.Stock[9] = deck.Card{Suit: deck.Spades, Rank: deck.Ace}

	err := g.DealRow()
	assert.NoError(t, err)

	assert.Len(t, g.Completed, 0)
	assert.Equal(t, 13, g.Tableau.Piles[0].Size())
}

func TestCheckCompletedRuns_IsIdempotent(t *testing.T) {
	g := &GameState{Tableau: Tableau{Piles: [10]Pile{}}}

	run := newSequence(deck.Spades)
	g.Tableau.Piles[0].AddCards(run)

	g.checkCompletedRuns()
	assert.Len(t, g.Completed, 1)

	// call again - should not duplicate
	g.checkCompletedRuns()
	assert.Len(t, g.Completed, 1)
}

func TestRunCompletion_ConservesTotalCards(t *testing.T) {
	g := &GameState{Tableau: Tableau{Piles: [10]Pile{}}}

	total := 0

	run := newSequence(deck.Spades)
	g.Tableau.Piles[0].AddCards(run)

	for _, pile := range g.Tableau.Piles {
		total += pile.Size()
	}

	g.checkCompletedRuns()

	after := 0
	for _, pile := range g.Tableau.Piles {
		after += pile.Size()
	}

	for _, r := range g.Completed {
		after += len(r)
	}

	assert.Equal(t, total, after)
}

func TestFullGame_DealRowAfterRunCompletion(t *testing.T) {
	g := &GameState{Tableau: Tableau{Piles: [10]Pile{}}}

	dst := newSequenceWithIgnoreRank(deck.Spades, deck.Ace)
	g.Tableau.Piles[0].AddCards(dst)
	g.Tableau.Piles[1].AddCard(deck.Card{Suit: deck.Spades, Rank: deck.Ace}, true)

	_ = g.MoveSequence(1, 0, 0)
	assert.Len(t, g.Completed, 1)

	// Add 10 dummy cards to stock for deal
	for range 10 {
		g.Stock = append(g.Stock, deck.Card{Suit: deck.Spades, Rank: deck.King})
	}

	err := g.DealRow()
	assert.NoError(t, err)
	assert.Len(t, g.Stock, 0)
}

func TestGame_WinTriggeredByMoveSequence(t *testing.T) {
	g := &GameState{Tableau: Tableau{Piles: [10]Pile{}}}

	// Preload 7 completed runs
	for range 7 {
		g.Completed = append(g.Completed, newSequence(deck.Hearts))
	}

	// build an almost complete run (missing ace)
	run := newSequenceWithIgnoreRank(deck.Hearts, deck.Ace)
	g.Tableau.Piles[0].AddCards(run)

	// add an ace to another pile
	g.Tableau.Piles[1].AddCard(deck.Card{Suit: deck.Hearts, Rank: deck.Ace}, true)

	// move the ace onto our almost completed run
	err := g.MoveSequence(1, 0, 0)
	assert.NoError(t, err)
	assert.True(t, g.Won)
}

func TestGame_WinTriggeredByDealRow(t *testing.T) {
	g := &GameState{Tableau: Tableau{Piles: [10]Pile{}}}

	// preload 7 completed runs
	for range 7 {
		g.Completed = append(g.Completed, newSequence(deck.Clubs))
	}

	// carefully construct the tableau so that dealing a row completed a run
	run := newSequenceWithIgnoreRank(deck.Diamonds, deck.Ace)
	g.Tableau.Piles[0].AddCards(run)

	// build the stock to provide the final ace
	g.Stock = []deck.Card{
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Clubs, Rank: deck.King},
		{Suit: deck.Diamonds, Rank: deck.Ace}, // pile 0
	}

	err := g.DealRow()
	assert.NoError(t, err)
	assert.True(t, g.Won)
}
