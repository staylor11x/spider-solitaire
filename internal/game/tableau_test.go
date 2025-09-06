package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	// stock should have 50 cards left
	assert.Equal(t, 50, len(state.Stock))
}

func TestDealRow(t *testing.T) {
	state, err := DealInitialGame()
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
	state, _ := DealInitialGame()

	// drain the stock
	state.Stock = state.Stock[:5]

	err := state.DealRow()
	assert.Error(t, err, "should fail when fewer than 10 cards remain")
}
