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

// CanAccept checks if a pile can accept the given sequence
func (p *Pile) CanAccept(seq []CardInPile) bool {

	if len(seq) == 0 {
		return false
	}

	// if the pile is empty - any sequence can be placed
	if len(p.cards) == 0 {
		return true
	}

	top := p.cards[len(p.cards)-1] // top card in the destination pile
	movingTop := seq[0]            // top card in the moving pile

	// must match suit and be exactly one rank lower
	// return top.Card.Suit == movingTop.Card.Suit &&
	return top.Card.Rank == movingTop.Card.Rank+1

}

func (p *Pile) RemoveCardsFrom(startIdx int) ([]CardInPile, error) {
	if startIdx < 0 || startIdx >= len((p.cards)) {
		return nil, ErrInvalidStartIndex
	}

	removed := make([]CardInPile, len(p.cards)-startIdx)
	copy(removed, p.cards[startIdx:])

	p.cards = p.cards[:startIdx]
	return removed, nil
}

func (p *Pile) AddCards(cards []CardInPile) {
	p.cards = append(p.cards, cards...)
}

func (p *Pile) FlipTopCardIfFaceDown() error {
	if len(p.cards) == 0 {
		return nil // no card to flip - this is ok
	}

	topIdx := len(p.cards) - 1
	if !p.cards[topIdx].FaceUp {
		p.cards[topIdx].FaceUp = true
	}

	return nil
}
