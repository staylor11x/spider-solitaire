package ui

// pileLayout describes computed vertical placement for cards in a tableau pile.
type pileLayout struct {
	Gap   int   // vertical distance between adjacent card origins
	CardY []int // y-origin for each card index in the pile
}

const (
	// defaultMinCompressedCardStackGap is used when theme does not specify a minimum.
	defaultMinCompressedCardStackGap = 10
	// pileBottomPadding reserves pixels below the last card. Keep 0 to maximize usable space.
	pileBottomPadding = 0
)

// computePileLayout calculates per-pile vertical card placement.
//
// Behavior:
// - Uses defaultGap when the pile fits.
// - Compresses gap when the pile would overflow.
// - Prioritizes keeping the last card fully visible on screen.
func computePileLayout(cardCount, startY, cardHeight, logicalHeight, defaultGap, minGap, bottomPadding int) pileLayout {
	if cardCount <= 0 {
		return pileLayout{Gap: defaultGap, CardY: nil}
	}

	if cardCount == 1 {
		return pileLayout{Gap: defaultGap, CardY: []int{startY}}
	}

	available := logicalHeight - bottomPadding - startY - cardHeight
	if available < 0 {
		available = 0
	}

	steps := cardCount - 1
	maxVisibleGap := available / steps // largest gap that still keeps last card visible

	gap := min(maxVisibleGap, defaultGap)

	if gap < minGap {
		// Apply readability floor only when it still fits.
		if minGap*steps <= available {
			gap = minGap
		} else {
			// Visibility has priority over readability floor.
			gap = max(maxVisibleGap, 1)
		}
	}

	cardY := make([]int, cardCount)
	for i := 0; i < cardCount; i++ {
		cardY[i] = startY + i*gap
	}

	return pileLayout{Gap: gap, CardY: cardY}
}

// computeTableauPileLayout computes pile layout from theme defaults.
func computeTableauPileLayout(theme *Theme, cardCount int) pileLayout {
	minGap := defaultMinCompressedCardStackGap
	if theme.Layout.MinCardStackGap > 0 {
		minGap = theme.Layout.MinCardStackGap
	}

	return computePileLayout(
		cardCount,
		theme.Layout.TableauStartY,
		theme.Layout.CardHeight,
		theme.Layout.LogicalHeight,
		theme.Layout.CardStackGap,
		minGap,
		pileBottomPadding,
	)
}
