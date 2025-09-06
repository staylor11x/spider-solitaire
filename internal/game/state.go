package game

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/staylor11x/spider-solitaire/internal/deck"
)

// APPROACH 1: IMMUTABLE GAME STATES
// =================================
// Create new GameState for each operation instead of mutating existing ones

type ImmutableGameState struct {
	Tableau Tableau
	Stock   []deck.Card
	// Add metadata for tracking
	MoveCount int
	LastMove  string
}

// DealRowImmutable returns a new GameState with the row dealt
func (g *ImmutableGameState) DealRowImmutable() (*ImmutableGameState, error) {
	if !g.CanDealRow() {
		return nil, errors.New("not enough cards in stock to deal a full row")
	}

	// Create new state (deep copy)
	newState := &ImmutableGameState{
		Tableau:   g.Tableau, // This needs deep copying - see below
		MoveCount: g.MoveCount + 1,
		LastMove:  "DealRow",
	}

	// Copy stock and remove cards
	newState.Stock = make([]deck.Card, len(g.Stock)-TableauPiles)
	copy(newState.Stock, g.Stock[:len(g.Stock)-TableauPiles])

	// Deal cards to new tableau
	for i := range TableauPiles {
		card := g.Stock[len(g.Stock)-1-i]
		newState.Tableau.Piles[i].AddCard(card, true)
	}

	return newState, nil
}

// APPROACH 2: COMMAND PATTERN WITH UNDO STACK
// ==========================================
// Each operation is a command that knows how to undo itself

type GameCommand interface {
	Execute(gs *GameState) error
	Undo(gs *GameState) error
	Description() string
}

type DealRowCommand struct {
	dealtCards []deck.Card // Store what was dealt for undo
}

func NewDealRowCommand() *DealRowCommand {
	return &DealRowCommand{}
}

func (cmd *DealRowCommand) Execute(gs *GameState) error {
	if !gs.CanDealRow() {
		return errors.New("not enough cards in stock to deal a full row")
	}

	// Store dealt cards for undo
	cmd.dealtCards = make([]deck.Card, TableauPiles)
	for i := range TableauPiles {
		cmd.dealtCards[i] = gs.Stock[len(gs.Stock)-1-i]
	}

	// Execute the deal
	for i := range TableauPiles {
		card := gs.Stock[len(gs.Stock)-1]
		gs.Stock = gs.Stock[:len(gs.Stock)-1]
		gs.Tableau.Piles[i].AddCard(card, true)
	}

	return nil
}

func (cmd *DealRowCommand) Undo(gs *GameState) error {
	if len(cmd.dealtCards) != TableauPiles {
		return errors.New("invalid undo state")
	}

	// Remove top card from each pile and return to stock
	for i := range TableauPiles {
		pile := &gs.Tableau.Piles[i]
		if pile.Size() == 0 {
			return errors.New("cannot undo: pile is empty")
		}

		// Verify the top card matches what we dealt
		topCard, err := pile.TopCard()
		if err != nil {
			return err
		}

		if topCard.Card != cmd.dealtCards[i] {
			return fmt.Errorf("undo validation failed: expected %v, got %v",
				cmd.dealtCards[i], topCard.Card)
		}

		// Remove the card (this method would need to be added to Pile)
		pile.RemoveTopCard()

		// Add back to stock in reverse order
		gs.Stock = append(gs.Stock, cmd.dealtCards[TableauPiles-1-i])
	}

	return nil
}

func (cmd *DealRowCommand) Description() string {
	return "Deal row from stock"
}

// Game with command history
type GameWithHistory struct {
	State      *GameState
	History    []GameCommand
	HistoryPos int // For redo functionality
}

func NewGameWithHistory() (*GameWithHistory, error) {
	initialState, err := DealInitialGame()
	if err != nil {
		return nil, err
	}

	return &GameWithHistory{
		State:      initialState,
		History:    make([]GameCommand, 0),
		HistoryPos: 0,
	}, nil
}

func (gwh *GameWithHistory) ExecuteCommand(cmd GameCommand) error {
	// Execute the command
	if err := cmd.Execute(gwh.State); err != nil {
		return err
	}

	// Truncate history if we're in the middle (user made new move after undo)
	gwh.History = gwh.History[:gwh.HistoryPos]

	// Add command to history
	gwh.History = append(gwh.History, cmd)
	gwh.HistoryPos++

	return nil
}

func (gwh *GameWithHistory) Undo() error {
	if gwh.HistoryPos == 0 {
		return errors.New("nothing to undo")
	}

	// Get the last command and undo it
	lastCmd := gwh.History[gwh.HistoryPos-1]
	if err := lastCmd.Undo(gwh.State); err != nil {
		return fmt.Errorf("undo failed: %w", err)
	}

	gwh.HistoryPos--
	return nil
}

func (gwh *GameWithHistory) Redo() error {
	if gwh.HistoryPos >= len(gwh.History) {
		return errors.New("nothing to redo")
	}

	// Re-execute the command
	cmd := gwh.History[gwh.HistoryPos]
	if err := cmd.Execute(gwh.State); err != nil {
		return fmt.Errorf("redo failed: %w", err)
	}

	gwh.HistoryPos++
	return nil
}

// APPROACH 3: SNAPSHOT-BASED UNDO
// ==============================
// Take full snapshots at key points

type GameStateSnapshot struct {
	State     *GameState
	Timestamp int64
	MoveDesc  string
}

type GameWithSnapshots struct {
	currentState *GameState
	snapshots    []GameStateSnapshot
	maxSnapshots int
}

func NewGameWithSnapshots(maxSnapshots int) (*GameWithSnapshots, error) {
	initialState, err := DealInitialGame()
	if err != nil {
		return nil, err
	}

	return &GameWithSnapshots{
		currentState: initialState,
		snapshots:    make([]GameStateSnapshot, 0),
		maxSnapshots: maxSnapshots,
	}, nil
}

// DeepCopyGameState creates a complete copy of the game state
func DeepCopyGameState(original *GameState) (*GameState, error) {
	// Use JSON marshaling for deep copy (simple but not most efficient)
	data, err := json.Marshal(original)
	if err != nil {
		return nil, err
	}

	var copy GameState
	if err := json.Unmarshal(data, &copy); err != nil {
		return nil, err
	}

	return &copy, nil
}

func (gws *GameWithSnapshots) TakeSnapshot(description string) error {
	// Create deep copy of current state
	snapshot, err := DeepCopyGameState(gws.currentState)
	if err != nil {
		return err
	}

	// Add to snapshots
	gws.snapshots = append(gws.snapshots, GameStateSnapshot{
		State:    snapshot,
		MoveDesc: description,
	})

	// Limit snapshot history
	if len(gws.snapshots) > gws.maxSnapshots {
		gws.snapshots = gws.snapshots[1:]
	}

	return nil
}

func (gws *GameWithSnapshots) RestoreLastSnapshot() error {
	if len(gws.snapshots) == 0 {
		return errors.New("no snapshots available")
	}

	// Get last snapshot
	lastSnapshot := gws.snapshots[len(gws.snapshots)-1]

	// Restore state
	restored, err := DeepCopyGameState(lastSnapshot.State)
	if err != nil {
		return err
	}

	gws.currentState = restored

	// Remove the used snapshot
	gws.snapshots = gws.snapshots[:len(gws.snapshots)-1]

	return nil
}

// APPROACH 4: STATE VALIDATION
// ===========================
// Add validation methods to catch inconsistencies

func (gs *GameState) Validate() error {
	// Check total card count
	totalCards := len(gs.Stock)
	for i := range TableauPiles {
		totalCards += gs.Tableau.Piles[i].Size()
	}

	if totalCards != TotalSpiderCards {
		return fmt.Errorf("invalid total cards: expected %d, got %d",
			TotalSpiderCards, totalCards)
	}

	// Check for duplicate cards (more complex validation)
	cardCounts := make(map[deck.Card]int)

	// Count stock cards
	for _, card := range gs.Stock {
		cardCounts[card]++
	}

	// Count tableau cards
	for i := range TableauPiles {
		for _, cardInPile := range gs.Tableau.Piles[i].Cards() {
			cardCounts[cardInPile.Card]++
		}
	}

	// In a 2-deck Spider game, each card should appear exactly twice
	for card, count := range cardCounts {
		if count != 2 {
			return fmt.Errorf("card %v appears %d times, expected 2", card, count)
		}
	}

	return nil
}

// Extension to Pile struct needed for command pattern
func (p *Pile) RemoveTopCard() error {
	if len(p.cards) == 0 {
		return errors.New("cannot remove card from empty pile")
	}
	p.cards = p.cards[:len(p.cards)-1]
	return nil
}

// RECOMMENDED HYBRID APPROACH
// ==========================
type RecommendedGameManager struct {
	state   *GameState
	history []GameCommand
	maxUndo int
}

func NewRecommendedGameManager(maxUndo int) (*RecommendedGameManager, error) {
	initialState, err := DealInitialGame()
	if err != nil {
		return nil, err
	}

	return &RecommendedGameManager{
		state:   initialState,
		history: make([]GameCommand, 0, maxUndo),
		maxUndo: maxUndo,
	}, nil
}

func (rgm *RecommendedGameManager) ExecuteMove(cmd GameCommand) error {
	// Validate current state before move
	if err := rgm.state.Validate(); err != nil {
		return fmt.Errorf("pre-move validation failed: %w", err)
	}

	// Execute command
	if err := cmd.Execute(rgm.state); err != nil {
		return err
	}

	// Validate state after move
	if err := rgm.state.Validate(); err != nil {
		// Try to undo the command
		cmd.Undo(rgm.state)
		return fmt.Errorf("post-move validation failed: %w", err)
	}

	// Add to history
	rgm.history = append(rgm.history, cmd)
	if len(rgm.history) > rgm.maxUndo {
		rgm.history = rgm.history[1:]
	}

	return nil
}

func (rgm *RecommendedGameManager) Undo() error {
	if len(rgm.history) == 0 {
		return errors.New("nothing to undo")
	}

	// Get last command
	lastCmd := rgm.history[len(rgm.history)-1]

	// Undo it
	if err := lastCmd.Undo(rgm.state); err != nil {
		return err
	}

	// Validate state after undo
	if err := rgm.state.Validate(); err != nil {
		// Re-execute to restore state
		lastCmd.Execute(rgm.state)
		return fmt.Errorf("undo validation failed: %w", err)
	}

	// Remove from history
	rgm.history = rgm.history[:len(rgm.history)-1]

	return nil
}
