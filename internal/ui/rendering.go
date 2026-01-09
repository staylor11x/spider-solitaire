package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/staylor11x/spider-solitaire/internal/deck"
	"github.com/staylor11x/spider-solitaire/internal/game"

	"golang.org/x/image/font/basicfont"
)

var uiTextFace = text.NewGoXFace(basicfont.Face7x13)

// drawTableau renders all 10 piles from the view snapshot
func drawTableau(screen *ebiten.Image, view game.GameViewDTO, atlas *CardAtlas) {
	for i, pile := range view.Tableau {
		x := TableauStartX + i*PileSpacing
		y := TableauStartY
		drawPile(screen, pile, x, y, atlas)
	}
}

// drawPile renders a single pile at the given position
func drawPile(screen *ebiten.Image, pile game.PileDTO, x, y int, atlas *CardAtlas) {
	// If the pile is empty, render a faint placeholder to indicate a valid drop target
	if len(pile.Cards) == 0 {
		drawEmptyPilePlaceholder(screen, x, y)
		return
	}
	for i, card := range pile.Cards {
		// stack the cards vertically with a small gap
		cardY := y + i*CardStackGap
		drawCard(screen, card, x, cardY, atlas)
	}
}

// drawCard renders a single card at the given position
func drawCard(screen *ebiten.Image, card game.CardDTO, x, y int, atlas *CardAtlas) {

	if atlas != nil {
		var img *ebiten.Image
		if card.FaceUp {
			if im, err := atlas.Card(int(card.Suit), int(card.Rank)); err == nil {
				img = im
			}
		} else {
			img = atlas.Back()
		}
		if img != nil {
			w := img.Bounds().Dx()
			h := img.Bounds().Dy()
			opts := &ebiten.DrawImageOptions{}
			// Scale to logical card size if asset size differs
			sx := float64(CardWidth) / float64(w)
			sy := float64(CardHeight) / float64(h)
			opts.GeoM.Scale(sx, sy)
			opts.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(img, opts)
			return
		}
	}

	// Fallback: vector rect + text
	bgColor := CardFaceUpColor
	if !card.FaceUp {
		bgColor = CardFaceDownColor
	}

	// card rectangle
	vector.FillRect(screen, float32(x), float32(y), float32(CardWidth), float32(CardHeight), bgColor, false)

	// card text
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
	c := deck.Card{Suit: deck.Suit(card.Suit), Rank: deck.Rank(card.Rank)}
	return fmt.Sprintf("%s%s", c.RankName(), c.SuitName())
}

// drawStats renders stock and completed counts at the top-left
func drawStats(screen *ebiten.Image, view game.GameViewDTO) {
	stats := fmt.Sprintf("Stock: %d | Completed: %d | Won: %v | Lost: %v",
		view.StockCount, view.CompletedCount, view.Won, view.Lost)

	drawOpts := &text.DrawOptions{}
	drawOpts.GeoM.Translate(float64(StatsX), float64(StatsY))
	drawOpts.ColorScale.ScaleWithColor(color.White)

	text.Draw(screen, stats, uiTextFace, drawOpts)
}

// drawError shows an ephemeral error message at the top-right (centered in pill)
func drawError(screen *ebiten.Image, msg string) {

	// near top right
	const margin = 20
	w := screen.Bounds().Dx() // return the width of the screen

	// background pill
	bgW, bgH := 380, 24
	bgX := float32(w - margin - bgW)
	bgY := float32(StatsY)

	// Semi-transparent background
	vector.FillRect(screen, bgX, bgY, float32(bgW), float32(bgH), color.RGBA{0, 0, 0, 160}, false)

	// center the text within the pill
	opts := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign:   text.AlignCenter,
			SecondaryAlign: text.AlignCenter,
		},
	}
	// translate to the pills center
	opts.GeoM.Translate(float64(bgX)+float64(bgW)/2, float64(bgY)+float64(bgH)/2)
	opts.ColorScale.ScaleWithColor(color.White)

	text.Draw(screen, msg, uiTextFace, opts)
}

// drawSelectionOverlay highlights the selected suffix (from selectedIndex to top) on a pile.
func drawSelectionOverlay(screen *ebiten.Image, view game.GameViewDTO, pileIdx, selectedIndex int, col color.RGBA) {
	if pileIdx < 0 || pileIdx >= len(view.Tableau) {
		return
	}
	pile := view.Tableau[pileIdx]
	if selectedIndex < 0 || selectedIndex >= len(pile.Cards) {
		return
	}
	x := TableauStartX + pileIdx*PileSpacing
	y := TableauStartY

	for i := selectedIndex; i < len(pile.Cards); i++ {
		cy := y + i*CardStackGap
		vector.FillRect(screen, float32(x), float32(cy), float32(CardWidth), float32(CardHeight), col, false)
	}
}

func drawEmptyPilePlaceholder(screen *ebiten.Image, x, y int) {
	const borderWidth = 2
	// Faint fill and border for visibility on table felt
	fill := color.RGBA{R: 0, G: 0, B: 0, A: 30}
	border := color.RGBA{R: 255, G: 255, B: 255, A: 90}

	vector.FillRect(screen, float32(x), float32(y), float32(CardWidth), float32(CardHeight), fill, false)
	vector.StrokeRect(screen, float32(x), float32(y), float32(CardWidth), float32(CardHeight), borderWidth, border, false)

}

// drawWinLossOverlay darkens the background and renders a centered message
func drawWinLossOverlay(screen *ebiten.Image, msg string) {
	b := screen.Bounds()
	w, h := b.Dx(), b.Dy()

	// Dim the whole screen for contrast
	vector.FillRect(screen, 0, 0, float32(w), float32(h), color.RGBA{0, 0, 0, 160}, false)

	// center the text
	opts := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign:   text.AlignCenter,
			SecondaryAlign: text.AlignCenter,
		},
	}
	opts.GeoM.Translate(float64(w)/2, float64(h)/2)
	opts.ColorScale.ScaleWithColor(color.RGBA{255, 255, 255, 255})

	text.Draw(screen, msg, uiTextFace, opts)
}
