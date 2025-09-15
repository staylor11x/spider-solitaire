package game

import (
	"slices"

	"github.com/staylor11x/spider-solitaire/internal/deck"
)

// Game constants
const (
	TableauPiles     = 10
	SpiderDeckCount  = 2
	TotalSpiderCards = 104
	FirstPileCards   = 6  // first 4 piles get 6 cards
	RestPileCards    = 5  // remaining 6 piles get 5 cards
	FirstPileCount   = 4  // number of piles that get 6 cards
	RunLength        = 13 // King to Ace
)

// GameState represents the complete state of a spider game
type GameState struct {
	Tableau   Tableau
	Stock     []deck.Card
	Completed [][]CardInPile
}

// DealInitialGame creates a new spider layout using two decks
func DealInitialGame() (*GameState, error) {

	d := deck.NewMultiDeck(2)
	d.Shuffle()

	if d.Size() != TotalSpiderCards {
		return nil, ErrNotEnoughCards
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

// DealRow deals one card face-up onto each tableau pile from the stock
func (g *GameState) DealRow() error {

	if !g.canDealRow() {
		return ErrInsufficientStock
	}

	for i := range TableauPiles {
		card := g.Stock[len(g.Stock)-1] // take from the end
		g.Stock = g.Stock[:len(g.Stock)-1]
		g.Tableau.Piles[i].AddCard(card, true)
	}
	return nil
}

func (g *GameState) canDealRow() bool {
	return len(g.Stock) >= TableauPiles
}

func (g *GameState) MoveSequence(srcIdx, startIdx, dstIdx int) error {

	if err := g.validateMoveIndices(srcIdx, startIdx, dstIdx); err != nil {
		return err
	}

	src := &g.Tableau.Piles[srcIdx]
	dst := &g.Tableau.Piles[dstIdx]

	sequence, err := g.validateMoveSequence(src, startIdx)
	if err != nil {
		return err
	}

	if !dst.CanAccept(sequence) {
		return ErrDestinationNotAccepting
	}

	// perform atomic move
	return g.executeMove(src, dst, startIdx, sequence)
}

func (g *GameState) validateMoveIndices(srcIdx, startIdx, dstIdx int) error {

	if srcIdx < 0 || srcIdx >= TableauPiles {
		return ErrInvalidSourceIndex
	}

	if dstIdx < 0 || dstIdx >= TableauPiles {
		return ErrInvalidDestinationIndex
	}

	if srcIdx == dstIdx {
		return ErrSamePileMove
	}

	src := &g.Tableau.Piles[srcIdx]
	if startIdx < 0 || startIdx >= src.Size() {
		return ErrInvalidStartIndex
	}

	return nil
}

func (g *GameState) validateMoveSequence(src *Pile, startIdx int) ([]CardInPile, error) {

	allCards := src.Cards() // returns a copy of the cards
	sequence := allCards[startIdx:]

	if len(sequence) == 0 {
		return nil, ErrNoCardsToMove
	}

	// check all cards are face up
	for i, card := range sequence {
		if !card.FaceUp {
			return nil, CardFaceDownError{Index: startIdx + i}
		}
	}

	// validate sequence is properly ordered
	if !isValidSequence(sequence) {
		return nil, ErrInvalidSequence
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
		return ErrRemoveCardsWithContext(err)
	}

	// Paranoid check: verify removed cards match expected sequence
	if !sequenceEqual(removedCards, sequence) {
		// critical error, restore the cards and fail
		src.AddCards(removedCards)
		return ErrSequenceMismatch
	}

	// add cards to destination
	dst.AddCards(removedCards)

	// flip top card of source if needed
	if err := src.FlipTopCardIfFaceDown(); err != nil {
		return ErrFlipWithContext(err)
	}

	// check for completed runs
	g.checkCompletedRuns()

	return nil
}

func sequenceEqual(a, b []CardInPile) bool {
	return slices.EqualFunc(a, b, func(x, y CardInPile) bool {
		return x.FaceUp == y.FaceUp && x.Card == y.Card
	})
}

// checkCompletedRuns scans each pile for a complete run and removed it if found, storing it in g.Completed
func (g *GameState) checkCompletedRuns() {

	for i := range g.Tableau.Piles {
		pile := &g.Tableau.Piles[i]
		if pile.Size() < RunLength {
			continue
		}

		// look at the last 13 cards
		last := pile.Cards()[pile.Size()-RunLength:]
		if isValidRun(last) {
			removed, _ := pile.RemoveCardsFrom(pile.Size() - RunLength)

			// add to completed
			g.Completed = append(g.Completed, removed)

			// flip to card if needed - why is there here?
			_ = pile.FlipTopCardIfFaceDown()
		}
	}
}

// isValidRun checks for a perfect King->Ace descending run in one suit
func isValidRun(cards []CardInPile) bool {

	if len(cards) != RunLength {
		return false
	}

	if !cards[0].FaceUp {
		return false
	}

	// must start with a king
	if cards[0].Card.Rank != deck.King {
		return false
	}

	// check each step -- there is probs some hyper efficient way to do this
	for i := 0; i < len(cards)-1; i++ {
		if !cards[i].FaceUp || !cards[i+1].FaceUp {
			return false
		}
		if cards[i].Card.Suit != cards[i+1].Card.Suit {
			return false
		}
		if cards[i].Card.Rank != cards[i+1].Card.Rank+1 {
			return false
		}
	}
	// must end on an ace
	return cards[len(cards)-1].Card.Rank == deck.Ace
}
