package game

import (
	"errors"

	"slices"

	"github.com/staylor11x/spider-solitaire/internal/deck"
)

type CardInPile struct {
	Card   deck.Card
	FaceUp bool
}

type Pile struct {
	cards []CardInPile
}

func (p *Pile) AddCard(c deck.Card, faceUp bool) {
	p.cards = append(p.cards, CardInPile{Card: c, FaceUp: faceUp})
}

func (p *Pile) TopCard() (CardInPile, error) {
	if len(p.cards) == 0 {
		return CardInPile{}, errors.New("pile is empty")
	}
	return p.cards[len(p.cards)-1], nil
}

func (p *Pile) Cards() []CardInPile {
	return slices.Clone(p.cards) // copy for safety
}

func (p *Pile) Size() int {
	return len(p.cards)
}

// Tableau represents the 10 piles in play
type Tableau struct {
	Piles [10]Pile
}

// DealInitialLayout deals the first 54 card into the tableau
func DealInitialLayout(d *deck.Deck) (*Tableau, error) {

	if d.Size() < 52 { // this will change as we update to two decks
		return nil, errors.New("not enough cards to deal tableau")
	}

	t := &Tableau{}

	// first 4 piles get 6 cards, last 6 piles get 5 cards
	for i := range 10 {
		numCards := 5
		if i < 4 {
			numCards = 6
		}

		for j := range numCards {
			card, err := d.Draw()
			if err != nil {
				return nil, err
			}
			// draw last card face up
			faceUp := j == numCards-1
			t.Piles[i].AddCard(card, faceUp)
		}
	}

	return t, nil
}

type GameState struct {
	Tableau Tableau
	Stock   []deck.Card
}

// DealinitialGame creates a new spider layout using two decks
func DealinitialGame() (*GameState, error) {

	d := deck.NewMultiDeck(2)
	d.Shuffle()

	if d.Size() < 104 {
		return nil, errors.New("not enough cards for spider")
	}

	// deal tableau
	t := &Tableau{}
	for i := range 10 {
		numCards := 5
		if i < 4 {
			numCards = 6
		}

		for j := range numCards {
			card, err := d.Draw()
			if err != nil {
				return nil, err
			}
			faceUp := j == numCards-1
			t.Piles[i].AddCard(card, faceUp)
		}
	}

	// remaining cards from the stock
	stock := make([]deck.Card, d.Size())
	copy(stock, d.Cards()) // need a safe accessor in the deck

	return &GameState{
		Tableau: *t,
		Stock:   stock,
	}, nil
}
