package deck

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStandardDeck(t *testing.T) {

	deck := NewStandardDeck()
	assert.Equal(t, 52, len(deck.cards), "deck should contain 52 cards")

	// ensure no duplicates
	seen := make(map[string]bool)
	for _, c := range deck.cards {
		if seen[c.String()] {
			t.Fatalf("duplicate card found: %s", c)
		}
		seen[c.String()] = true
	}
}

func TestShuffle(t *testing.T) {
	deck1 := NewStandardDeck()
	deck2 := NewStandardDeck()
	deck2.Shuffle()

	// check that the order is not the same
	same := true
	for i := range deck1.cards {
		if deck1.cards[i] != deck2.cards[i] {
			same = false
			break
		}
	}
	assert.False(t, same, "shuffled deck should not be in the same order")
}

func TestDraw(t *testing.T) {
	deck := NewStandardDeck()
	firstCard, err := deck.Draw()
	assert.NoError(t, err)
	assert.Equal(t, 51, len(deck.cards))
	assert.Equal(t, "Ace of Spades", firstCard.String())

	// draw all the cards
	for i := 0; i < 51; i++ {
		_, err := deck.Draw()
		assert.NoError(t, err)
	}

	// deck should be empty
	_, err = deck.Draw()
	assert.Error(t, err, "drawing from empty deck should error")
}
