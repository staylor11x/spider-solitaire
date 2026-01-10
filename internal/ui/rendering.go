package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/staylor11x/spider-solitaire/internal/deck"
	"github.com/staylor11x/spider-solitaire/internal/game"
)

// drawTableau renders all 10 piles from the view snapshot
func drawTableau(screen *ebiten.Image, view game.GameViewDTO, atlas *CardAtlas, theme *Theme) {
	for i, pile := range view.Tableau {
		x := theme.Layout.TableauStartX + i*theme.Layout.PileSpacing
		y := theme.Layout.TableauStartY
		drawPile(screen, pile, x, y, atlas, theme)
	}
}

// drawPile renders a single pile at the given position
func drawPile(screen *ebiten.Image, pile game.PileDTO, x, y int, atlas *CardAtlas, theme *Theme) {
	// If the pile is empty, render a faint placeholder to indicate a valid drop target
	if len(pile.Cards) == 0 {
		drawEmptyPilePlaceholder(screen, x, y, theme)
		return
	}
	for i, card := range pile.Cards {
		// stack the cards vertically with a small gap
		cardY := y + i*theme.Layout.CardStackGap
		drawCard(screen, card, x, cardY, atlas, theme)
	}
}

// drawCard renders a single card at the given position
func drawCard(screen *ebiten.Image, card game.CardDTO, x, y int, atlas *CardAtlas, theme *Theme) {

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
			sx := float64(theme.Layout.CardWidth) / float64(w)
			sy := float64(theme.Layout.CardHeight) / float64(h)
			opts.GeoM.Scale(sx, sy)
			opts.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(img, opts)
			return
		}
	}

	// Fallback: vector rect + text
	bgColor := theme.Colors.CardFaceUp
	if !card.FaceUp {
		bgColor = theme.Colors.CardFaceDown
	}

	// card rectangle
	vector.FillRect(screen, float32(x), float32(y), float32(theme.Layout.CardWidth), float32(theme.Layout.CardHeight), bgColor, false)

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
	drawOpts.GeoM.Translate(float64(x+theme.Layout.CardWidth/2), float64(y+theme.Layout.CardHeight/2))
	drawOpts.ColorScale.ScaleWithColor(theme.Colors.CardText)

	text.Draw(screen, cardText, theme.Font, drawOpts)
}

// formatCard converts a CardDTO to a display string (rank + suit)
func formatCard(card game.CardDTO) string {
	c := deck.Card{Suit: deck.Suit(card.Suit), Rank: deck.Rank(card.Rank)}
	return fmt.Sprintf("%s%s", c.RankName(), c.SuitName())
}

// drawStats renders stock and completed counts at the top-left
func drawStats(screen *ebiten.Image, view game.GameViewDTO, theme *Theme) {
	stats := fmt.Sprintf("Stock: %d | Completed: %d | Won: %v | Lost: %v",
		view.StockCount, view.CompletedCount, view.Won, view.Lost)

	drawOpts := &text.DrawOptions{}
	drawOpts.GeoM.Translate(float64(theme.Layout.StatsX), float64(theme.Layout.StatsY))
	drawOpts.ColorScale.ScaleWithColor(theme.Colors.HelpOverlayText)

	text.Draw(screen, stats, theme.Font, drawOpts)
}

// drawError shows an ephemeral error message at the top-right (centered in pill)
func drawError(screen *ebiten.Image, msg string, theme *Theme) {

	// near top right
	const margin = 20
	w := screen.Bounds().Dx() // return the width of the screen

	// background pill
	bgW, bgH := 380, 24
	bgX := float32(w - margin - bgW)
	bgY := float32(theme.Layout.StatsY)

	// Semi-transparent background
	vector.FillRect(screen, bgX, bgY, float32(bgW), float32(bgH), theme.Colors.ErrorPillBG, false)

	// center the text within the pill
	opts := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign:   text.AlignCenter,
			SecondaryAlign: text.AlignCenter,
		},
	}
	// translate to the pills center
	opts.GeoM.Translate(float64(bgX)+float64(bgW)/2, float64(bgY)+float64(bgH)/2)
	opts.ColorScale.ScaleWithColor(theme.Colors.ErrorPillText)

	text.Draw(screen, msg, theme.Font, opts)
}

// drawSelectionOverlay highlights the selected suffix (from selectedIndex to top) on a pile.
func drawSelectionOverlay(screen *ebiten.Image, view game.GameViewDTO, pileIdx, selectedIndex int, theme *Theme) {
	if pileIdx < 0 || pileIdx >= len(view.Tableau) {
		return
	}
	pile := view.Tableau[pileIdx]
	if selectedIndex < 0 || selectedIndex >= len(pile.Cards) {
		return
	}
	x := theme.Layout.TableauStartX + pileIdx*theme.Layout.PileSpacing
	y := theme.Layout.TableauStartY

	for i := selectedIndex; i < len(pile.Cards); i++ {
		cy := y + i*theme.Layout.CardStackGap
		vector.FillRect(screen, float32(x), float32(cy), float32(theme.Layout.CardWidth), float32(theme.Layout.CardHeight), theme.Colors.SelectionOverlay, false)
	}
}

func drawEmptyPilePlaceholder(screen *ebiten.Image, x, y int, theme *Theme) {
	const borderWidth = 2
	// Faint fill and border for visibility on table felt
	fill := theme.Colors.PlaceholderBG
	border := theme.Colors.PlaceholderBorder

	vector.FillRect(screen, float32(x), float32(y), float32(theme.Layout.CardWidth), float32(theme.Layout.CardHeight), fill, false)
	vector.StrokeRect(screen, float32(x), float32(y), float32(theme.Layout.CardWidth), float32(theme.Layout.CardHeight), borderWidth, border, false)

}

// drawWinLossOverlay darkens the background and renders a centered message
func drawWinLossOverlay(screen *ebiten.Image, msg string, theme *Theme) {
	b := screen.Bounds()
	w, h := b.Dx(), b.Dy()

	// Dim the whole screen for contrast
	vector.FillRect(screen, 0, 0, float32(w), float32(h), theme.Colors.HelpOverlayBG, false)

	// center the text
	opts := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign:   text.AlignCenter,
			SecondaryAlign: text.AlignCenter,
		},
	}
	opts.GeoM.Translate(float64(w)/2, float64(h)/2)
	opts.ColorScale.ScaleWithColor(theme.Colors.HelpOverlayText)

	text.Draw(screen, msg, theme.Font, opts)
}

func drawHelpOverlay(screen *ebiten.Image, theme *Theme) {
	b := screen.Bounds()
	w, h := b.Dx(), b.Dy()

	vector.FillRect(screen, 0, 0, float32(w), float32(h), theme.Colors.HelpOverlayBG, false)
	helpLines := []string{
		"Controls",
		"",
		"Click - Select/Move Cards",
		"[D] - Deal Row",
		"[R] - Reset Game",
		"[H] - Toggle Help",
		"",
		"Press [H] to close",
	}

	lineHeight := theme.Font.Metrics().HLineGap + theme.Font.Metrics().HAscent + theme.Font.Metrics().HDescent
	totalHeight := float64(len(helpLines)) * lineHeight
	startY := (float64(h) - totalHeight) / 2

	for i, line := range helpLines {
		opts := &text.DrawOptions{
			LayoutOptions: text.LayoutOptions{
				PrimaryAlign: text.AlignCenter,
			},
		}
		opts.GeoM.Translate(float64(w)/2, startY+float64(i)*lineHeight)
		opts.ColorScale.ScaleWithColor(theme.Colors.HelpOverlayText)
		text.Draw(screen, line, theme.Font, opts)
	}

}
