package game

import (
	"testing"

	"github.com/staylor11x/spider-solitaire/internal/deck"
	"github.com/stretchr/testify/assert"
)

func TestGameStateView_FidelityAndDefensiveCopies(t *testing.T) {

	p0 := Pile{
		cards: []CardInPile{
			{Card: deck.Card{Suit: deck.Spades, Rank: deck.King}, FaceUp: true},
			{Card: deck.Card{Suit: deck.Spades, Rank: deck.Queen}, FaceUp: true},
			{Card: deck.Card{Suit: deck.Spades, Rank: deck.Jack}, FaceUp: true},
		},
	}
	p1 := Pile{
		cards: []CardInPile{
			{Card: deck.Card{Suit: deck.Hearts, Rank: deck.Ace}, FaceUp: true},
		},
	}

	var tableau Tableau
	tableau.Piles[0] = p0
	tableau.Piles[1] = p1

	stock := []deck.Card{
		{Suit: deck.Hearts, Rank: deck.Two},
		{Suit: deck.Clubs, Rank: deck.Three},
	}

	completed := make([][]CardInPile, 2)

	g := &GameState{
		Tableau:   tableau,
		Stock:     stock,
		Completed: completed,
		Won:       false,
		Lost:      false,
	}

	view := g.View()

	// Check the gamestate elements
	assert.Equal(t, len(stock), view.StockCount)
	assert.Equal(t, len(completed), view.CompletedCount)
	assert.False(t, view.Lost)
	assert.False(t, view.Won)
	assert.Equal(t, len(g.Tableau.Piles), len(view.Tableau))

	// Check the piles
	got0 := view.Tableau[0].Cards
	want0 := g.Tableau.Piles[0].cards
	assert.Equal(t, len(want0), len(got0))
	for i := range got0 {
		assert.Equal(t, int(want0[i].Card.Rank), int(got0[i].Rank))
		assert.Equal(t, int(want0[i].Card.Suit), int(got0[i].Suit))
		assert.Equal(t, want0[i].FaceUp, got0[i].FaceUp)
	}
	got1 := view.Tableau[1].Cards
	want1 := g.Tableau.Piles[1].cards
	assert.Equal(t, len(want1), len(got1))
	for i := range got1 {
		assert.Equal(t, int(want1[i].Card.Rank), int(got1[i].Rank))
		assert.Equal(t, int(want1[i].Card.Suit), int(got1[i].Suit))
		assert.Equal(t, want1[i].FaceUp, got1[i].FaceUp)
	}

	// Defensive: mutating the DTO should not affect the engine
	view.Tableau[0].Cards[0].FaceUp = false
	assert.True(t, g.Tableau.Piles[0].cards[0].FaceUp)

	// Defensive: mutating engine after snapshot should not retroactively change DTO
	g.Tableau.Piles[0].cards[1].FaceUp = false
	assert.True(t, view.Tableau[0].Cards[1].FaceUp)

	// A fresh snapshot should reflect the engine mutation
	view2 := g.View()
	assert.False(t, view2.Tableau[0].Cards[1].FaceUp)
}

func TestGameStateView_EmptyState(t *testing.T) {
	var g GameState // nothing

	view := g.View()

	assert.Equal(t, 0, view.StockCount)
	assert.Equal(t, 0, view.CompletedCount)
	assert.False(t, view.Lost)
	assert.False(t, view.Won)

	// all piles should be present and empty
	assert.Equal(t, len(g.Tableau.Piles), len(view.Tableau))

	for i := range view.Tableau {
		assert.Len(t, view.Tableau[i].Cards, 0)
	}
}
