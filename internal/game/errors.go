package game

import (
	"errors"
	"fmt"
)

// setup errors
var (
	ErrNotEnoughCards    = errors.New("not enough cards to start spider")
	ErrInsufficientStock = errors.New("not enough cards in stock to deal a row")
)

// validation errors
var (
	ErrInvalidSourceIndex      = errors.New("invalid source pile index")
	ErrInvalidDestinationIndex = errors.New("invalid destination pile index")
	ErrSamePileMove            = errors.New("cannot move cards within the same pile")
	ErrInvalidStartIndex       = errors.New("invalid start index")
	ErrNoCardsToMove           = errors.New("no cards to move")
	ErrInvalidSequence         = errors.New("invalid move: sequence not ordered")
	ErrDestinationNotAccepting = errors.New("invalid move: destination cannot accept")
)

// internal errors
var (
	ErrSequenceMismatch  = errors.New("internal error: removed cards don't match expected sequence")
	ErrFlipFailed        = errors.New("failed to flip source card")
	ErrRemoveCardsFailed = errors.New("failed to remove cards from the pile")
)

// error helper functions to better context

func ErrRemoveCardsWithContext(err error) error {
	return fmt.Errorf("%w: %v", ErrRemoveCardsFailed, err)
}

func ErrFlipWithContext(err error) error {
	return fmt.Errorf("%w: %v", ErrFlipFailed, err)
}

// typed errors

type CardFaceDownError struct {
	Index int
}

func (e CardFaceDownError) Error() string {
	return fmt.Sprintf("card at position %d is face down", e.Index)
}
