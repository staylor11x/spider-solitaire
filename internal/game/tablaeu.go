package game

import (
	"errors"

	"slices"

	"github.com/staylor11x/spider-solitaire/internal/deck"
)

// Game constants
const (
	TableauPiles     = 10
	SpiderDeckCount  = 2
	TotalSpiderCards = 104
	FirstPileCards   = 6 // first 4 piles get 6 cards
	RestPileCards    = 5 // remanining 6 piles get 5 cards
	FirstPileCount   = 4 // number of piles that get 6 cards
)

type CardInPile struct {
	Card   deck.Card
	FaceUp bool
}

type Pile struct {
	cards []CardInPile
}

// AddCard adds a card to the top of the pile
func (p *Pile) AddCard(c deck.Card, faceUp bool) {
	p.cards = append(p.cards, CardInPile{Card: c, FaceUp: faceUp})
}

// TopCard returns the top card without removing it
func (p *Pile) TopCard() (CardInPile, error) {
	if len(p.cards) == 0 {
		return CardInPile{}, errors.New("pile is empty")
	}
	return p.cards[len(p.cards)-1], nil
}

// Cards returns a defensive copy of all cards in the pile
func (p *Pile) Cards() []CardInPile {
	return slices.Clone(p.cards)
}

// Size returns the number of cards in the pile
func (p *Pile) Size() int {
	return len(p.cards)
}

// Tableau represents the 10 piles in play
type Tableau struct {
	Piles [10]Pile
}

// GameState represents the complete state of a spider game
type GameState struct {
	Tableau Tableau
	Stock   []deck.Card
}

// DealInitialGame creates a new spider layout using two decks
func DealInitialGame() (*GameState, error) {

	d := deck.NewMultiDeck(2)
	d.Shuffle()

	if d.Size() != TotalSpiderCards {
		return nil, errors.New("not enough cards for spider")
	}

	// deal tableau
	t := &Tableau{}
	for i := range TableauPiles {
		numCards := RestPileCards
		if i < FirstPileCount {
			numCards = FirstPileCards
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
	stock := d.DrawAll()

	return &GameState{
		Tableau: *t,
		Stock:   stock,
	}, nil
}

// DealRow deals one card face-up onto each tabeau pile from the stock
func (g *GameState) DealRow() error {

	if !g.CanDealRow() {
		return errors.New("not enough cards in the stock to deal a full row")
	}

	for i := range TableauPiles {
		card := g.Stock[len(g.Stock)-1] // take from the end
		g.Stock = g.Stock[:len(g.Stock)-1]
		g.Tableau.Piles[i].AddCard(card, true)
	}
	return nil
}

func (g *GameState) CanDealRow() bool {
	return len(g.Stock) >= TableauPiles
}
