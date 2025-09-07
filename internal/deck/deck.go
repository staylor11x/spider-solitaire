package deck

import (
	"errors"
	"math/rand"
	"time"
)

// Deck represents a standard deck (or n decks) of Cards
type Deck struct {
	Cards []Card
}

// NewStandardDeck creates a standard 52-card deck.
func NewStandardDeck() Deck {

	// pre allocate the slice
	Cards := make([]Card, 0, 52)

	for _, s := range []Suit{Spades, Hearts, Diamonds, Clubs} {
		for r := Ace; r <= King; r++ {
			Cards = append(Cards, Card{Suit: s, Rank: r})
		}
	}
	return Deck{Cards}
}

// NewMultiDeck creates n standard decks, combined into one.
func NewMultiDeck(n int) Deck {

	if n <= 0 {
		return Deck{Cards: []Card{}}
	}

	Cards := make([]Card, 0, 52*n)

	for range n {
		for _, s := range []Suit{Spades, Hearts, Diamonds, Clubs} {
			for r := Ace; r <= King; r++ {
				Cards = append(Cards, Card{Suit: s, Rank: r})
			}
		}
	}
	return Deck{Cards}
}

// Shuffle randomises the order of the deck
func (d *Deck) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

// Draw removes and returns the top card from the deck
func (d *Deck) Draw() (Card, error) {

	if len(d.Cards) == 0 {
		return Card{}, errors.New("deck is empty")
	}

	lastIndex := len(d.Cards) - 1
	card := d.Cards[lastIndex]
	d.Cards = d.Cards[:lastIndex]
	return card, nil
}

func (d *Deck) Size() int {
	return len(d.Cards)
}

// Cards returns a copy of the Cards in the deck
func (d *Deck) GetCards() []Card {
	c := make([]Card, len(d.Cards))
	copy(c, d.Cards)
	return c
}

// DrawAll removes and returns all remaining Cards from the deck
func (d *Deck) DrawAll() []Card {
	if len(d.Cards) == 0 {
		return nil
	}

	remaining := make([]Card, len(d.Cards))
	copy(remaining, d.Cards)
	d.Cards = d.Cards[:0] // clear the deck
	return remaining
}
