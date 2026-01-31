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
	TotalRunsToWin   = 8
	maxHistorySize   = 25
)

// GameState represents the complete state of a spider game
type GameState struct {
	Tableau   Tableau
	Stock     []deck.Card
	Completed [][]CardInPile
	Won       bool
	Lost      bool
	history   []GameState // would be good to explain the context of this recursive style structure.
}

// DealInitialGame creates a new spider layout using two decks
func DealInitialGame(suitCount deck.SuitCount) (*GameState, error) {

	d := deck.NewSpiderDeck(suitCount)
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
	g.pushHistory()

	for i := range TableauPiles {
		card := g.Stock[len(g.Stock)-1] // take from the end
		g.Stock = g.Stock[:len(g.Stock)-1]
		g.Tableau.Piles[i].AddCard(card, true)
	}
	if err := g.checkCompletedRuns(); err != nil {
		return err
	}

	// only check when there is no more stock
	if len(g.Stock) == 0 {
		g.checkLossCondition()
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
	g.pushHistory()

	// perform atomic move
	err = g.executeMove(src, dst, startIdx, sequence)
	if err != nil {
		return err
	}

	// only check when there is no more stock
	if len(g.Stock) == 0 {
		g.checkLossCondition()
	}
	return nil
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
	if err := g.checkCompletedRuns(); err != nil {
		return err
	}
	return nil
}

func sequenceEqual(a, b []CardInPile) bool {
	return slices.EqualFunc(a, b, func(x, y CardInPile) bool {
		return x.FaceUp == y.FaceUp && x.Card == y.Card
	})
}

// checkCompletedRuns scans each pile for a complete run and removed it if found, storing it in g.Completed
func (g *GameState) checkCompletedRuns() error {

	for i := range g.Tableau.Piles {
		pile := &g.Tableau.Piles[i]
		if pile.Size() < RunLength {
			continue
		}

		// look at the last 13 cards
		last := pile.Cards()[pile.Size()-RunLength:]
		if isValidRun(last) {
			removed, err := pile.RemoveCardsFrom(pile.Size() - RunLength)
			if err != nil {
				return ErrRemoveCardsWithContext(err)
			}
			g.Completed = append(g.Completed, removed)

			// flip top card if needed
			if err := pile.FlipTopCardIfFaceDown(); err != nil {
				return ErrFlipWithContext(err)
			}
		}
	}
	g.checkWinCondition()
	return nil
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

// checkWinCondition checks if the the number of piles in g.Completed is >= to 8 (The total number of runs to win)
func (g *GameState) checkWinCondition() {
	if g.Won || g.Lost {
		return
	}

	if len(g.Completed) >= TotalRunsToWin {
		g.Won = true
	}
}

// checkLossCondition checks if it is game over for the user
func (g *GameState) checkLossCondition() {
	if g.Won || g.Lost {
		return
	}
	if hasStock(g.Stock) {
		return
	}
	if hasEmptyPile(g.Tableau.Piles) {
		return
	}
	if hasAnyValidMove(g.Tableau.Piles) {
		return
	}
	g.Lost = true
}

// snapshot creates a deep copy of the current GameState for undo history
func (g *GameState) snapshot() GameState {
	snap := GameState{
		Won:   g.Won,
		Lost:  g.Lost,
		Stock: make([]deck.Card, len(g.Stock)),
	}
	copy(snap.Stock, g.Stock)

	// Deep copy tableau piles
	for i := range g.Tableau.Piles {
		snap.Tableau.Piles[i] = g.Tableau.Piles[i].Clone()
	}

	// deep copy completed runs (slice of slices)
	snap.Completed = make([][]CardInPile, len(g.Completed))
	for i, run := range g.Completed {
		snap.Completed[i] = make([]CardInPile, len(run))
		copy(snap.Completed[i], run)
	}

	// don't copy history itself (would cause exponential memory growth!)
	return snap
}

// push history saves the current state before an action
func (g *GameState) pushHistory() {
	snap := g.snapshot()
	g.history = append(g.history, snap)

	// Enforce size limit (FIFO - remove oldest object if needed)
	if len(g.history) > maxHistorySize {
		g.history = g.history[1:]
	}
}

func (g *GameState) Undo() error {
	if len(g.history) == 0 {
		return ErrNoHistory
	}

	// Pop the last state from history
	lastIdx := len(g.history) - 1
	previous := g.history[lastIdx]
	g.history = g.history[:lastIdx]

	// restore the previous state (but preserve remaining history)
	g.Tableau = previous.Tableau
	g.Stock = previous.Stock
	g.Completed = previous.Completed
	g.Won = previous.Won
	g.Lost = previous.Lost

	return nil
}
