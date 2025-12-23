package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/staylor11x/spider-solitaire/internal/game"

	"golang.org/x/image/font/basicfont"
)

var uiTextFace = text.NewGoXFace(basicfont.Face7x13)

// drawTableau renders all 10 piles from the view snapshot
func drawTableau(screen *ebiten.Image, view game.GameViewDTO) {
	for i, pile := range view.Tableau {
		x := TableauStartX + i*PileSpacing
		y := TableauStartY
		drawPile(screen, pile, x, y)
	}
}

// drawPile renders a single pile at the given position
func drawPile(screen *ebiten.Image, pile game.PileDTO, x, y int) {
	for i, card := range pile.Cards {
		// stack the cards vertically with a small gap
		cardY := y + i*CardStackGap
		drawCard(screen, card, x, cardY)
	}
}

// drawCard renders a single pile at the given position
func drawCard(screen *ebiten.Image, card game.CardDTO, x, y int) {
	bgColor := CardFaceUpColor
	if !card.FaceUp {
		bgColor = CardFaceDownColor
	}

	// draw card rectangle
	vector.FillRect(screen, float32(x), float32(y), float32(CardWidth), float32(CardHeight), bgColor, false)

	// draw card text
	var cardText string
	if card.FaceUp {
		cardText = formatCard(card)
	} else {
		cardText = "##"
	}

	// center text on the card
	drawOpts := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign:   text.AlignCenter,
			SecondaryAlign: text.AlignCenter,
		},
	}
	drawOpts.GeoM.Translate(float64(x+CardWidth/2), float64(y+CardHeight/2))
	drawOpts.ColorScale.ScaleWithColor(TextColor)

	text.Draw(screen, cardText, uiTextFace, drawOpts)
}

// formatCard converts a CardDTO to a display string (rank + suit)
func formatCard(card game.CardDTO) string {
	rank := rankStr(int(card.Rank))
	suit := suitStr(int(card.Suit))
	return fmt.Sprintf("%s%s", rank, suit)
}

// rankStr converts rank int to display string
func rankStr(r int) string {
	switch r {
	case 1:
		return "A"
	case 11:
		return "J"
	case 12:
		return "Q"
	case 13:
		return "K"
	default:
		return fmt.Sprintf("%d", r)
	}
}

// suitStr converts rank into to display string
func suitStr(s int) string {
	switch s {
	case 0:
		return "S"
	case 1:
		return "H"
	case 2:
		return "D"
	case 3:
		return "C"
	default:
		return "?"
	}
}

func drawStats(screen *ebiten.Image, view game.GameViewDTO) {
	stats := fmt.Sprintf("Stock: %d | Completed: %d | Won: %v | Lost: %v",
		view.StockCount, view.CompletedCount, view.Won, view.Lost)

	drawOpts := &text.DrawOptions{}
	drawOpts.GeoM.Translate(float64(StatsX), float64(StatsY))
	drawOpts.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, stats, uiTextFace, drawOpts)
}
