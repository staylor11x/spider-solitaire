package game

import (
	"testing"

	"github.com/staylor11x/spider-solitaire/internal/deck"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUndo_EmptyHistory(t *testing.T) {
	g, err := DealInitialGame(deck.OneSuit)
	assert.NoError(t, err)

	// should fail when no history exists
	err = g.Undo()
	assert.ErrorIs(t, err, ErrNoHistory)
}

func TestSnapshot_IsolatesState(t *testing.T) {
	g, err := DealInitialGame(deck.OneSuit)
	require.NoError(t, err)

	snap := g.snapshot()

	// modify original state
	g.Won = true
	g.Lost = true // bit wild
	originalStockSize := len(g.Stock)
	if originalStockSize > 0 {
		g.Stock = g.Stock[:originalStockSize-1] // remove a single card
	}

	// modify original state
	assert.False(t, snap.Won, "snapshot should not reflect Won change")
	assert.False(t, snap.Lost, "snapshot should not reflect Lost change")
	assert.Equal(t, originalStockSize, len(snap.Stock), "snapshot stock should be unchanged")
}

func TestUndo_AfterDealRow(t *testing.T) {
	g, err := DealInitialGame(deck.OneSuit)
	require.NoError(t, err)

	// record initial state
	initialStockSize := len(g.Stock)
	initialPileSizes := make([]int, TableauPiles)
	for i := range TableauPiles {
		initialPileSizes[i] = g.Tableau.Piles[i].Size()
	}

	// deal a row
	err = g.DealRow()
	require.NoError(t, err)

	// verify state changed
	assert.Equal(t, initialStockSize-TableauPiles, len(g.Stock), "stock should decrease by 10")
	for i := range TableauPiles {
		assert.Equal(t, initialPileSizes[i]+1, g.Tableau.Piles[i].Size(), "pile %d should have one more card", i)
	}

	// undo the deal
	err = g.Undo()
	require.NoError(t, err)

	// verify state restored
	assert.Equal(t, initialStockSize, len(g.Stock), "stock should be restored")
	for i := range TableauPiles {
		assert.Equal(t, initialPileSizes[i], g.Tableau.Piles[i].Size(), "pile %d size should be restored", i)
	}
}

func TestUndo_AfterValidMove(t *testing.T) {
	// create a custom game state with a known valid move
	g := &GameState{
		Tableau: Tableau{},
		Stock:   []deck.Card{},
	}

	// pile 0: 7 spades
	g.Tableau.Piles[0].AddCard(deck.Card{Suit: deck.Spades, Rank: deck.Seven}, true)

	// Pile 1: 8 spades
	g.Tableau.Piles[1].AddCard(deck.Card{Suit: deck.Spades, Rank: deck.Eight}, true)

	// record initial state
	pile0Initial := g.Tableau.Piles[0].Size()
	pile1Initial := g.Tableau.Piles[1].Size()

	// move 7 from pile 0 into 8 @ pile 1
	err := g.MoveSequence(0, 0, 1)
	require.NoError(t, err)

	// verify move happened
	assert.Equal(t, 0, g.Tableau.Piles[0].Size(), "source pile should be empty")
	assert.Equal(t, 2, g.Tableau.Piles[1].Size(), "destination pile should have 2 cards")

	// undo the move
	err = g.Undo()
	require.NoError(t, err)

	assert.Equal(t, pile0Initial, g.Tableau.Piles[0].Size(), "source pile size restored")
	assert.Equal(t, pile1Initial, g.Tableau.Piles[1].Size(), "destination pile size restored")

	top0, _ := g.Tableau.Piles[0].TopCard()
	assert.Equal(t, deck.Rank(deck.Seven), top0.Card.Rank, "source pile should have 7")

	top1, _ := g.Tableau.Piles[1].TopCard()
	assert.Equal(t, deck.Rank(deck.Eight), top1.Card.Rank, "source pile should have 8")
}

func TestUndo_MultipleUndos(t *testing.T) {
	g, err := DealInitialGame(deck.OneSuit)
	require.NoError(t, err)

	initialStock := len(g.Stock)

	// Deal 3 rows
	require.NoError(t, g.DealRow())
	require.NoError(t, g.DealRow())
	require.NoError(t, g.DealRow())

	// Stock should have decreased by 30
	assert.Equal(t, initialStock-30, len(g.Stock))

	// undo all three
	require.NoError(t, g.Undo())
	require.NoError(t, g.Undo())
	require.NoError(t, g.Undo())

	// should be back to the initial state
	assert.Equal(t, initialStock, len(g.Stock))

	// 4th undo should fail
	err = g.Undo()
	assert.ErrorIs(t, err, ErrNoHistory)
}

func TestUndo_HistoryLimit(t *testing.T) {
	g, err := DealInitialGame(deck.OneSuit)
	require.NoError(t, err)

	dealsPerformed := 0
	for range maxHistorySize + 5 {
		if len(g.Stock) >= TableauPiles {
			err := g.DealRow()
			if err == nil {
				dealsPerformed++
			}
		}
	}

	// history should be capped at maxHistorySize
	assert.LessOrEqual(t, len(g.history), maxHistorySize, "history size should not exceed max size")

	// should be able to undo up to maxHistorySize times
	successfulUndos := 0
	for range maxHistorySize + 5 {
		err := g.Undo()
		if err == nil {
			successfulUndos++
		} else {
			assert.ErrorIs(t, err, ErrNoHistory)
			break
		}
	}

	assert.LessOrEqual(t, successfulUndos, maxHistorySize, "should not undo more than the max history size")
}

func TestUndo_PreserveWonLostFlags(t *testing.T) {
	g, err := DealInitialGame(deck.OneSuit)
	require.NoError(t, err)

	// manually set the run flag
	g.Won = true

	// deal a row (this saves snapshot with Won=true)
	err = g.DealRow()
	require.NoError(t, err)

	// Won flag should still be true
	assert.True(t, g.Won)

	// Undo
	err = g.Undo()
	require.NoError(t, err)

	// won flag should be restored to true (from snapshot)
	assert.True(t, g.Won, "won flag should be preserved in undo")
}

func TestUndo_DealWithRunCompletion_SingleUndo(t *testing.T) {
	// create a game state where dealing will complete a run
	g := &GameState{
		Tableau: Tableau{},
		Stock:   make([]deck.Card, 10),
	}

	// build an almost-complete run in pile 0 (missing ace)
	g.Tableau.Piles[0].AddCards(newSequenceWithIgnoreRank(deck.Spades, deck.Ace))

	// put an ace in the stock (will complete the run)
	g.Stock[9] = deck.Card{Suit: deck.Spades, Rank: deck.Ace}
	for i := 0; i < 9; i++ {
		g.Stock[i] = deck.Card{Suit: deck.Hearts, Rank: deck.Two}
	}

	initialPileSize := g.Tableau.Piles[0].Size()
	initialStockSize := len(g.Stock)

	// deal row will complete run and remove it
	err := g.DealRow()
	require.NoError(t, err)

	// verify that the run was completed
	assert.Equal(t, 1, len(g.Completed), "run should be completed")
	assert.Equal(t, 0, len(g.Stock), "stock should be empty")

	// ONE undo should revert BOTH the deal and the run removal
	err = g.Undo()
	require.NoError(t, err)

	assert.Equal(t, 0, len(g.Completed), "completed runs should be empty after undo")
	assert.Equal(t, initialPileSize, g.Tableau.Piles[0].Size(), "pile 0 should have original size")
	assert.Equal(t, initialStockSize, len(g.Stock), "stock should be restored")

	// second undo should fail, no history
	err = g.Undo()
	assert.ErrorIs(t, err, ErrNoHistory, "second undo should fail")
}

func TestUndo_DeepCopyVerification(t *testing.T) {
	g, err := DealInitialGame(deck.OneSuit)
	require.NoError(t, err)

	// save initial top card of pile 0
	initialTop, err := g.Tableau.Piles[0].TopCard()
	require.NoError(t, err)

	// deal a row (creates snapshot)
	err = g.DealRow()
	require.NoError(t, err)

	// undo
	err = g.Undo()
	require.NoError(t, err)

	// top card should match the initial card, not the modified one
	restoredTop, _ := g.Tableau.Piles[0].TopCard()
	assert.Equal(t, initialTop.Card, restoredTop.Card, "deep copy should preserve card identity")
}
