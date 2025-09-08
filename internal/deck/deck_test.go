package deck

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStandardDeck(t *testing.T) {

	d := NewStandardDeck()
	assert.Equal(t, 52, d.Size(), "deck should contain 52 cards")

	// ensure no duplicates
	seen := make(map[string]bool)
	for _, c := range d.cards {
		if seen[c.String()] {
			t.Fatalf("duplicate card found: %s", c)
		}
		seen[c.String()] = true
	}
}

func TestNewMultiDeck(t *testing.T) {

	t.Run("two decks", func(t *testing.T) {
		d := NewMultiDeck(2)
		assert.Equal(t, 104, d.Size(), "two decks should have 104 cards")

		// count occurences of each card
		counts := make(map[string]int)
		for _, c := range d.Cards() {
			counts[c.String()]++
		}

		// Each unique card should appear two times
		assert.Equal(t, 52, len(counts), "should have 52 unique cards")
		for card, count := range counts {
			assert.Equal(t, 2, count, "card %s should appear exactly twice", card)
		}
	})

	t.Run("single deck", func(t *testing.T) {
		d := NewMultiDeck(1)
		assert.Equal(t, 52, d.Size(), "one deck should have 52 cards")
	})

	t.Run("zero deck", func(t *testing.T) {
		d := NewMultiDeck(0)
		assert.Equal(t, 0, d.Size(), "zero deck should produce empty deck")
	})
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
	assert.Equal(t, "King of Clubs", firstCard.String())

	// draw all the cards
	for i := 0; i < 51; i++ {
		_, err := deck.Draw()
		assert.NoError(t, err)
	}

	// deck should be empty
	_, err = deck.Draw()
	assert.Error(t, err, "drawing from empty deck should error")
}
