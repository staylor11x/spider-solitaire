package deck

import (
	"errors"
	"math/rand"
	"time"
)

type SuitCount int
type DeckCount int

const (
	OneSuit   SuitCount = 1
	TwoSuits  SuitCount = 2
	FourSuits SuitCount = 4

	OneSuitDeckCount  DeckCount = 8
	TwoSuitDeckCount  DeckCount = 4
	FourSuitDeckCount DeckCount = 2
)

// NewSpiderDeck creates a 104-card deck for Spider Solitaire with the specified number of suits
func NewSpiderDeck(suitCount SuitCount) *Deck {
	cards := make([]Card, 0, 104)
	switch suitCount {
	case OneSuit:
		// 8 decks of one suit
		for range OneSuitDeckCount {
			for r := Ace; r <= King; r++ {
				cards = append(cards, Card{Suit: Spades, Rank: r})
			}
		}
	case TwoSuits:
		suits := []Suit{Spades, Hearts}
		for range TwoSuitDeckCount {
			for _, s := range suits {
				for r := Ace; r <= King; r++ {
					cards = append(cards, Card{Suit: s, Rank: r})
				}
			}
		}
	case FourSuits:
		for range FourSuitDeckCount {
			for s := Spades; s <= Clubs; s++ {
				for r := Ace; r <= King; r++ {
					cards = append(cards, Card{Suit: s, Rank: r})
				}
			}
		}
	}
	return &Deck{cards: cards}
}

// Deck represents a standard deck (or n decks) of cards
type Deck struct {
	cards []Card
}

// NewStandardDeck creates a standard 52-card deck.
func NewStandardDeck() Deck {

	// pre allocate the slice
	cards := make([]Card, 0, 52)

	for _, s := range []Suit{Spades, Hearts, Diamonds, Clubs} {
		for r := Ace; r <= King; r++ {
			cards = append(cards, Card{Suit: s, Rank: r})
		}
	}
	return Deck{cards}
}

// NewMultiDeck creates n standard decks, combined into one.
func NewMultiDeck(n int) Deck {

	if n <= 0 {
		return Deck{cards: []Card{}}
	}

	cards := make([]Card, 0, 52*n)

	for range n {
		for _, s := range []Suit{Spades, Hearts, Diamonds, Clubs} {
			for r := Ace; r <= King; r++ {
				cards = append(cards, Card{Suit: s, Rank: r})
			}
		}
	}
	return Deck{cards}
}

// Shuffle randomizes the order of the deck
func (d *Deck) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

// Draw removes and returns the top card from the deck
func (d *Deck) Draw() (Card, error) {

	if len(d.cards) == 0 {
		return Card{}, errors.New("deck is empty")
	}

	lastIndex := len(d.cards) - 1
	card := d.cards[lastIndex]
	d.cards = d.cards[:lastIndex]
	return card, nil
}

func (d *Deck) Size() int {
	return len(d.cards)
}

// Cards returns a copy of the cards in the deck
func (d *Deck) Cards() []Card {
	c := make([]Card, len(d.cards))
	copy(c, d.cards)
	return c
}

// DrawAll removes and returns all remaining cards from the deck
func (d *Deck) DrawAll() []Card {
	if len(d.cards) == 0 {
		return nil
	}

	remaining := make([]Card, len(d.cards))
	copy(remaining, d.cards)
	d.cards = d.cards[:0] // clear the deck
	return remaining
}
