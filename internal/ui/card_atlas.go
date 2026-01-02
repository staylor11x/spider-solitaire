package ui

import (
	"bytes"
	"fmt"
	"image/png"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/staylor11x/spider-solitaire/internal/deck"
)

type CardAtlas struct {
	fs    fs.FS
	cache map[string]*ebiten.Image
	back  *ebiten.Image
}

func NewCardAtlas(fsys fs.FS) (*CardAtlas, error) {
	a := &CardAtlas{
		fs:    fsys,
		cache: make(map[string]*ebiten.Image),
	}
	b, err := a.loadPNG("images/card-back/card-back-red.png")
	if err != nil {
		return nil, fmt.Errorf("load card back: %w", err)
	}
	a.back = b
	return a, nil
}

func (a *CardAtlas) Back() *ebiten.Image {
	return a.back
}

func (a *CardAtlas) Card(suit, rank int) (*ebiten.Image, error) {
	name := fmt.Sprintf("images/cards/%s", filenameForCard(suit, rank))
	if img, ok := a.cache[name]; ok {
		return img, nil
	}
	img, err := a.loadPNG(name)
	if err != nil {
		return nil, err
	}
	a.cache[name] = img
	return img, nil
}

func (a *CardAtlas) loadPNG(path string) (*ebiten.Image, error) {
	data, err := fs.ReadFile(a.fs, path)
	if err != nil {
		return nil, err
	}
	src, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(src), nil
}

func filenameForCard(suit, rank int) string {
	c := deck.Card{Suit: deck.Suit(suit), Rank: deck.Rank(rank)}
	return fmt.Sprintf("%s_of_%s.png", c.RankName(), c.SuitName())
}
