package game

import (
	"errors"
	"fmt"

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

func (g *GameState) MoveSequence(srcIdx, startIdx, dstIdx int) error {

	if err := g.validateMoveIndicies(srcIdx, startIdx, dstIdx); err != nil {
		return err
	}

	src := &g.Tableau.Piles[srcIdx]
	dst := &g.Tableau.Piles[dstIdx]

	sequence, err := g.validateMoveSequence(src, startIdx)
	if err != nil {
		return err
	}

	if !dst.CanAccept(sequence) {
		return errors.New("invalid move: destination cannot accept")
	}

	// perform atomic move
	return g.executeMove(src, dst, startIdx, sequence)
}

func (g *GameState) validateMoveIndicies(srcIdx, startIdx, dstIdx int) error {

	if srcIdx < 0 || srcIdx >= TableauPiles {
		return errors.New("invalid source pile index")
	}

	if dstIdx < 0 || dstIdx >= TableauPiles {
		return errors.New("invalid destination pile index")
	}

	if srcIdx == dstIdx {
		return errors.New("cannot move cards to the same pile")
	}

	src := &g.Tableau.Piles[srcIdx]
	if startIdx < 0 || startIdx >= src.Size() {
		return errors.New("invalid start index")
	}

	return nil
}

func (g *GameState) validateMoveSequence(src *Pile, startIdx int) ([]CardInPile, error) {

	allCards := src.Cards() // returns a copy of the cards
	sequence := allCards[startIdx:]

	if len(sequence) == 0 {
		return nil, errors.New("no cards to move")
	}

	// check all cards are face up
	for i, card := range sequence {
		if !card.FaceUp {
			return nil, fmt.Errorf("card at position %d is face down", startIdx+i)
		}
	}

	// validate sequence is properly ordered
	if !isValidSequence(sequence) {
		return nil, errors.New("invalid move: sequence not ordered")
	}

	return sequence, nil
}

func isValidSequence(cards []CardInPile) bool {
	if len(cards) <= 1 {
		return true
	}

	for i := 0; i < len(cards)-1; i++ {
		current := cards[i].Card
		next := cards[i+1].Card

		// must be same suit
		if current.Suit != next.Suit {
			return false
		}

		// must be descending rank
		if current.Rank != next.Rank+1 {
			return false
		}
	}

	return true
}

func (g *GameState) executeMove(src, dst *Pile, startIdx int, sequence []CardInPile) error {

	removedCards, err := src.RemoveCardsFrom(startIdx)
	if err != nil {
		return fmt.Errorf("failed to remove cards: %w", err)
	}

	// Paranoid check: verify removed cards match expected sequence
	if !sequenceEqual(removedCards, sequence) {
		// critical error, restore the cards and fail
		src.AddCards(removedCards)
		return errors.New("internal error: removed cards don't match expected sequence")
	}

	// add cards to destination
	dst.AddCards(removedCards)

	// flip top card of ssource if needed
	if err := src.FlipTopCardIfFaceDown(); err != nil {
		return fmt.Errorf("failed to flip source card: %w", err)
	}

	return nil
}

func sequenceEqual(a, b []CardInPile) bool {
	return slices.EqualFunc(a, b, func(x, y CardInPile) bool {
		return x.FaceUp == y.FaceUp && x.Card == y.Card
	})
}
