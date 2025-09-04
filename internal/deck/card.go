package deck

import "fmt"

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

func (c Card) String() string {
	return fmt.Sprintf("%s of %s", c.Rank, c.Suit)
}
