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
	Cards []CardInPile
}

// AddCard adds a card to the top of the pile
func (p *Pile) AddCard(c deck.Card, faceUp bool) {
	p.Cards = append(p.Cards, CardInPile{Card: c, FaceUp: faceUp})
}

// TopCard returns the top card without removing it
func (p *Pile) TopCard() (CardInPile, error) {
	if len(p.Cards) == 0 {
		return CardInPile{}, errors.New("pile is empty")
	}
	return p.Cards[len(p.Cards)-1], nil
}

// Cards returns a defensive copy of all cards in the pile
func (p *Pile) GetCards() []CardInPile {
	return slices.Clone(p.Cards)
}

// Size returns the number of cards in the pile
func (p *Pile) Size() int {
	return len(p.Cards)
}

// CanAccept checks if a pile can accept the given sequence
func (p *Pile) CanAccept(seq []CardInPile) bool {

	if len(seq) == 0 {
		return false
	}

	// if the pile is empty - any sequence can be placed
	if len(p.Cards) == 0 {
		return true
	}

	top := p.Cards[len(p.Cards)-1] // top card in the destination pile
	movingTop := seq[0]            // top card in the moving pile

	// must match suit and be exactly one rank lower
	return top.Card.Suit == movingTop.Card.Suit &&
		top.Card.Rank == movingTop.Card.Rank+1

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

// MoveSequence moves cards starting at index from one pile to another
func (g *GameState) MoveSequence(srcIdx, startIdx int, dsIdx int) error {

	if srcIdx < 0 || srcIdx >= TableauPiles || dsIdx < 0 || dsIdx >= TableauPiles {
		return errors.New("invalid pile index")
	}

	src := &g.Tableau.Piles[srcIdx]
	dst := &g.Tableau.Piles[dsIdx]

	if startIdx < 0 || startIdx >= src.Size() {
		return errors.New("invalid start index")
	}

	if !src.Cards[startIdx].FaceUp {
		return errors.New("cannot move face-down cards")
	}

	seq := src.GetCards()[startIdx:]
	if !isValidSequence(seq) {
		return errors.New("invalid move: sequence not ordered") // this is an error that the user will "interact" with
	}

	if !dst.CanAccept(seq) {
		return errors.New("invalid move: destination cannot accept") //same here!
	}

	// perform move
	src.Cards = src.Cards[:startIdx]
	dst.Cards = append(dst.Cards, seq...)

	if len(src.Cards) > 0 {
		top := &src.Cards[len(src.Cards)-1]
		if !top.FaceUp {
			top.FaceUp = true
		}
	}

	return nil
}

func isValidSequence(seq []CardInPile) bool {

	for i := 0; i < len(seq)-1; i++ {
		if seq[i].Card.Suit != seq[i+1].Card.Suit {
			return false
		}
		if seq[i].Card.Rank != seq[i+1].Card.Rank+1 {
			return false
		}
	}
	return true
}
