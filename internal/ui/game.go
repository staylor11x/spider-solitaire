package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/staylor11x/spider-solitaire/internal/game"
	"github.com/staylor11x/spider-solitaire/internal/logger"
)

const (
	CardWidth     = 80
	CardHeight    = 120
	PileSpacing   = 100
	CardStackGap  = 30
	TableauStartX = 50
	TableauStartY = 150
	StatsX        = 20
	StatsY        = 20

	errorDisplayDuration = 180
)

var (
	CardFaceUpColor     = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	CardFaceDownColor   = color.RGBA{R: 50, G: 50, B: 150, A: 255}
	TextColor           = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	SelectionOverlayCol = color.RGBA{R: 255, G: 215, B: 0, A: 90}
)

// Game implements the ebiten.Game interface for Spider Solitaire
type Game struct {
	state     *game.GameState  // Engine state, mutated only in Update
	view      game.GameViewDTO // Read-only snapshot for rendering
	lastErr   string           // Ephemeral error text
	errFrames int              // Frames left to display lastErr

	// Mouse selection state for clock-to-move
	selecting     bool
	selectedPile  int
	selectedIndex int
}

// NewGame create a new Ebiten game instance
func NewGame() *Game {
	state, err := game.DealInitialGame()
	if err != nil {
		panic(err) // TODO: Handle this error gracefully
	}

	view := state.View()

	logger.Info("NewGame: initial deal (stock=%d, completed=%d, won=%v, lost=%v)", view.StockCount, view.CompletedCount, view.Won, view.Lost)

	return &Game{
		state: state,
		view:  view,
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

	// R = reset game
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		logger.Debug("Reset: requested")
		state, err := game.DealInitialGame()
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
	}
	g.clearSelection()
}

// logicalCursor maps the OS/window cursor to logical coordinates (handles resize scaling)
func (g *Game) logicalCursor() (lx, ly int) {
	wx, wy := ebiten.WindowSize()
	if wx <= 0 || wy <= 0 {
		return 0, 0
	}
	mx, my := ebiten.CursorPosition()
	lx = mx * 1280 / wx
	ly = my * 720 / wy
	return
}

// hitTest finds the top-most card under the cursor, returning pile and card indices
func (g *Game) hitTest(mx, my int) (pileIdx, cardIdx int, ok bool) {
	for i, pile := range g.view.Tableau {
		x := TableauStartX + i*PileSpacing
		// quick horizontal reject
		if mx < x || mx >= x+CardWidth {
			continue
		}
		y := TableauStartY
		// Cards overlap; check top-most first
		for j := len(pile.Cards) - 1; j >= 0; j-- {
			cy := y + j*CardStackGap
			if mx >= x && mx < x+CardWidth && my >= cy && my < cy+CardHeight {
				return i, j, true
			}
		}
		// if clicking below all cards but within column, treat as click on topmost card area
		if len(pile.Cards) > 0 {
			topY := y + (len(pile.Cards)-1)*CardStackGap
			if mx >= x && mx < x+CardWidth && my >= topY+CardHeight {
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
	g.errFrames = errorDisplayDuration
	logger.Warn("Error: %s", msg)
}

func (g *Game) clearSelection() {
	g.selecting = false
	g.selectedPile = -1
	g.selectedIndex = -1
}

// Draw renders the current frame to the screen
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0, G: 100, B: 0, A: 255})

	drawTableau(screen, g.view)

	if g.selecting {
		drawSelectionOverlay(screen, g.view, g.selectedPile, g.selectedIndex, SelectionOverlayCol)
	}

	drawStats(screen, g.view)

	if g.lastErr != "" && g.errFrames > 0 {
		drawError(screen, g.lastErr)
	}
}

// Layout return the logical screen dimensions
// Ebiten will scale this to fit the actual window size, but you always draw in 1280x720 coordinates
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720
}

// performMove executes the engine move via MoveSequence
func (g *Game) performMove(srcPile, startIdx, dstPile int) error {
	return g.state.MoveSequence(srcPile, startIdx, dstPile)
}
