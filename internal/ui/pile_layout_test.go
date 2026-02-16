package ui

import "testing"

func TestComputePileLayout_UsesDefaultGapWhenPileFits(t *testing.T) {
	layout := computePileLayout(
		5,   // cardCount
		150, // startY
		120, // cardHeight
		720, // logicalHeight
		30,  // defaultGap
		8,   // minGap
		0,   // bottomPadding
	)

	if layout.Gap != 30 {
		t.Fatalf("expected default gap 30, got %d", layout.Gap)
	}

	if got, want := layout.CardY[len(layout.CardY)-1], 270; got != want {
		t.Fatalf("expected last card Y %d, got %d", want, got)
	}
}

func TestComputePileLayout_CompressesWhenPileWouldOverflow(t *testing.T) {
	layout := computePileLayout(
		20,  // cardCount
		150, // startY
		120, // cardHeight
		720, // logicalHeight
		30,  // defaultGap
		8,   // minGap
		0,   // bottomPadding
	)

	// available = 450, steps = 19, maxVisibleGap = 23
	if layout.Gap != 23 {
		t.Fatalf("expected compressed gap 23, got %d", layout.Gap)
	}

	lastY := layout.CardY[len(layout.CardY)-1]
	lastBottom := lastY + 120
	if lastBottom > 720 {
		t.Fatalf("expected last card bottom <= 720, got %d", lastBottom)
	}
}

func TestComputePileLayout_VisibilityWinsWhenMinGapCannotFit(t *testing.T) {
	layout := computePileLayout(
		80,  // cardCount
		150, // startY
		120, // cardHeight
		720, // logicalHeight
		30,  // defaultGap
		8,   // minGap
		0,   // bottomPadding
	)

	if layout.Gap < 1 {
		t.Fatalf("expected gap >= 1, got %d", layout.Gap)
	}

	lastY := layout.CardY[len(layout.CardY)-1]
	lastBottom := lastY + 120
	if lastBottom > 720 {
		t.Fatalf("expected last card bottom <= 720, got %d", lastBottom)
	}
}

func TestComputeTableauPileLayout_UsesThemeDefaults(t *testing.T) {
	theme := DefaultTheme
	layout := computeTableauPileLayout(&theme, 20)

	if layout.Gap != 23 {
		t.Fatalf("expected tableau gap 23, got %d", layout.Gap)
	}
}

func TestComputeTableauPileLayout_UsesConfiguredMinGapWhenNeeded(t *testing.T) {
	theme := DefaultTheme
	theme.Layout.MinCardStackGap = 12

	// For 40 cards: available=450, steps=39, maxVisibleGap=11.
	// Min gap cannot fit, so visibility must win and use 11.
	layout := computeTableauPileLayout(&theme, 40)
	if layout.Gap != 11 {
		t.Fatalf("expected visibility gap 11, got %d", layout.Gap)
	}

	// For 30 cards: available=450, steps=29, maxVisibleGap=15.
	// Compression is required, but maxVisibleGap is already above min gap.
	// The algorithm keeps the tightest visible spacing (15).
	layout = computeTableauPileLayout(&theme, 30)
	if layout.Gap != 15 {
		t.Fatalf("expected compressed gap 15, got %d", layout.Gap)
	}
}
