package deck

import (
	"fmt"
	"strings"
)

type Suit int

const (
	Spades Suit = iota
	Hearts
	Diamonds
	Clubs
)

func (s Suit) String() string {
	suits := [...]string{"Spades", "Hearts", "Diamonds", "Clubs"}

	if int(s) < 0 || int(s) >= len(suits) {
		return "Unknown"
	}
	return suits[s]
}

type Rank int

const (
	Ace Rank = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

func (r Rank) String() string {
	ranks := [...]string{
		"Ace", "2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "Queen", "King"}

	if r < Ace || r > King {
		return "Unknown"
	}
	return ranks[r-1]
}

type Card struct {
	Suit Suit
	Rank Rank
}

// String returns a human-readable card name (e.g. "Ace of Spades")
func (c Card) String() string {
	return fmt.Sprintf("%s of %s", c.Rank, c.Suit)
}

// RankName returns the rank as a word or number (lowercase)
func (c Card) RankName() string {
	return strings.ToLower(c.Rank.String())
}

// SuitName returns the suit as a word (lowercase)
func (c Card) SuitName() string {
	return strings.ToLower(c.Suit.String())
}

// RankSymbol return the rank as a display character
func (c Card) RankSymbol() string {
	switch c.Rank {
	case Ace:
		return "A"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	default:
		return fmt.Sprintf("%d", c.Rank)
	}
}

// SuitSymbol return the Unicode suit symbol
func (c Card) SuitSymbol() string {
	switch c.Suit {
	case Spades:
		return "♠"
	case Hearts:
		return "♥"
	case Diamonds:
		return "♦"
	case Clubs:
		return "♣"
	default:
		return "?"
	}
}
