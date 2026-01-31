package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/staylor11x/spider-solitaire/internal/assets"
	"github.com/staylor11x/spider-solitaire/internal/deck"
	"github.com/staylor11x/spider-solitaire/internal/game"
	"github.com/staylor11x/spider-solitaire/internal/logger"
)

// Game implements the ebiten.Game interface for Spider Solitaire
type Game struct {
	state     *game.GameState  // Engine state, mutated only in Update
	view      game.GameViewDTO // Read-only snapshot for rendering
	atlas     *CardAtlas       // The cards
	suitCount deck.SuitCount   // Store the difficulty
	theme     *Theme

	lastErr   string // Ephemeral error text
	errFrames int    // Frames left to display lastErr

	// Mouse selection state for click-to-move
	selecting     bool
	selectedPile  int
	selectedIndex int
	showHelp      bool
}

// NewGame create a new Ebiten game instance
func NewGame(suitCount deck.SuitCount) *Game {
	state, err := game.DealInitialGame(suitCount)
	if err != nil {
		panic(err) // TODO: Handle this error gracefully
	}
	atlas, err := NewCardAtlas(assets.Files)
	if err != nil {
		panic(err) // TODO: handle this gracefully too
	}

	view := state.View()

	logger.Info("NewGame: initial deal (stock=%d, completed=%d, won=%v, lost=%v)", view.StockCount, view.CompletedCount, view.Won, view.Lost)

	return &Game{
		state:     state,
		view:      view,
		atlas:     atlas,
		suitCount: suitCount,
		theme:     &DefaultTheme,
		showHelp:  false,
	}
}

// Update runs game logic at 60 FPS
func (g *Game) Update() error {
	g.handleKeyboard()
	g.handleMouse()
	g.tickError()
	return nil
}

func (g *Game) handleKeyboard() {
	// D = deal a row
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		logger.Debug("DealRow: requested")
		if err := g.state.DealRow(); err != nil {
			g.setError(err.Error())
			logger.Error("DealRow: error: %s", err.Error())
		} else {
			g.view = g.state.View()
			g.clearSelection()
			logger.Info("DealRow: success (stock=%d, completed=%d)", g.view.StockCount, g.view.CompletedCount)
		}
	}

	// U = undo last move
	if inpututil.IsKeyJustPressed(ebiten.KeyU) {
		logger.Debug("Undo: requested")
		err := g.state.Undo()
		if err != nil {
			g.setError("No moved to undo")
			logger.Warn("Undo: no history available")
		} else {
			g.view = g.state.View()
			g.selecting = false
			logger.Info("Undo: reverted to previous state")
		}
	}

	// R = reset game
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		logger.Debug("Reset: requested")
		state, err := game.DealInitialGame(g.suitCount)
		if err != nil {
			g.setError(err.Error())
			logger.Error("Reset: error: %s", err.Error())
		} else {
			g.state = state
			g.view = state.View()
			g.clearSelection()
			logger.Info("Reset: success (stock=%d, completed=%d)", g.view.StockCount, g.view.CompletedCount)
		}
	}

	// H = toggle help
	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		g.showHelp = !g.showHelp
	}

	// ESC = cancel selection or close help overlay
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if g.showHelp {
			g.showHelp = false
			logger.Debug("Help overlay closed via ESC key")
		} else if g.selecting {
			g.clearSelection()
			logger.Debug("Selection canceled via ESC key")
		}
	}
}

func (g *Game) handleMouse() {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	mx, my := g.logicalCursor()
	pileIdx, cardIdx, ok := g.hitTest(mx, my)

	if !g.selecting {
		// Start selection on a face-up card
		if !ok {
			return
		}
		// Check for empty pile first
		if len(g.view.Tableau[pileIdx].Cards) == 0 {
			g.setError("cannot select from empty pile")
			logger.Warn("Select: empty pile (pile=%d)", pileIdx)
			return
		}
		if cardIdx >= len(g.view.Tableau[pileIdx].Cards) {
			g.setError("invalid card")
			logger.Warn("Select: invalid card (pile=%d, idx=%d)", pileIdx, cardIdx)
			return
		}
		if !g.view.Tableau[pileIdx].Cards[cardIdx].FaceUp {
			g.setError("select a face-up card")
			logger.Warn("Select: not face-up (pile=%d, idx=%d)", pileIdx, cardIdx)
			return
		}
		g.selecting = true
		g.selectedPile = pileIdx
		g.selectedIndex = cardIdx
		logger.Debug("Select: start (pile=%d, idx=%d)", pileIdx, cardIdx)
		return
	}
	// finish selection, attempt move
	if ok {
		logger.Debug("Move: attempt %d:%d -> %d", g.selectedPile, g.selectedIndex, pileIdx)
		if err := g.performMove(g.selectedPile, g.selectedIndex, pileIdx); err != nil {
			g.setError(err.Error())
			logger.Error("Move: error: %s", err.Error())
		} else {
			g.view = g.state.View()
			logger.Info("Move: success %d:%d -> %d (completed=%d)", g.selectedPile, g.selectedIndex, pileIdx, g.view.CompletedCount)
		}
	} else {
		logger.Debug("Selection canceled by clicking empty space")
	}
	g.clearSelection()
}

// logicalCursor maps the OS/window cursor to logical coordinates
// Ebiten returns cursor positions in Layout-space, so no manual scaling is needed!
func (g *Game) logicalCursor() (lx, ly int) {
	return ebiten.CursorPosition()
}

// hitTest finds the top-most card under the cursor, returning pile and card indices
func (g *Game) hitTest(mx, my int) (pileIdx, cardIdx int, ok bool) {
	for i, pile := range g.view.Tableau {
		x := g.theme.Layout.TableauStartX + i*g.theme.Layout.PileSpacing
		// quick horizontal reject
		if mx < x || mx >= x+g.theme.Layout.CardWidth {
			continue
		}
		y := g.theme.Layout.TableauStartY
		// If the pile is empty, treat clicks within the column's base area as valid
		if len(pile.Cards) == 0 {
			if my >= y && my < y+g.theme.Layout.CardHeight {
				return i, 0, true
			}
			// no cards and click outside base area: continue searching
			continue
		}
		// Cards overlap; check top-most first
		for j := len(pile.Cards) - 1; j >= 0; j-- {
			cy := y + j*g.theme.Layout.CardStackGap
			if mx >= x && mx < x+g.theme.Layout.CardWidth && my >= cy && my < cy+g.theme.Layout.CardHeight {
				return i, j, true
			}
		}
		// if clicking below all cards but within column, treat as click on topmost card area
		// Allow a small margin (one CardStackGap) below the bottom of the top card for easier targeting
		if len(pile.Cards) > 0 {
			topY := y + (len(pile.Cards)-1)*g.theme.Layout.CardStackGap
			bottomY := topY + g.theme.Layout.CardHeight
			clickMarginBelow := g.theme.Layout.CardStackGap // 30px margin
			if mx >= x && mx < x+g.theme.Layout.CardWidth && my >= bottomY && my < bottomY+clickMarginBelow {
				return i, len(pile.Cards) - 1, true
			}
		}
	}
	return 0, 0, false
}

func (g *Game) tickError() {

	// tickle down the error overlay timer
	if g.errFrames > 0 {
		g.errFrames--
		if g.errFrames == 0 {
			g.lastErr = ""
		}
	}
}

func (g *Game) setError(msg string) {
	g.lastErr = msg
	g.errFrames = g.theme.Layout.ErrorDisplayDuration
	logger.Warn("Error: %s", msg)
}

func (g *Game) clearSelection() {
	g.selecting = false
	g.selectedPile = -1
	g.selectedIndex = -1
}

// Draw renders the current frame to the screen
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.theme.Colors.Background)
	drawTableau(screen, g.view, g.atlas, g.theme)

	if g.selecting {
		drawSelectionOverlay(screen, g.view, g.selectedPile, g.selectedIndex, g.theme)
	}

	drawStats(screen, g.view, g.theme)

	if g.lastErr != "" && g.errFrames > 0 {
		drawError(screen, g.lastErr, g.theme)
	}

	if g.view.Won {
		drawWinLossOverlay(screen, "You Win!", g.theme)
	} else if g.view.Lost {
		drawWinLossOverlay(screen, "Game Over :(", g.theme)
	}

	if g.showHelp {
		drawHelpOverlay(screen, g.theme)
	}

}

// Layout returns the logical screen dimensions (fixed).
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.theme.Layout.LogicalWidth, g.theme.Layout.LogicalHeight
}

// performMove executes the engine move via MoveSequence
func (g *Game) performMove(srcPile, startIdx, dstPile int) error {
	return g.state.MoveSequence(srcPile, startIdx, dstPile)
}
