package game_test

import (
	"testing"

	"github.com/staylor11x/spider-solitaire/internal/deck"
	"github.com/staylor11x/spider-solitaire/internal/game"
	"github.com/staylor11x/spider-solitaire/internal/testtools"
	"github.com/stretchr/testify/assert"
)

func TestDealInitialGame(t *testing.T) {
	state, err := game.DealInitialGame()
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

	// stock should have 50 cards left
	assert.Equal(t, 50, len(state.Stock))
}

func TestDealRow(t *testing.T) {
	state, err := game.DealInitialGame()
	assert.NoError(t, err)

	originalStock := len(state.Stock)

	err = state.DealRow()
	assert.NoError(t, err)

	// stock should be reduced by 10
	assert.Equal(t, originalStock-10, len(state.Stock))

	// each pile should be incremented by 1
	for _, pile := range state.Tableau.Piles {
		assert.NotEmpty(t, pile.Cards())
		top := pile.Cards()[len(pile.Cards())-1]
		assert.True(t, top.FaceUp, "newly dealt cards should be face up")
	}
}

func TestDealRowFailedWIthInsufficientStock(t *testing.T) {
	state, _ := game.DealInitialGame()

	// drain the stock
	state.Stock = state.Stock[:5]

	err := state.DealRow()
	assert.Error(t, err, "should fail when fewer than 10 cards remain")
}

func TestMoveSequence_ValidMove(t *testing.T) {

	src := testtools.NewPile(
		testtools.MakeCardInPile(deck.Spades, deck.Ten, true),
		testtools.MakeCardInPile(deck.Spades, deck.Nine, true),
	)

	dst := testtools.NewPile(
		testtools.MakeCardInPile(deck.Spades, deck.Jack, true),
	)

	g := game.GameState{Tableau: game.Tableau{Piles: [10]game.Pile{src, dst}}}

	err := g.MoveSequence(0, 0, 1)
	assert.NoError(t, err)

	assert.Equal(t, 0, g.Tableau.Piles[0].Size(), "source should be empty")
	assert.Equal(t, 3, g.Tableau.Piles[1].Size(), "destination should have three cards")

	top, _ := g.Tableau.Piles[1].TopCard()
	assert.Equal(t, deck.Nine, top.Card.Rank, "top card should be 9 of hearts")
}

func TestMoveSequence_InvalidSequence_NotDescending(t *testing.T) {

	// src pile: 10S, 8S (gap invalid)
	src := testtools.NewPile(
		testtools.MakeCardInPile(deck.Spades, deck.Ten, true),
		testtools.MakeCardInPile(deck.Spades, deck.Eight, true),
	)

	dst := testtools.NewPile(
		testtools.MakeCardInPile(deck.Spades, deck.Jack, true),
	)

	g := game.GameState{Tableau: game.Tableau{Piles: [10]game.Pile{src, dst}}}

	err := g.MoveSequence(0, 0, 1)
	assert.ErrorIs(t, err, game.ErrInvalidSequence)
}

func TestMoveSequence_InvalidSequence_WrongSuit(t *testing.T) {

	// src pile 10S, 9H
	src := testtools.NewPile(
		testtools.MakeCardInPile(deck.Spades, deck.Ten, true),
		testtools.MakeCardInPile(deck.Hearts, deck.Nine, true),
	)

	dst := testtools.NewPile(
		testtools.MakeCardInPile(deck.Spades, deck.Jack, true),
	)

	g := game.GameState{Tableau: game.Tableau{Piles: [10]game.Pile{src, dst}}}

	err := g.MoveSequence(0, 0, 1)
	assert.ErrorIs(t, err, game.ErrInvalidSequence)
}

func TestMoveSequence_InvalidDestination(t *testing.T) {

	// src pile 10S
	src := testtools.NewPile(
		testtools.MakeCardInPile(deck.Spades, deck.Ten, true),
	)

	// dst pile: JH (wrong suit)
	dst := testtools.NewPile(
		testtools.MakeCardInPile(deck.Hearts, deck.Jack, true),
	)

	g := game.GameState{Tableau: game.Tableau{Piles: [10]game.Pile{src, dst}}}

	err := g.MoveSequence(0, 0, 1)
	assert.ErrorIs(t, err, game.ErrDestinationNotAccepting)
}

func TestMoveSequence_MoveIntoEmptyPile(t *testing.T) {

	//src pile: 10S
	src := testtools.NewPile(
		testtools.MakeCardInPile(deck.Spades, deck.Ten, true),
	)

	//dst pile: empty
	dst := testtools.NewPile()

	g := game.GameState{Tableau: game.Tableau{Piles: [10]game.Pile{src, dst}}}

	err := g.MoveSequence(0, 0, 1)
	assert.NoError(t, err)

	assert.Equal(t, 0, g.Tableau.Piles[0].Size())
	assert.Equal(t, 1, g.Tableau.Piles[1].Size())
}

func TestMoveSequence_FaceDownCardDisallowed(t *testing.T) {
	src := testtools.NewPile(
		testtools.MakeCardInPile(deck.Spades, deck.Jack, false), // face down
		testtools.MakeCardInPile(deck.Spades, deck.Ten, true),   // face up
	)

	dst := testtools.NewPile()

	g := &game.GameState{Tableau: game.Tableau{Piles: [10]game.Pile{src, dst}}}

	err := g.MoveSequence(0, 0, 1)
	assert.ErrorIs(t, err, game.CardFaceDownError{})
}

func TestMoveSequence_FlipsTopCard(t *testing.T) {

	src := testtools.NewPile(
		testtools.MakeCardInPile(deck.Spades, deck.Ace, false),
		testtools.MakeCardInPile(deck.Spades, deck.Ten, true),
	)

	dst := testtools.NewPile(
		testtools.MakeCardInPile(deck.Spades, deck.Jack, true),
	)

	g := &game.GameState{Tableau: game.Tableau{Piles: [10]game.Pile{src, dst}}}

	err := g.MoveSequence(0, 1, 1) // move the ten onto the jack
	assert.NoError(t, err)

	// check that top card has been flipped
	top, _ := g.Tableau.Piles[0].TopCard()
	assert.True(t, top.FaceUp)
}

func TestMoveSequence_CompletedRun(t *testing.T) {

	g := &game.GameState{Tableau: game.Tableau{Piles: [10]game.Pile{}}}

	// make a pile all cards apart from an ace
	dst := testtools.NewSequenceWithIgnoreRank(deck.Spades, deck.Ace)
	g.Tableau.Piles[0].AddCards(dst)

	// add the ace to another pile
	g.Tableau.Piles[1].AddCard(deck.Card{deck.Spades, deck.Ace}, true)

	// move the ace to the almost complete pile
	err := g.MoveSequence(1, 0, 0)
	assert.NoError(t, err, "unexpected error moving cards: %v", err)

	// pile 1 should be empty
	assert.Equal(t, 0, g.Tableau.Piles[1].Size(), "pile 1 is not empty")

	// run should be recorded as completed
	assert.Equal(t, 1, len(g.Completed))

	// there should be 13 cards in the completed run
	assert.Equal(t, 13, len(g.Completed[0]))

}
