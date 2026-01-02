package ui

import (
	"io/fs"
	"testing"

	"github.com/staylor11x/spider-solitaire/internal/assets"
	"github.com/stretchr/testify/assert"
)

func TestFilenameForCard_Mapping(t *testing.T) {
	cases := []struct {
		suit int
		rank int
		want string
	}{
		{1, 10, "10_of_hearts.png"},
		{2, 12, "queen_of_diamonds.png"},
		{0, 1, "ace_of_spades.png"},
		{3, 13, "king_of_clubs.png"},
	}
	for _, c := range cases {
		got := filenameForCard(c.suit, c.rank)
		assert.Equal(t, c.want, got)
	}
}

func TestCardAssets_Exists(t *testing.T) {
	// Back image
	if _, err := fs.ReadFile(assets.Files, "images/card-back/card-back-red.png"); err != nil {
		t.Fatalf("missing back image: %v", err)
	}

	// All 52 faces
	for s := 0; s < 4; s++ {
		for r := 1; r <= 13; r++ {
			path := "images/cards/" + filenameForCard(s, r)
			if _, err := fs.ReadFile(assets.Files, path); err != nil {
				t.Errorf("missing card asset: %s (%v)", path, err)
			}
		}
	}
}

func TestCardAtlas_InvalidCard(t *testing.T) {
	atlas, err := NewCardAtlas(assets.Files)
	assert.NoError(t, err)

	_, err = atlas.Card(99, 1) // invalid suit
	assert.Error(t, err)

	_, err = atlas.Card(1, 99) // invalid rank
	assert.Error(t, err)
}

func TestCardAtlas_Caching(t *testing.T) {
	atlas, err := NewCardAtlas(assets.Files)
	assert.NoError(t, err)

	img1, err := atlas.Card(1, 10)
	assert.NoError(t, err)

	img2, err := atlas.Card(1, 10) // same card
	assert.NoError(t, err)

	assert.Same(t, img1, img2) // should be the same pointer
}
