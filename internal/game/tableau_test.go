package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDealInitialGame(t *testing.T) {
	state, err := DealinitialGame()
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
