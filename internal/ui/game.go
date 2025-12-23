package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/staylor11x/spider-solitaire/internal/game"
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
	CardFaceUpColor   = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	CardFaceDownColor = color.RGBA{R: 50, G: 50, B: 150, A: 255}
	TextColor         = color.RGBA{R: 0, G: 0, B: 0, A: 255}
)

// Game implements the ebiten.Game interface for Spider Solitaire
type Game struct {
	state     *game.GameState  // Engine state, mutated only in Update
	view      game.GameViewDTO // Read-only snapshot for rendering
	lastErr   string           // Ephemeral error text
	errFrames int              // Frames left to display lastErr
}

// NewGame create a new Ebiten game instance
func NewGame() *Game {
	state, err := game.DealInitialGame()
	if err != nil {
		panic(err) // TODO: Handle this error gracefully
	}

	view := state.View()

	return &Game{
		state: state,
		view:  view,
	}
}

// Update runs game logic at 60 FPS
func (g *Game) Update() error {

	// D = deal a row
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		if err := g.state.DealRow(); err != nil {
			g.setError(err.Error())
		} else {
			g.view = g.state.View()
		}
	}

	// R = reset game
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		state, err := game.DealInitialGame()
		if err != nil {
			g.setError(err.Error())
		} else {
			g.state = state
			g.view = state.View()
		}
	}

	// tickle down the error overlay timer
	if g.errFrames > 0 {
		g.errFrames--
		if g.errFrames == 0 {
			g.lastErr = ""
		}
	}
	return nil
}

func (g *Game) setError(msg string) {
	g.lastErr = msg
	g.errFrames = errorDisplayDuration
}

// Draw renders the current frame to the screen
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0, G: 100, B: 0, A: 255})

	drawTableau(screen, g.view)
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
