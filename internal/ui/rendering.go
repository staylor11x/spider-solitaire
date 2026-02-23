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
func drawTableau(screen *ebiten.Image, view game.GameViewDTO, atlas *CardAtlas, theme *Theme, selectedPile, selectedIndex, hoveredPile, hoveredCardIdx int) {
	// When a selection is active, suppress hover overlays to avoid visual noise
	selectionActive := selectedPile >= 0 && selectedIndex >= 0

	for i, pile := range view.Tableau {
		x := theme.Layout.TableauStartX + i*theme.Layout.PileSpacing
		y := theme.Layout.TableauStartY
		isSelected := (i == selectedPile) // is this pile the selected one?
		isHovered := (i == hoveredPile)   // is this pile hovered?
		hvdCardIdx := -1
		if isHovered && !selectionActive {
			hvdCardIdx = hoveredCardIdx
		}
		drawPile(screen, pile, x, y, atlas, theme, isSelected, selectedIndex, hvdCardIdx)
	}
}

// drawPile renders a single pile at the given position
// If isSelected is true, cards from selectedIndex onwards are skipped (they'll be drawn lifted in drawSelectionOverlay)
// hoveredCardIdx is the index of the card being hovered (-1 for none), overlay is drawn on that card only
func drawPile(screen *ebiten.Image, pile game.PileDTO, x, y int, atlas *CardAtlas, theme *Theme, isSelected bool, selectedIndex int, hoveredCardIdx int) {
	// If the pile is empty, render a faint placeholder to indicate a valid drop target
	if len(pile.Cards) == 0 {
		drawEmptyPilePlaceholder(screen, x, y, theme)
		return
	}

	// Only show hover if the hovered card itself is face-up
	// Don't highlight anything when hovering over face-down cards
	showHover := false
	movableHoverEnd := -1
	if hoveredCardIdx >= 0 && hoveredCardIdx < len(pile.Cards) {
		showHover = pile.Cards[hoveredCardIdx].FaceUp
		if showHover {
			movableHoverEnd = computeMovableHoverEnd(pile.Cards, hoveredCardIdx)
		}
	}

	layout := computeTableauPileLayout(theme, len(pile.Cards))

	for i, card := range pile.Cards {
		// Skip cards that are part of a selection (they'll be drawn lifted by drawSelectionOverlay)
		if isSelected && i >= selectedIndex {
			continue
		}
		cardY := layout.CardY[i]
		drawCard(screen, card, x, cardY, atlas, theme)
		// Draw hover overlay only over the movable sequence from the hovered card.
		if showHover && movableHoverEnd >= hoveredCardIdx && i >= hoveredCardIdx && i <= movableHoverEnd {
			vector.FillRect(screen, float32(x), float32(cardY), float32(theme.Layout.CardWidth), float32(theme.Layout.CardHeight), theme.Colors.HoverOverlay, false)
		}
	}
}

// computeMovableHoverEnd returns the top index (len(cards)-1) only when the
// sequence from startIdx to top is fully movable.
// A movable sequence is face-up, same suit, descending rank.
// Returns -1 when startIdx is invalid, not face-up, or blocked before the top.
func computeMovableHoverEnd(cards []game.CardDTO, startIdx int) int {
	if startIdx < 0 || startIdx >= len(cards) {
		return -1
	}

	if !cards[startIdx].FaceUp {
		return -1
	}

	for i := startIdx; i < len(cards)-1; i++ {
		current := cards[i]
		next := cards[i+1]

		if !next.FaceUp {
			return -1
		}

		if current.Suit != next.Suit {
			return -1
		}

		if current.Rank != next.Rank+1 {
			return -1
		}
	}

	return len(cards) - 1
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
// Cards are lifted 8 pixels upward and outlined with a goldenrod border for visual feedback.
func drawSelectionOverlay(screen *ebiten.Image, view game.GameViewDTO, pileIdx, selectedIndex int, atlas *CardAtlas, theme *Theme) {
	if pileIdx < 0 || pileIdx >= len(view.Tableau) {
		return
	}
	pile := view.Tableau[pileIdx]
	if selectedIndex < 0 || selectedIndex >= len(pile.Cards) {
		return
	}
	x := theme.Layout.TableauStartX + pileIdx*theme.Layout.PileSpacing
	layout := computeTableauPileLayout(theme, len(pile.Cards))

	for i := selectedIndex; i < len(pile.Cards); i++ {
		cy := layout.CardY[i] - theme.Layout.SelectionLiftPx
		// Redraw the card at the lifted position
		drawCard(screen, pile.Cards[i], x, cy, atlas, theme)
		// Gold tint overlay
		vector.FillRect(screen, float32(x), float32(cy), float32(theme.Layout.CardWidth), float32(theme.Layout.CardHeight), theme.Colors.SelectionOverlay, false)
		// Goldenrod border for contrast
		vector.StrokeRect(screen, float32(x), float32(cy), float32(theme.Layout.CardWidth), float32(theme.Layout.CardHeight), float32(theme.Layout.SelectionBorderPx), theme.Colors.SelectionBorder, false)
	}
}

func drawEmptyPilePlaceholder(screen *ebiten.Image, x, y int, theme *Theme) {
	// Faint fill and border for visibility on table felt
	fill := theme.Colors.PlaceholderBG
	border := theme.Colors.PlaceholderBorder

	vector.FillRect(screen, float32(x), float32(y), float32(theme.Layout.CardWidth), float32(theme.Layout.CardHeight), fill, false)
	vector.StrokeRect(screen, float32(x), float32(y), float32(theme.Layout.CardWidth), float32(theme.Layout.CardHeight), float32(theme.Layout.PlaceholderBorderPx), border, false)

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
		"[U] - Undo Move",
		"[R] - Reset Game",
		"[H] - Toggle Help",
		"[ESC] - Cancel Selection / Close Help",
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

// drawStockPile renders the stock pile visual in the bottom-right corner
func drawStockPile(screen *ebiten.Image, stockCount int, atlas *CardAtlas, theme *Theme, isHovered bool) {
	// Position: bottom-right corner with 20px margin
	stockX := 1120
	stockY := 580

	// If stock is empty, show placeholder
	if stockCount == 0 {
		drawEmptyPilePlaceholder(screen, stockX, stockY, theme)
		return
	}

	// Determine number of layers based on stock count for visual depletion
	var offsets []int
	if stockCount > 20 {
		offsets = []int{4, 2, 0} // 3 layers
	} else if stockCount > 10 {
		offsets = []int{2, 0} // 2 layers
	} else {
		offsets = []int{0} // 1 layer
	}

	// Draw layered card backs
	for _, offset := range offsets {
		x := stockX + offset
		y := stockY + offset

		back := atlas.Back()
		if back != nil {
			w := back.Bounds().Dx()
			h := back.Bounds().Dy()
			opts := &ebiten.DrawImageOptions{}
			// Scale to logical card size
			sx := float64(theme.Layout.CardWidth) / float64(w)
			sy := float64(theme.Layout.CardHeight) / float64(h)
			opts.GeoM.Scale(sx, sy)
			opts.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(back, opts)
		}
	}

	// Hover overlay (only when not disabled)
	if isHovered && stockCount >= 10 {
		hoverColor := theme.Colors.HoverOverlay
		vector.FillRect(screen, float32(stockX), float32(stockY),
			float32(theme.Layout.CardWidth), float32(theme.Layout.CardHeight),
			hoverColor, false)
	}

	// Disabled overlay when stock < 10 (insufficient for deal)
	if stockCount < 10 {
		disabledColor := theme.Colors.Background // Use background color with higher opacity
		disabledColor.A = 180                    // Semi-transparent
		vector.FillRect(screen, float32(stockX), float32(stockY),
			float32(theme.Layout.CardWidth), float32(theme.Layout.CardHeight),
			disabledColor, false)
	}
}
